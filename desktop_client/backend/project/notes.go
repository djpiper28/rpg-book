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

func (p *Project) CreateNote(name, markdown string, characterIds []uuid.UUID) (*model.Note, error) {
	note := &model.Note{
		Id:       uuid.New(),
		Name:     name,
		Markdown: markdown,
		Created:  time.Now().String(),
	}

	note.Normalise()

	tx, err := p.db.Db.Beginx()
	if err != nil {
		return nil, errors.Join(errors.New("Cannot start transaction"), err)
	}
	defer tx.Rollback()

	_, err = tx.NamedExec(`
    INSERT INTO notes 
      (id, name, name_normalised, markdown, markdown_normalised, created) 
    VALUES 
      (:id, :name, :name_normalised, :markdown, :markdown_normalised, :created);`, note)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot insert note"), err)
	}

	for _, id := range characterIds {
		_, err = tx.Exec(`
      INSERT INTO note_relations
        (note_id, character_id)
      VALUES
        (?, ?);
      `, note.Id, id)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot insert note relation"), err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Join(errors.New("Cannot commit transaction"), err)
	}

	return note, nil
}

func (p *Project) GetNotes() ([]*model.Note, error) {
	notes := make([]*model.Note, 0)

	rows, err := p.db.Db.Queryx("SELECT * FROM notes;")
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get notes"), err)
	}

	for rows.Next() {
		var note model.Note
		err = rows.StructScan(&note)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot scan note into struct"), err)
		}

		notes = append(notes, &note)
	}

	return notes, nil
}

func (p *Project) GetNote(noteId uuid.UUID) (*model.CompleteNote, error) {
	tx, err := p.db.Db.BeginTxx(context.Background(), &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Join(errors.New("Cannot start transaction"), err)
	}
	defer tx.Rollback()

	row := tx.QueryRowx("SELECT * FROM notes WHERE id=?;", noteId)

	var note model.Note
	err = row.StructScan(&note)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot scan note into struct"), err)
	}

	rows, err := tx.Queryx(`
    SELECT characters.*
    FROM characters 
    INNER JOIN note_relations 
      ON characters.id=note_relations.character_id
    WHERE note_id=?;`, noteId)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get related characters"), err)
	}

	characters := make([]*model.Character, 0)
	for rows.Next() {
		var character model.Character
		err = rows.StructScan(&character)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot scan character into struct"), err)
		}

		characters = append(characters, &character)
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Join(errors.New("Cannot commit transaction"), err)
	}

	return &model.CompleteNote{
		Note:       &note,
		Characters: characters,
	}, nil
}

const (
	noteMarkdown = "notes.markdown"
	noteName     = "notes.name"
)

func noteColumns() map[string]string {
	return map[string]string{
		"name":        noteName,
		"desc":        noteMarkdown,
		"description": noteMarkdown,
		"markdown":    noteMarkdown,
		"contents":    noteMarkdown,
	}
}
func qualifiedNoteColumns() map[string]string {
	unmodifiedColumns := noteColumns()
	modifiedColumns := make(map[string]string)

	for key, value := range unmodifiedColumns {
		modifiedColumns["note."+key] = value
	}
	return modifiedColumns
}

func (p *Project) SearchNote(query string) ([]uuid.UUID, error) {
	ast, err := parser.Parse(query)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot parse query"), err)
	}

	searchSql, args, err := sqlsearch.AsSql(ast,
		sqlsearch.SqlTableData{
			FieldsToScan: []string{"notes.id"},
			TableName:    "notes",
      JoinClauses: `
  LEFT JOIN note_relations ON note_relations.note_id = notes.id
  LEFT JOIN characters ON characters.id = note_relations.character_id
      `,
		},
		sqlsearch.SqlColmnMap{
			TextColumns:      mergeColumns(noteColumns(), qualifiedNoteColumns(), qualifiedCharacterColumns()),
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

	rows, err := tx.Queryx(searchSql, args...)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot execute SQL query"), err)
	}

	defer rows.Close()
	noteIds := make([]uuid.UUID, 0)
	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(&id)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot read notes"), err)
		}

		noteIds = append(noteIds, id)
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Join(errors.New("Cannot commit transaction"), err)
	}

	return noteIds, nil
}
