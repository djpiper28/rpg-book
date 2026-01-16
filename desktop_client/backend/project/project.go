package project

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/djpiper28/rpg-book/common/database/sqlite3"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project/model"
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

func (p *Project) Close() {
	defer p.db.Close()
}

func mergeColumns(columnMaps ...map[string]string) map[string]string {
	output := make(map[string]string)

	for _, columnMap := range columnMaps {
		for key, value := range columnMap {
			_, found := output[key]
			if found {
				panic(fmt.Sprintf("Illegal merge operation - duplicate key %s", key))
			}

			output[key] = value
		}
	}

	return output
}
