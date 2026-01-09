package project_test

import (
	"fmt"
	"testing"

	imagecompression "github.com/djpiper28/rpg-book/common/image/image_compression"
	"github.com/djpiper28/rpg-book/common/normalisation"
	testutils "github.com/djpiper28/rpg-book/common/test_utils"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

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
func TestSearchCharacterQualified(t *testing.T) {
	t.Parallel()

	project, remove := testutils.NewProject(t)
	defer remove()

	name := uuid.New().String()
	desc := uuid.New().String()
	character, err := project.CreateCharacter(name, desc, nil)
	require.NoError(t, err)

	res, err := project.SearchCharacter(fmt.Sprintf("character.description:%s", desc))
	require.NoError(t, err)
	require.Len(t, res, 1)
	require.Equal(t, character.Id, res[0])
}
