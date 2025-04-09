package database_test

import (
	"os"
	"testing"

	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDbSavesToDisk(t *testing.T) {
	dbName := uuid.New().String() + database.DbExtension
	defer os.Remove(dbName)

	db, err := database.New(dbName)
	require.NoError(t, err)

	_, err = db.Db.Exec(`CREATE TABLE test (
    id string primary key,
    test string
);`)
	require.NoError(t, err)
	db.Close()

	_, err = os.Lstat(dbName)
	require.NoError(t, err)
}
