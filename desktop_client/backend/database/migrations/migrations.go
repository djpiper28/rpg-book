package migrations

import (
	"database/sql"
	"errors"

	"github.com/charmbracelet/log"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
)

type Migration struct {
	PreProcess  func() error
	Sql         string
	PostProcess func() error
	Test        func(*sql.Tx) error // Don't edit the database in tests plx
}

type DbMigrator struct {
	migrations []Migration
}

func New(migrations []Migration) *DbMigrator {
	return &DbMigrator{migrations: migrations}
}

func (m *DbMigrator) Migrate(db *database.Db) error {
	tx, err := db.Db.Begin()
	if err != nil {
		return errors.Join(errors.New("Cannot start transaction"), err)
	}

	for i, migration := range m.migrations {
		log.Info("Migrating database", loggertags.TagCurrent, i+1, loggertags.TagCount, len(m.migrations))
		if migration.PreProcess != nil {
			err = migration.PreProcess()
			if err != nil {
				return errors.Join(errors.New("Cannot run pre-processing"), err)
			}
		}

		_, err = tx.Exec(migration.Sql)
		if err != nil {
			return errors.Join(errors.New("Cannot execute migration SQL"), err)
		}

		if migration.PostProcess != nil {
			err = migration.PostProcess()
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
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(errors.New("Cannot commit migrations"), err)
	}

	return err
}
