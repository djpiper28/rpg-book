package backend

import (
	"errors"

	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/djpiper28/rpg-book/desktop_client/backend/database/migrations"
	"github.com/djpiper28/rpg-book/desktop_client/backend/model"
	"github.com/jmoiron/sqlx"
)

func Migrate(db *database.Db) error {
	m := migrations.New([]migrations.Migration{
		{
			Sql: `
  CREATE TABLE settings (
    dev_mode BOOLEAN NOT NULL DEFAULT('false'),
    dark_mode BOOLEAN NOT NULL DEFAULT('true')
  );
      `,
			PostProcess: func(tx *sqlx.Tx) error {
				_, err := tx.Exec("INSERT INTO settings DEFAULT VALUES;")
				return err
			},
			Test: func(tx *sqlx.Tx) error {
				rows := tx.QueryRowx("SELECT * FROM settings;")

				var settings model.Settings
				err := rows.StructScan(&settings)
				return err
			},
		},
		{
			Sql: `
  CREATE TABLE recently_opened (
    file_name TEXT PRIMARY KEY,
    project_name TEXT NOT NULL,
    last_opened TIMESTAMP WITH TIME ZONE NOT NULL
  );
      `,
		},
	})

	err := m.Migrate(db)
	if err != nil {
		return errors.Join(errors.New("Cannot migrate tables"), err)
	}

	return nil
}
