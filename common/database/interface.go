package database

import "github.com/jmoiron/sqlx"

// Database abstracts the database connection.
// Currently it exposes the underlying sqlx.DB, but this interface allows for future abstraction.
type Database interface {
	GetSqlxDb() *sqlx.DB
	Close()
}
