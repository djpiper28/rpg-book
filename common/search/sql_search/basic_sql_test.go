package sqlsearch_test

import (
	"os"
	"testing"

	sqlsearch "github.com/djpiper28/rpg-book/common/search/sql_search"
	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBasicQueryText(t *testing.T) {
	t.Parallel()

	type TestTable struct {
		Name string `db:"name"`
	}

	querySql, sqlArgs := asSql(t, "test", sqlsearch.SqlTableData{
		FieldsToScan: []string{"name"},
		TableName:    "test",
	},
		sqlsearch.SqlColmnMap{
			BasicQueryColumn: "name",
			TextColumns: map[string]string{
				"name": "name",
			},
		})

	require.Equal(t, `SELECT name
FROM test
WHERE
name LIKE ?;`, querySql)

	dbName := uuid.NewString() + database.DbExtension
	defer os.Remove(dbName)

	db, err := database.New(dbName)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Db.Exec("CREATE TABLE test (name TEXT NOT NULL);")
	require.NoError(t, err)

	_, err = db.Db.Exec("INSERT INTO test (name) VALUES('testing123');")
	require.NoError(t, err)

	rows, err := db.Db.Queryx(querySql, sqlArgs...)
	require.NoError(t, err)

	require.True(t, rows.Next())

	var item TestTable
	require.NoError(t, rows.StructScan(&item))
}

func TestSetGenerator(t *testing.T) {
	t.Parallel()

	type TestTable struct {
		Name string `db:"name"`
	}

	querySql, sqlArgs := asSql(t, "name:\"123\"", sqlsearch.SqlTableData{
		FieldsToScan: []string{"name"},
		TableName:    "test",
	},
		sqlsearch.SqlColmnMap{
			BasicQueryColumn: "name",
			TextColumns: map[string]string{
				"name": "name",
			},
		})

	require.Equal(t, `SELECT name
FROM test
WHERE
name LIKE ?;`, querySql)

	dbName := uuid.NewString() + database.DbExtension
	defer os.Remove(dbName)

	db, err := database.New(dbName)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Db.Exec("CREATE TABLE test (name TEXT NOT NULL);")
	require.NoError(t, err)

	_, err = db.Db.Exec("INSERT INTO test (name) VALUES('testing123');")
	require.NoError(t, err)

	rows, err := db.Db.Queryx(querySql, sqlArgs...)
	require.NoError(t, err)

	require.True(t, rows.Next())

	var item TestTable
	require.NoError(t, rows.StructScan(&item))
}
