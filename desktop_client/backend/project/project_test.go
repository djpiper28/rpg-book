package project_test

import (
	"os"
	"testing"

	"github.com/djpiper28/rpg-book/common/database/sqlite3"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project"
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
