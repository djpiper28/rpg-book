package project

import (
	"errors"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/djpiper28/rpg-book/common/database/sqlite3"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/djpiper28/rpg-book/common/search/parser"
	sqlsearch "github.com/djpiper28/rpg-book/common/search/sql_search"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project/model"
	"github.com/google/uuid"
)

type Project struct {
	db       *sqlite3.Db
	Settings model.ProjectSettings
	Filename string
}

func Open(filename string) (*Project, error) {
	_, err := os.Stat(filename)
	if err != nil {
		return nil, errors.New("Project does not exist so cannot be opned")
	}

	db, err := openDatabase(filename)
	if err != nil {
		return nil, err
	}

	var settings model.ProjectSettings
	rows := db.Db.QueryRowx("SELECT * FROM project_settings;")
	// TODO: change to full struct scan when there are more than one fields
	err = rows.Scan(&settings.Name)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get project settings"), err)
	}

	log.Info("Opened a project", loggertags.TagFileName, filename)
	return &Project{
		db:       db,
		Settings: settings,
		Filename: filename,
	}, nil
}

func openDatabase(filename string) (*sqlite3.Db, error) {
	db, err := sqlite3.New(filename)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot open project database"), err)
	}

	err = Migrate(db)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot migrate project database"), err)
	}

	return db, nil
}

func Create(filename string, projectName string) (*Project, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return nil, errors.New("Cannot override a project as it is dangerous")
	}

	if !errors.Is(err, os.ErrNotExist) {
		return nil, errors.Join(errors.New("Expected the file to not exist"), err)
	}

	db, err := openDatabase(filename)
	if err != nil {
		return nil, err
	}

	_, err = db.Db.Exec("INSERT INTO project_settings(name) VALUES(?);", projectName)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot create project settings"), err)
	}

	log.Info("Created a project", loggertags.TagFileName, filename)
	return &Project{
		db: db,
		Settings: model.ProjectSettings{
			Name: projectName,
		},
		Filename: filename,
	}, nil
}

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

	row := p.db.Db.QueryRowx(`SELECT * FROM characters WHERE id=?;`, id)
	err := row.StructScan(&character)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get character"), err)
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

func (p *Project) Close() {
	defer p.db.Close()
}

func (p *Project) SearchCharacter(query string) ([]uuid.UUID, error) {
	ast, err := parser.Parse(query)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot parse query"), err)
	}

	const (
		description = "description"
		name        = "name"
	)

	sql, args, err := sqlsearch.AsSql(ast,
		sqlsearch.SqlTableData{
			FieldsToScan: []string{"id"},
			TableName:    "characters",
		},
		sqlsearch.SqlColmnMap{
			TextColumns: map[string]string{
				"name":        name,
				"desc":        description,
				"description": description,
			},
			BasicQueryColumn: name,
		})
	if err != nil {
		return nil, errors.Join(errors.New("Cannot process query"), err)
	}

	log.Debug("Executing search", "sql", sql, "args", args)
	rows, err := p.db.Db.Queryx(sql, args...)
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

	return characterIds, nil
}

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
