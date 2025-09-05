package projectsvc

import (
	"errors"

	"github.com/charmbracelet/log"
	"github.com/djpiper28/rpg-book/desktop_client/backend/model"
)

func (p *ProjectSvc) GetSettings() (model.Settings, error) {
	row := p.primaryDb.Db.QueryRowx("SELECT * FROM settings;")
	if row == nil {
		log.Error("Cannot get settings")
		return model.Settings{}, errors.New("Rows are nil")
	}

	var settings model.Settings
	err := row.StructScan(&settings)
	if err != nil {
		log.Error("Cannot get settings")
		return model.Settings{}, errors.Join(errors.New("Cannot get settings"), err)
	}

	return settings, nil
}
