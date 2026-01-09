package project

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/charmbracelet/log"
	"github.com/djpiper28/rpg-book/common/search/parser"
	sqlsearch "github.com/djpiper28/rpg-book/common/search/sql_search"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project/model"
	"github.com/google/uuid"
)

func (p *Project) CreateCharacter(name, description string, icon []byte) (*model.Character, error) {
	character := model.NewCharacter(name)
	character.Description = description
	character.Icon = icon
	character.Created = time.Now().String()
	character.Normalise()

	_, err := p.db.Db.NamedExec(`
    INSERT INTO 
    characters (id, name, name_normalised, created, description, description_normalised, icon)
    VALUES(:id, :name, :name_normalised, :created, :description, :description_normalised, :icon);`,
		character)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot create character"), err)
	}

	return character, nil
}

func (p *Project) GetCharacters() ([]*model.Character, error) {
	// TODO: do not select icon
	rows, err := p.db.Db.Queryx(`SELECT * FROM characters;`)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot query characters"), err)
	}

	characters := make([]*model.Character, 0)

	for rows.Next() {
		var character model.Character
		err := rows.StructScan(&character)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot scan characters"), err)
		}

		characters = append(characters, &character)
	}

	return characters, nil
}

func (p *Project) GetCharacter(id uuid.UUID) (*model.Character, error) {
	character := model.Character{}

	tx, err := p.db.Db.Beginx()
	if err != nil {
		return nil, errors.Join(errors.New("Cannot begin transaction"), err)
	}
	defer tx.Rollback()

	row := tx.QueryRowx(`SELECT * FROM characters WHERE id=?;`, id)
	err = row.StructScan(&character)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get character"), err)
	}

	// Select the related notes
	notes := make([]*model.Note, 0)
	rows, err := tx.Queryx(`
    SELECT notes.*
    FROM
    (
      (
        characters
        INNER JOIN note_relations
          ON note_relations.character_id = characters.id
      )
      INNER JOIN notes
        ON notes.id = note_relations.note_id
    )
    WHERE characters.id = ?;
    `, id)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get related notes"), err)
	}

	for rows.Next() {
		var note model.Note
		err = rows.StructScan(&note)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot scan notes"), err)
		}

		notes = append(notes, &note)
	}

	character.Notes = notes

	err = tx.Commit()
	if err != nil {
		return nil, errors.Join(errors.New("Cannot commit transaction"), err)
	}

	return &character, nil
}

func (p *Project) UpdateCharacter(character *model.Character, setIcon bool) error {
	character.Normalise()

	sql := `
    UPDATE characters
      SET icon=:icon, name=:name, description=:description,
          name_normalised=:name_normalised, description_normalised=:description_normalised
      WHERE id=:id;
    `
	if !setIcon {
		sql = `
    UPDATE characters
      SET name=:name, description=:description,
          name_normalised=:name_normalised, description_normalised=:description_normalised
      WHERE id=:id;
    `
	}

	_, err := p.db.Db.NamedExec(sql, character)
	if err != nil {
		return errors.Join(errors.New("Cannot update character"), err)
	}

	return nil
}

func (p *Project) DeleteCharacter(id uuid.UUID) error {
	_, err := p.db.Db.Exec("DELETE FROM characters WHERE id=?;", id)
	if err != nil {
		return errors.Join(errors.New("Cannot delete character"), err)
	}

	return nil
}

const (
	characterDescription = "characters.description"
	characterName        = "characters.name"
)

func characterColumns() map[string]string {
	return map[string]string{
		"name":        characterName,
		"desc":        characterDescription,
		"description": characterDescription,
	}
}

func qualifiedCharacterColumns() map[string]string {
	unmodifiedColumns := characterColumns()
	modifiedColumns := make(map[string]string)

	for key, value := range unmodifiedColumns {
		modifiedColumns["character."+key] = value
	}
	return modifiedColumns
}

func (p *Project) SearchCharacter(query string) ([]uuid.UUID, error) {
	ast, err := parser.Parse(query)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot parse query"), err)
	}

	searchSql, args, err := sqlsearch.AsSql(ast,
		sqlsearch.SqlTableData{
			FieldsToScan: []string{"characters.id"},
			TableName:    "characters",
			JoinClauses: `
  LEFT JOIN note_relations ON note_relations.character_id = characters.id
  LEFT JOIN notes ON notes.id = note_relations.note_id
      `,
		},
		sqlsearch.SqlColmnMap{
			TextColumns:      mergeColumns(characterColumns(), qualifiedCharacterColumns(), qualifiedNoteColumns()),
			BasicQueryColumn: "name",
		})
	if err != nil {
		return nil, errors.Join(errors.New("Cannot process query"), err)
	}

	log.Debug("Executing search", "sql", searchSql, "args", args)
	tx, err := p.db.Db.BeginTxx(context.Background(), &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Join(errors.New("Cannot start read only SQL transaction"), err)
	}

	defer tx.Rollback()

	rows, err := tx.Queryx(searchSql, args...)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot execute SQL query"), err)
	}

	defer rows.Close()
	characterIds := make([]uuid.UUID, 0)
	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(&id)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot read characters"), err)
		}

		characterIds = append(characterIds, id)
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Join(errors.New("Cannot commit transaction"), err)
	}

	return characterIds, nil
}
