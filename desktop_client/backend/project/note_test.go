package project_test

import (
	"fmt"
	"testing"

	"github.com/djpiper28/rpg-book/common/normalisation"
	testutils "github.com/djpiper28/rpg-book/common/test_utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateNote(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	name := uuid.New().String()
	markdown := uuid.New().String()
	note, err := project.CreateNote(name, markdown, []uuid.UUID{})
	require.NoError(t, err)
	require.Equal(t, name, note.Name)
	require.Equal(t, markdown, note.Markdown)
	require.Equal(t, normalisation.Normalise(name), note.NameNormalised)
	require.Equal(t, normalisation.Normalise(markdown), note.MarkdownNormalised)
	require.NotEmpty(t, note.Created)
	require.NotEmpty(t, note.Id)

	res, err := project.GetNotes()
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, note, res[0])
}

func TestCreateNoteWithRelatedCharacters(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	const relatedIds = 10
	characterIds := make([]uuid.UUID, 0)

	for i := range relatedIds {
		character, err := project.CreateCharacter(uuid.New().String(), fmt.Sprintf("This is a test %d", i), nil)
		require.NoError(t, err)
		require.NotNil(t, character)

		characterIds = append(characterIds, character.Id)
	}

	name := uuid.New().String()
	markdown := uuid.New().String()
	note, err := project.CreateNote(name, markdown, characterIds)
	require.NoError(t, err)
	require.Equal(t, name, note.Name)
	require.Equal(t, markdown, note.Markdown)
	require.Equal(t, normalisation.Normalise(name), note.NameNormalised)
	require.Equal(t, normalisation.Normalise(markdown), note.MarkdownNormalised)
	require.NotEmpty(t, note.Created)
	require.NotEmpty(t, note.Id)

	res, err := project.GetNotes()
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, note, res[0])

	for _, id := range characterIds {
		character, err := project.GetCharacter(id)
		require.NoError(t, err)
		require.Len(t, character.Notes, 1)
		require.Equal(t, note, character.Notes[0])
	}
}

func TestGetNote(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	character, err := project.CreateCharacter("Crazy Dave", "pan on head", nil)
	require.NoError(t, err)
	require.NotEmpty(t, character.Id)

	name := uuid.New().String()
	markdown := uuid.New().String()
	expectedNote, err := project.CreateNote(name, markdown, []uuid.UUID{character.Id})
	require.NoError(t, err)

	require.NotEmpty(t, expectedNote.Id)
	require.Equal(t, name, expectedNote.Name)
	require.Equal(t, markdown, expectedNote.Markdown)

	completeNote, err := project.GetNote(expectedNote.Id)
	require.NoError(t, err)
	require.Equal(t, expectedNote, completeNote.Note)
	require.Len(t, completeNote.Characters, 1)
	require.Equal(t, character, completeNote.Characters[0])
}

func TestSearchNoteBasic(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	character, err := project.CreateCharacter("Crazy Dave", "pan on head", nil)
	require.NoError(t, err)
	require.NotEmpty(t, character.Id)

	name := uuid.New().String()
	markdown := uuid.New().String()
	expectedNote, err := project.CreateNote(name, markdown, []uuid.UUID{character.Id})
	require.NoError(t, err)

	ids, err := project.SearchNote(name)
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{expectedNote.Id}, ids)
}

func TestSearchNoteName(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	character, err := project.CreateCharacter("Crazy Dave", "pan on head", nil)
	require.NoError(t, err)
	require.NotEmpty(t, character.Id)

	name := uuid.New().String()
	markdown := uuid.New().String()
	expectedNote, err := project.CreateNote(name, markdown, []uuid.UUID{character.Id})
	require.NoError(t, err)

	ids, err := project.SearchNote(fmt.Sprintf("name:%s", name))
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{expectedNote.Id}, ids)
}

func TestSearchNoteMarkdown(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	character, err := project.CreateCharacter("Crazy Dave", "pan on head", nil)
	require.NoError(t, err)
	require.NotEmpty(t, character.Id)

	name := uuid.New().String()
	markdown := uuid.New().String()
	expectedNote, err := project.CreateNote(name, markdown, []uuid.UUID{character.Id})
	require.NoError(t, err)

	ids, err := project.SearchNote(fmt.Sprintf("markdown:%s", markdown))
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{expectedNote.Id}, ids)

	ids, err = project.SearchNote(fmt.Sprintf("contents:%s", markdown))
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{expectedNote.Id}, ids)

	ids, err = project.SearchNote(fmt.Sprintf("desc:%s", markdown))
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{expectedNote.Id}, ids)

	ids, err = project.SearchNote(fmt.Sprintf("description:%s", markdown))
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{expectedNote.Id}, ids)
}

func TestSearchNoteQualified(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	character, err := project.CreateCharacter("Crazy Dave", "pan on head", nil)
	require.NoError(t, err)
	require.NotEmpty(t, character.Id)

	name := uuid.New().String()
	markdown := uuid.New().String()
	expectedNote, err := project.CreateNote(name, markdown, []uuid.UUID{character.Id})
	require.NoError(t, err)

	ids, err := project.SearchNote(fmt.Sprintf("note.markdown:%s", markdown))
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{expectedNote.Id}, ids)
}

func TestSearchNoteCharacterName(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	character, err := project.CreateCharacter("Crazy Dave", "pan on head", nil)
	require.NoError(t, err)
	require.NotEmpty(t, character.Id)

	name := uuid.New().String()
	markdown := uuid.New().String()
	expectedNote, err := project.CreateNote(name, markdown, []uuid.UUID{character.Id})
	require.NoError(t, err)

	ids, err := project.SearchNote("character.name:dave")
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{expectedNote.Id}, ids)
}

func TestSearchNoteCharacterDescription(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	character, err := project.CreateCharacter("Crazy Dave", "pan on head", nil)
	require.NoError(t, err)
	require.NotEmpty(t, character.Id)

	name := uuid.New().String()
	markdown := uuid.New().String()
	expectedNote, err := project.CreateNote(name, markdown, []uuid.UUID{character.Id})
	require.NoError(t, err)

	ids, err := project.SearchNote("character.description:pan")
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{expectedNote.Id}, ids)

	ids, err = project.SearchNote("character.desc:pan")
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{expectedNote.Id}, ids)
}
