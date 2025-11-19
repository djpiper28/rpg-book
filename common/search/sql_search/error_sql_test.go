package sqlsearch_test

import (
	"testing"

	sqlsearch "github.com/djpiper28/rpg-book/common/search/sql_search"
	"github.com/stretchr/testify/require"
)

func TestMissingColumn(t *testing.T) {
	t.Parallel()

	tableData := sqlsearch.SqlTableData{
		FieldsToScan: []string{"name"},
		TableName:    "test",
	}
	columnMap := sqlsearch.SqlColmnMap{
		TextColumns: map[string]string{"name": "name"},
	}

	// "age" is not in the map
	// We cannot use asSql helper here because it asserts NoError
	_, _, err := sqlsearch.AsSql(parse(t, `age:10`), tableData, columnMap)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not in the column map")
}

func TestSqlInjection(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	maliciousInput := "'; DROP TABLE test; --"
	// We use a text column so it goes through the text path
	// The query will be: name LIKE '%'; DROP TABLE test; --%'
	// This is a valid string literal for LIKE, so it should execute as a search.
	// It should NOT execute the DROP TABLE command.
	sql, args := asSql(t, `name:"`+maliciousInput+`"`, TestTableData, TestColumnMap)

	// The SQL should NOT contain the malicious input directly (it should be parameterized)
	require.NotContains(t, sql, "DROP TABLE")
	require.Contains(t, sql, "name_normalised LIKE ?")

	// Execute the query
	rows, err := db.Db.Queryx(sql, args...)
	require.NoError(t, err)
	rows.Close()

	// Verify table still exists and has data
	var count int
	err = db.Db.Get(&count, "SELECT count(*) FROM test")
	require.NoError(t, err)
	require.Equal(t, 5, count) // newTestDb inserts 5 rows
}

func TestInvalidIncludesSyntax(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	// Case 1: Just a basic query with weird chars
	// The query will be: name LIKE '%name_normalised LIKE ?%'
	// This tests that special characters in the input are treated as literals in the search
	weirdInput := "%_'"
	sql, args := asSql(t, `"`+weirdInput+`"`, TestTableData, TestColumnMap)

	require.Contains(t, sql, "name LIKE ?")
	require.Len(t, args, 1)
	// It should be wrapped in %
	require.Equal(t, "%"+weirdInput+"%", args[0])

	// Execute the query - it should not error, but might not find anything
	rows, err := db.Db.Queryx(sql, args...)
	require.NoError(t, err)
	rows.Close()
}
