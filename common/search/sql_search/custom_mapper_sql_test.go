package sqlsearch_test

import (
	"errors"
	"testing"

	"github.com/djpiper28/rpg-book/common/search/parser"
	sqlsearch "github.com/djpiper28/rpg-book/common/search/sql_search"
	"github.com/stretchr/testify/require"
)

type TestCustomColumn struct{}

func (TestCustomColumn) Map(operator parser.GeneratorOperator, value string) (query string, args []any, err error) {
	if operator != parser.GeneratorOperator_Equals {
		return "", nil, errors.New("Invalid operator for custom column, it must be =")
	}

	return "(friends.friend_a_id = (SELECT users.id FROM users WHERE name LIKE ?) OR friends.friend_b_id = (SELECT users.id FROM users WHERE name LIKE ?))", []any{value, value}, nil
}

func TestJoinAndCustomSql(t *testing.T) {
	t.Parallel()

	db, close := newJoinTestDb(t)
	defer close()

	tableData := sqlsearch.SqlTableData{
		FieldsToScan: []string{"users.name"},
		TableName:    "users",
		JoinClauses:  "JOIN friends ON (friends.friend_a_id = users.id) OR (friend_b_id = users.id)",
	}

	columnMap := sqlsearch.SqlColmnMap{
		TextColumns: map[string]string{
			"name": "users.name",
		},
		BasicQueryColumn:    "users.name",
		BasicQueryOperation: parser.GeneratorOperator_Includes,
		CustomColumns: map[string]sqlsearch.CustomColumn{
			"friend": TestCustomColumn{},
		},
	}

	sql, args := asSql(t, `friend=Alice`, tableData, columnMap)

	rows, err := db.Db.Queryx(sql, args...)
	require.NoError(t, err)
	defer rows.Close()

	var results []User
	for rows.Next() {
		var res User
		err := rows.StructScan(&res)
		require.NoError(t, err)
		results = append(results, res)
	}

	require.Contains(t, results, User{Name: "Alice"})
	require.Contains(t, results, User{Name: "Bob"})
	require.Len(t, results, 2)
}
