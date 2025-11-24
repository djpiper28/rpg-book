package project

import (
	"errors"
	"os"

	"github.com/charmbracelet/log"
	"github.com/djpiper28/rpg-book/common/database/sqlite3"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
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
