package sqlsearch_test

import (
	"os"
	"testing"

	"github.com/djpiper28/rpg-book/common/normalisation"
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

func asSql(t *testing.T, query string, tableData sqlsearch.SqlTableData, columnMap sqlsearch.SqlColmnMap) (string, []any) {
	t.Helper()

	res, args, err := sqlsearch.AsSql(parse(t, query), tableData, columnMap)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.NotEmpty(t, args)

	t.Logf("Query: %s", query)
	t.Logf("SQL: %s", res)
	t.Logf("Args: %+v", args)
	return res, args
}

type TestModel struct {
	Id             int     `db:"id"`
	Name           string  `db:"name"`
	NameNormalised string  `db:"name_normalised"`
	Age            float32 `db:"age"`
	Gender         string  `db:"gender"`
	LikesCats      bool    `db:"likes_cats"`
}

func newTestDb(t *testing.T) (*database.Db, func()) {
	t.Helper()

	dbName := uuid.NewString() + database.DbExtension
	db, err := database.New(dbName)
	require.NoError(t, err)

	closeFn := func() {
		db.Close()
		os.Remove(dbName)
	}

	_, err = db.Db.Exec(`
CREATE TABLE test (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  name_normalised TEXT NOT NULL,
  age FLOAT,
  gender TEXT DEFAULT('non-binary'),
  likes_cats BOOLEAN
);
    `)
	require.NoError(t, err)

	people := []TestModel{
		{
			Name:      "Danny Piper",
			Age:       23,
			Gender:    "Male (mostly)",
			LikesCats: true,
		},
		{
			Name:      "Steph",
			Age:       25, // ish
			Gender:    "Female",
			LikesCats: true,
		},
		{
			Name:      "John Smith",
			Age:       40,
			Gender:    "Male",
			LikesCats: false,
		},
		{
			Name:      "Leondro Lio",
			Age:       26,
			Gender:    "Male",
			LikesCats: true,
		},
		{
			Name:      "Amber Piper",
			Age:       22,
			Gender:    "Female",
			LikesCats: true,
		},
	}

	for i, person := range people {
		person.Id = i
		person.NameNormalised = normalisation.Normalise(person.Name)
		_, err := db.Db.NamedExec("INSERT INTO test (id, name, name_normalised, age, gender, likes_cats) VALUES(:id, :name, :name_normalised, :age, :gender, :likes_cats);", person)
		require.NoError(t, err)
	}

	return db, closeFn
}

var TestTableData sqlsearch.SqlTableData = sqlsearch.SqlTableData{
	FieldsToScan: []string{"name", "age", "gender", "likes_cats"},
	TableName:    "test",
}

var TestColumnMap sqlsearch.SqlColmnMap = sqlsearch.SqlColmnMap{
	TextColumns: map[string]string{
		"name": "name_normalised",
		"nom":  "name_normalised",
	},
	NumberColumns: map[string]string{
		"age": "age",
		"ans": "age",
	},
	BasicQueryColumn:    "name", // This tests the alias system
	BasicQueryOperation: parser.GeneratorOperator_Includes,
}
