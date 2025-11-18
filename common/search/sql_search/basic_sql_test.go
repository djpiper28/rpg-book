package sqlsearch_test

import (
	"os"
	"testing"

	"github.com/djpiper28/rpg-book/common/search/parser"
	sqlsearch "github.com/djpiper28/rpg-book/common/search/sql_search"
	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func parse(t *testing.T, query string) *parser.Node {
	t.Helper()

	ast, err := parser.Parse(query)
	require.NoError(t, err)
	require.NotEmpty(t, ast)
	return ast
}

func asSql[T any](t *testing.T, query string, tableData sqlsearch.SqlTableData, columnMap sqlsearch.SqlColmnMap) (string, []any) {
	t.Helper()

	res, args, err := sqlsearch.AsSql[T](parse(t, query), tableData, columnMap)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.NotEmpty(t, args)

	t.Logf("Query: %s", query)
	t.Logf("SQL: %s", res)
	t.Logf("Args: %+v", args)
	return res, args
}

func TestBasicQueryText(t *testing.T) {
	type TestTable struct {
		Name string `db:"name"`
	}

	querySql, sqlArgs := asSql[TestTable](t, "test", sqlsearch.SqlTableData{
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

	dbName := uuid.NewString()
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
	type TestTable struct {
		Name string `db:"name"`
	}

	querySql, sqlArgs := asSql[TestTable](t, "name:\"123\"", sqlsearch.SqlTableData{
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

	dbName := uuid.NewString()
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
