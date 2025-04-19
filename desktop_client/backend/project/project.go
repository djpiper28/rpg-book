package project

import (
	"errors"
	"os"

	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project/model"
)

type Project struct {
	db       *database.Db
	Settings model.ProjectSettings
	Filename string
}

func openDatabase(filename string) (*database.Db, error) {
	db, err := database.New(filename)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot open project database"), err)
	}

	err = Migrate(db)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot migrate project database"), err)
	}

	return db, err
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
	// TODO: cahnge to full struct scan when there are more than one fields
	err = rows.Scan(&settings.Name)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get project settings"), err)
	}

	return &Project{
		db:       db,
		Settings: settings,
		Filename: filename,
	}, nil
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

	return &Project{
		db: db,
		Settings: model.ProjectSettings{
			Name: projectName,
		},
		Filename: filename,
	}, nil
}

func (p *Project) Close() {
	defer p.db.Close()
}
