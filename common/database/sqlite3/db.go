package sqlite3

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	Db       *sqlx.DB
	filename string
}

const (
	dbTimeout       = time.Second * 10
	DbExtension     = ".rpg"
	migrationsTable = `
CREATE TABLE IF NOT EXISTS migrations (
  version INTEGER PRIMARY KEY,
  date TIMESTAMPTZ NOT NULL
);
  `
)

type Migrations struct {
	Version int       `db:"version"`
	Date    time.Time `db:"date"`
}

func New(filename string) (*Db, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "sqlite3", fmt.Sprintf("file:%s?_journal_mode=WAL&parseTime=true&_auto_vacuum=2", filename))
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Cannot create database connection for %s", filename), err)
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Cannot ping database %s", filename), err)
	}

	log.Info("Connected to new database", loggertags.TagFileName, filename)
	db.SetMaxOpenConns(1)

	ret := &Db{
		Db:       db,
		filename: filename,
	}

	err = ret.prepareMigrations()
	if err != nil {
		return nil, errors.Join(errors.New("Cannot check initial migrations table"), err)
	}

	return ret, nil
}

func (d *Db) prepareMigrations() error {
	_, err := d.Db.Exec(migrationsTable)
	if err != nil {
		return errors.Join(errors.New("Cannot check for migrations table"), err)
	}

	return nil
}

func (d *Db) Close() {
	d.Db.Close()
	log.Info("Disconnected from database", loggertags.TagFileName, d.filename)
}

func (d *Db) GetSqlxDb() *sqlx.DB {
	return d.Db
}
