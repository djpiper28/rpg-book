package sqlsearch_test

import (
	"testing"

	"github.com/djpiper28/rpg-book/common/database/sqlite3"
	"github.com/stretchr/testify/require"
)

func doBooleanSqltest(t *testing.T, db *sqlite3.Db, query string, expected int) {
	t.Helper()

	sql, args := asSql(t, query, TestTableData, TestColumnMap)

	var people []TestModel
	err := db.Db.Select(&people, sql, args...)
	require.NoError(t, err)
	require.Len(t, people, expected)
}

func TestSqlBooleanEquals(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	t.Run("t value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats=t", 4)
	})

	t.Run("true value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats=true", 4)
	})

	t.Run("mixed case true value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats=True", 4)
		doBooleanSqltest(t, db, "likes_cats=TruE", 4)
		doBooleanSqltest(t, db, "likes_cats=TRUE", 4)
	})

	t.Run("1 value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats=1", 4)
	})

	t.Run("f value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats=f", 1)
	})

	t.Run("false value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats=false", 1)
	})

	t.Run("mixed case false value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats=False", 1)
		doBooleanSqltest(t, db, "likes_cats=FALse", 1)
		doBooleanSqltest(t, db, "likes_cats=FALSE", 1)
	})

	t.Run("0 value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats=0", 1)
	})
}

func TestSqlBooleanNotEquals(t *testing.T) {
	t.Parallel()

	db, close := newTestDb(t)
	defer close()

	t.Run("t value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats~t", 1)
	})

	t.Run("true value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats~true", 1)
	})

	t.Run("mixed case true value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats~True", 1)
		doBooleanSqltest(t, db, "likes_cats~TruE", 1)
		doBooleanSqltest(t, db, "likes_cats~TRUE", 1)
	})

	t.Run("1 value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats~1", 1)
	})

	t.Run("f value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats~f", 4)
	})

	t.Run("false value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats~false", 4)
	})

	t.Run("mixed case false value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats~False", 4)
		doBooleanSqltest(t, db, "likes_cats~FALse", 4)
		doBooleanSqltest(t, db, "likes_cats~FALSE", 4)
	})

	t.Run("0 value", func(t *testing.T) {
		doBooleanSqltest(t, db, "likes_cats~0", 4)
	})
}
