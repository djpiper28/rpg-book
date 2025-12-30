package project_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/djpiper28/rpg-book/common/database/sqlite3"
	imagecompression "github.com/djpiper28/rpg-book/common/image/image_compression"
	"github.com/djpiper28/rpg-book/common/normalisation"
	testutils "github.com/djpiper28/rpg-book/common/test_utils"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestOpenProjectNotFound(t *testing.T) {
	t.Parallel()

	dbName := uuid.NewString() + sqlite3.DbExtension
	defer os.Remove(dbName)

	_, err := project.Open(dbName)
	require.Error(t, err)
}

func TestCreateNewProject(t *testing.T) {
	t.Parallel()

	filename := uuid.New().String() + sqlite3.DbExtension
	defer os.Remove(filename)

	name := uuid.New().String()
	project, err := project.Create(filename, name)
	require.NoError(t, err)
	require.Equal(t, project.Filename, filename)

	require.Equal(t, project.Settings.Name, name)
	defer project.Close()
}

func TestReOpenProject(t *testing.T) {
	t.Parallel()

	filename := uuid.New().String() + sqlite3.DbExtension
	defer os.Remove(filename)

	name := uuid.New().String()

	project1, err := project.Create(filename, name)
	require.NoError(t, err)
	require.Equal(t, project1.Settings.Name, name)
	project1.Close()

	project2, err := project.Open(filename)
	require.NoError(t, err)
	require.Equal(t, project2.Settings.Name, name)
	defer project2.Close()

	require.Equal(t, project1.Filename, filename)
	require.Equal(t, project2.Filename, filename)
}

func TestCreateCharacter(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	name := uuid.New().String()
	desc := uuid.New().String()
	icon := []byte{0x01, 0x02, 0x03, 0x04}
	c, err := project.CreateCharacter(name, desc, icon)
	require.NoError(t, err)

	require.Equal(t, name, c.Name)
	require.Equal(t, desc, c.Description)
	require.Equal(t, normalisation.Normalise(name), c.NameNormalised)
	require.Equal(t, desc, c.Description)
	require.Equal(t, normalisation.Normalise(desc), c.DescriptionNormalised)
	require.Equal(t, icon, c.Icon)
	require.NotEmpty(t, c.Created)
	require.NotEmpty(t, c.Id)

	characters, err := project.GetCharacters()
	require.NoError(t, err)
	require.Len(t, characters, 1)
	require.Equal(t, c, characters[0])
}

func TestCreateCharacterNilIcon(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	name := uuid.New().String()
	desc := uuid.New().String()
	var icon []byte = nil
	c, err := project.CreateCharacter(name, desc, icon)
	require.NoError(t, err)

	require.Equal(t, name, c.Name)
	require.Equal(t, desc, c.Description)
	require.Equal(t, icon, c.Icon)
	require.NotEmpty(t, c.Created)
	require.NotEmpty(t, c.Id)

	characters, err := project.GetCharacters()
	require.NoError(t, err)
	require.Len(t, characters, 1)
	require.Equal(t, c, characters[0])
}

func TestUpdateCharacter(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	name := uuid.New().String()
	desc := uuid.New().String()
	character, err := project.CreateCharacter(name, desc, nil)
	require.NoError(t, err)

	img := testutils.NewTestImage(100, 100)
	icon, err := imagecompression.CompressIcon(img)
	require.NoError(t, err)
	name = "new name " + uuid.New().String()
	description := "new description " + uuid.New().String()

	character.Name = name
	character.Icon = icon
	character.Description = description
	character.Notes = make([]*model.Note, 0)

	err = project.UpdateCharacter(character, true)
	require.NoError(t, err)

	readCharacter, err := project.GetCharacter(character.Id)
	require.NoError(t, err)
	require.Equal(t, *character, *readCharacter)
}

func TestUpdateCharacterNoIconChange(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	name := uuid.New().String()
	desc := uuid.New().String()

	img := testutils.NewTestImage(100, 100)
	icon, err := imagecompression.CompressIcon(img)
	require.NoError(t, err)

	character, err := project.CreateCharacter(name, desc, icon)
	require.NoError(t, err)

	name = "new name " + uuid.New().String()
	description := "new description " + uuid.New().String()

	character.Name = name
	character.Icon = nil
	character.Description = description
	character.Notes = make([]*model.Note, 0)

	err = project.UpdateCharacter(character, false)
	require.NoError(t, err)

	readCharacter, err := project.GetCharacter(character.Id)
	require.NoError(t, err)
	character.Icon = icon
	require.Equal(t, *character, *readCharacter)
}

func TestDeleteCharacter(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	name := uuid.New().String()
	desc := uuid.New().String()
	character, err := project.CreateCharacter(name, desc, nil)
	require.NoError(t, err)

	err = project.DeleteCharacter(character.Id)
	require.NoError(t, err)

	_, err = project.GetCharacter(character.Id)
	require.Error(t, err)
}

func TestSearchCharacterBasic(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	name := uuid.New().String()
	desc := uuid.New().String()
	character, err := project.CreateCharacter(name, desc, nil)
	require.NoError(t, err)

	res, err := project.SearchCharacter(name)
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, character.Id, res[0])
}

func TestSearchCharacterName(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	name := uuid.New().String()
	desc := uuid.New().String()
	character, err := project.CreateCharacter(name, desc, nil)
	require.NoError(t, err)

	res, err := project.SearchCharacter(fmt.Sprintf("name:%s", name))
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, character.Id, res[0])
}

func TestSearchCharacterDesc(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	name := uuid.New().String()
	desc := uuid.New().String()
	character, err := project.CreateCharacter(name, desc, nil)
	require.NoError(t, err)

	res, err := project.SearchCharacter(fmt.Sprintf("desc:%s", desc))
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, character.Id, res[0])
}

func TestSearchCharacterDescription(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	name := uuid.New().String()
	desc := uuid.New().String()
	character, err := project.CreateCharacter(name, desc, nil)
	require.NoError(t, err)

	res, err := project.SearchCharacter(fmt.Sprintf("description:%s", desc))
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, character.Id, res[0])
}

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
