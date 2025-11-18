package sqlsearch_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSqlTextEqual(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	sql, args := asSql(t, `name="DANNY piper"`, TestTableData, TestColumnMap)

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

func TestSqlTextNotEqual(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	sql, args := asSql(t, `name~"Danny Piper"`, TestTableData, TestColumnMap)

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

func TestSqlTextBasic1(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	sql, args := asSql(t, `"Danny Piper"`, TestTableData, TestColumnMap)

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

func TestSqlTextBasic2(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	sql, args := asSql(t, `Danny Piper`, TestTableData, TestColumnMap)

  require.Contains(t, sql, "AND")

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

func TestSqlTextAlias(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	sql, args := asSql(t, `nom="DANNY piper"`, TestTableData, TestColumnMap)

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
