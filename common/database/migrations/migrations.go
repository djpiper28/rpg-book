package migrations

import (
	"database/sql"
	"errors"
	"time"

	"github.com/charmbracelet/log"
	"github.com/djpiper28/rpg-book/common/database"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/jmoiron/sqlx"
)

type Migration struct {
	PreProcess  func(*sqlx.Tx) error
	Sql         string
	PostProcess func(*sqlx.Tx) error
	Test        func(*sqlx.Tx) error // Don't edit the database in tests plx
}

type DbMigrator struct {
	migrations []Migration
	Rebinder   int // See https://github.com/jmoiron/sqlx/blob/master/bind.go, defaults to ? (sqlite)
}

// You will need to change the Rebinder if you are not using sqlite, for example:
// if you are using Postgres you need to set the DbMigrator.Rebinder to DOLLAR.
func New(migrations []Migration) *DbMigrator {
	return &DbMigrator{
		migrations: migrations,
		Rebinder:   sqlx.QUESTION,
	}
}

func (m *DbMigrator) Migrate(db database.Database) error {
	tx, err := db.GetSqlxDb().Beginx()
	if err != nil {
		return errors.Join(errors.New("Cannot start transaction"), err)
	}
	defer tx.Rollback()

	var currentMigration int
	var migrationDate string
	row := tx.QueryRow("SELECT version, date FROM migrations ORDER BY version DESC LIMIT 1;")
	err = row.Scan(&currentMigration, &migrationDate)
	if errors.Is(err, sql.ErrNoRows) {
		migrationDate = time.Now().Local().String()
		currentMigration = -1 // causes the migration slice to start from 0
	} else if err != nil {
		return errors.Join(errors.New("Cannot get current migration version"), err)
	}

	log.Info("Starting migrations from last version onwards", loggertags.TagDate, migrationDate, loggertags.TagVersion, currentMigration)
	for version, migration := range m.migrations[currentMigration+1:] {
		log.Info("Migrating database", loggertags.TagCurrent, version+currentMigration+1, loggertags.TagCount, len(m.migrations))
		if migration.PreProcess != nil {
			err = migration.PreProcess(tx)
			if err != nil {
				return errors.Join(errors.New("Cannot run pre-processing"), err)
			}
		}

		_, err = tx.Exec(migration.Sql)
		if err != nil {
			return errors.Join(errors.New("Cannot execute migration SQL"), err)
		}

		if migration.PostProcess != nil {
			err = migration.PostProcess(tx)
			if err != nil {
				return errors.Join(errors.New("Cannot run post-processing"), err)
			}
		}

		if migration.Test != nil {
			err = migration.Test(tx)
			if err != nil {
				return errors.Join(errors.New("Tests for the migration failed"), err)
			}
		}

		stmt := sqlx.Rebind(m.Rebinder, `
      INSERT INTO migrations (version, date)
      VALUES (?, ?);
    `)
		_, err = tx.Exec(stmt, version+currentMigration+1, time.Now())
		if err != nil {
			return errors.Join(errors.New("Cannot insert migration record"), err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(errors.New("Cannot commit migrations"), err)
	}

	return err
}
