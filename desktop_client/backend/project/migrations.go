package project

import (
	"errors"

	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/djpiper28/rpg-book/desktop_client/backend/database/migrations"
)

func Migrate(db *database.Db) error {
	m := migrations.New([]migrations.Migration{
		{
			Sql: `
  CREATE TABLE project_settings (
    name TEXT NOT NULL
  );

  CREATE TABLE characters (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL,
    icon BLOB,
    description TEXT NOT NULL
  );
      `,
		},
	})

	err := m.Migrate(db)
	if err != nil {
		return errors.Join(errors.New("Cannot migrate project database"), err)
	}

	return nil
}
