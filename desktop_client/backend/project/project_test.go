package project_test

import (
	"os"
	"testing"

	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestOpenProjectNotFound(t *testing.T) {
	name := uuid.New().String() + database.DbExtension
	defer os.Remove(name)

	_, err := project.Open(name)
	require.Error(t, err)
}

func TestCreateNewProject(t *testing.T) {
	filename := uuid.New().String() + database.DbExtension
	defer os.Remove(filename)

	name := uuid.New().String()

	project, err := project.Create(filename, name)
	require.NoError(t, err)
	require.Equal(t, project.Settings.Name, name)
	require.Equal(t, project.Filename, filename)
	defer project.Close()
}

func TestReOpenProject(t *testing.T) {
	filename := uuid.New().String() + database.DbExtension
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
	filename := uuid.New().String() + database.DbExtension
	defer os.Remove(filename)

	projectName := uuid.New().String()

	project, err := project.Create(filename, projectName)
	require.NoError(t, err)
	require.Equal(t, project.Settings.Name, projectName)
	defer project.Close()

	name := uuid.New().String()
	c, err := project.CreateCharacter(name)
	require.NoError(t, err)

	require.Equal(t, name, c.Name)
	require.NotEmpty(t, c.Created)
	require.NotEmpty(t, c.Id)
}
