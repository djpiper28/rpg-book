package project

import (
	"errors"
	"fmt"

	"github.com/djpiper28/rpg-book/common/database"
	"github.com/djpiper28/rpg-book/common/database/migrations"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project/model"
	"github.com/jmoiron/sqlx"
)

func Migrate(db database.Database) error {
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
		{
			Sql: `
  ALTER TABLE characters ADD COLUMN name_normalised TEXT NOT NULL DEFAULT('');
  ALTER TABLE characters ADD COLUMN description_normalised TEXT NOT NULL DEFAULT('');
      `,
			PostProcess: func(tx *sqlx.Tx) error {
				rows, err := tx.Queryx(`SELECT * FROM characters;`)
				if err != nil {
					return errors.Join(errors.New("Cannot query characters"), err)
				}

				characters := make([]*model.Character, 0)

				for rows.Next() {
					var character model.Character
					err := rows.StructScan(&character)
					if err != nil {
						return errors.Join(errors.New("Cannot scan characters"), err)
					}

					characters = append(characters, &character)
				}

				for _, character := range characters {
					character.Normalise()

					_, err = tx.Exec("UPDATE characters SET name_normalised=?, description_normalised=? WHERE id=?;",
						character.NameNormalised,
						character.DescriptionNormalised,
						character.Id)
					if err != nil {
						return errors.Join(fmt.Errorf("Cannot update character %s (id %s)", character.Name, character.Id), err)
					}
				}

				return nil
			},
		},
	})

	err := m.Migrate(db)
	if err != nil {
		return errors.Join(errors.New("Cannot execute migrations on project database"), err)
	}

	return nil
}
