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
    description TEXT NOT NULL,
		PRIMARY KEY(id)
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
		{
			/*
				When the note_relations TABLE is updated its constraint should be updated to make its CHECK looks like below

				CREATE TABLE note_relations (
				    note_id TEXT NOT NULL,
				    id_1 INTEGER,
				    id_2 INTEGER,
				    id_3 INTEGER,

				    CHECK (
				        (id_1 IS NOT NULL) +
				        (id_2 IS NOT NULL) +
				        (id_3 IS NOT NULL) <= 1
				    )
				);
			*/
			Sql: `
	CREATE TABLE notes (
		id TEXT NOT NULL,
		name TEXT NOT NULL,
		name_normalised TEXT NOT NULL,
		markdown TEXT NOT NULL,
		markdown_normalised TEXT NOT NULL,
    created TIMESTAMP WITH TIME ZONE NOT NULL,
		PRIMARY KEY(id)
	);

	CREATE TABLE note_relations (
		note_id TEXT NOT NULL,
		character_id TEXT,
		FOREIGN KEY (character_id) REFERENCES characters(id),
		FOREIGN KEY (note_id) REFERENCES notes(id),
		UNIQUE(note_id, character_id),
		CHECK (
      (character_id IS NOT NULL) <= 1
    )
	);
			`,
		},
	})

	err := m.Migrate(db)
	if err != nil {
		return errors.Join(errors.New("Cannot execute migrations on project database"), err)
	}

	return nil
}
