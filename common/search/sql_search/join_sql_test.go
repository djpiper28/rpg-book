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

type User struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

type Order struct {
	Id     int    `db:"id"`
	UserId int    `db:"user_id"`
	Item   string `db:"item"`
}

type UserOrder struct {
	Name string `db:"name"`
	Item string `db:"item"`
}

func newJoinTestDb(t *testing.T) (*database.Db, func()) {
	t.Helper()

	dbName := uuid.NewString()
	db, err := database.New(dbName)
	require.NoError(t, err)

	closeFn := func() {
		db.Close()
		os.Remove(dbName)
	}

	_, err = db.Db.Exec(`
CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE orders (
  id INTEGER PRIMARY KEY,
  user_id INTEGER,
  item TEXT NOT NULL,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
    `)
	require.NoError(t, err)

	users := []User{
		{Id: 1, Name: "Alice"},
		{Id: 2, Name: "Bob"},
	}

	for _, user := range users {
		_, err := db.Db.NamedExec("INSERT INTO users (id, name) VALUES(:id, :name);", user)
		require.NoError(t, err)
	}

	orders := []Order{
		{Id: 1, UserId: 1, Item: "Apple"},
		{Id: 2, UserId: 2, Item: "Banana"},
		{Id: 3, UserId: 1, Item: "Apricot"},
	}

	for _, order := range orders {
		_, err := db.Db.NamedExec("INSERT INTO orders (id, user_id, item) VALUES(:id, :user_id, :item);", order)
		require.NoError(t, err)
	}

	return db, closeFn
}

func TestJoinSql(t *testing.T) {
	t.Parallel()

	db, close := newJoinTestDb(t)
	defer close()

	tableData := sqlsearch.SqlTableData{
		FieldsToScan: []string{"users.name", "orders.item"},
		TableName:    "users",
		JoinClauses:  "INNER JOIN orders ON users.id = orders.user_id",
	}

	columnMap := sqlsearch.SqlColmnMap{
		TextColumns: map[string]string{
			"item": "orders.item",
			"name": "users.name",
		},
		BasicQueryColumn:    "users.name",
		BasicQueryOperation: parser.GeneratorOperator_Includes,
	}

	// Test finding orders for Alice (Apple, Apricot)
	sql, args := asSql(t, `item:"Apple"`, tableData, columnMap)

	rows, err := db.Db.Queryx(sql, args...)
	require.NoError(t, err)
	defer rows.Close()

	var results []UserOrder
	for rows.Next() {
		var res UserOrder
		err := rows.StructScan(&res)
		require.NoError(t, err)
		results = append(results, res)
	}

	require.Len(t, results, 1)
	require.Equal(t, "Alice", results[0].Name)
	require.Equal(t, "Apple", results[0].Item)

	// Test finding orders for Bob (Banana)
	sql, args = asSql(t, `item:"Banana"`, tableData, columnMap)

	rows, err = db.Db.Queryx(sql, args...)
	require.NoError(t, err)
	defer rows.Close()

	results = nil
	for rows.Next() {
		var res UserOrder
		err := rows.StructScan(&res)
		require.NoError(t, err)
		results = append(results, res)
	}

	require.Len(t, results, 1)
	require.Equal(t, "Bob", results[0].Name)
	require.Equal(t, "Banana", results[0].Item)
}
