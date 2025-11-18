package sqlsearch_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSqlGreaterThan(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	sql, args := asSql(t, "age>3", TestTableData, TestColumnMap)

	rows, err := db.Db.Queryx(sql, args...)
	require.NoError(t, err)

	var people []TestModel
	for rows.Next() {
		var person TestModel
		err := rows.StructScan(&person)
		require.NoError(t, err)

		people = append(people, person)
	}

	require.Len(t, people, 5)
}

func TestSqlGreaterThanOrEqual(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	sql, args := asSql(t, "age>=40", TestTableData, TestColumnMap)

	rows, err := db.Db.Queryx(sql, args...)
	require.NoError(t, err)

	var people []TestModel
	for rows.Next() {
		var person TestModel
		err := rows.StructScan(&person)
		require.NoError(t, err)

		people = append(people, person)
	}

	require.Len(t, people, 1)
}

func TestSqlLessThan(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	sql, args := asSql(t, "age<25", TestTableData, TestColumnMap)

	rows, err := db.Db.Queryx(sql, args...)
	require.NoError(t, err)

	var people []TestModel
	for rows.Next() {
		var person TestModel
		err := rows.StructScan(&person)
		require.NoError(t, err)

		people = append(people, person)
	}

	require.Len(t, people, 2)
}

func TestSqlLessThanOrEqual(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	sql, args := asSql(t, "age<=26", TestTableData, TestColumnMap)

	rows, err := db.Db.Queryx(sql, args...)
	require.NoError(t, err)

	var people []TestModel
	for rows.Next() {
		var person TestModel
		err := rows.StructScan(&person)
		require.NoError(t, err)

		people = append(people, person)
	}

	require.Len(t, people, 4)
}

func TestSqlNumberEqual(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	sql, args := asSql(t, "age=22", TestTableData, TestColumnMap)

	rows, err := db.Db.Queryx(sql, args...)
	require.NoError(t, err)

	var people []TestModel
	for rows.Next() {
		var person TestModel
		err := rows.StructScan(&person)
		require.NoError(t, err)

		people = append(people, person)
	}

	require.Len(t, people, 1)
}

func TestSqlNotNumberEqual(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	sql, args := asSql(t, "age~22", TestTableData, TestColumnMap)

	rows, err := db.Db.Queryx(sql, args...)
	require.NoError(t, err)

	var people []TestModel
	for rows.Next() {
		var person TestModel
		err := rows.StructScan(&person)
		require.NoError(t, err)

		people = append(people, person)
	}

	require.Len(t, people, 4)
}
