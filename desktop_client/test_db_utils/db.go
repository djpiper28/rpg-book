package testdbutils

import (
	"os"

	"github.com/djpiper28/rpg-book/desktop_client/backend"
	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project"
	"github.com/google/uuid"
)

func GetPrimaryDb() (*database.Db, func()) {
	name := uuid.New().String() + database.DbExtension
	db, err := database.New(name)
	if err != nil {
		panic(err)
	}

	err = backend.Migrate(db)
	if err != nil {
		panic(err)
	}

	return db, func() {
		os.Remove(name)
	}
}

func GetProjectDb() (*database.Db, func()) {
	name := uuid.New().String() + database.DbExtension
	db, err := database.New(name)
	if err != nil {
		panic(err)
	}

	err = project.Migrate(db)
	if err != nil {
		panic(err)
	}

	return db, func() {
		os.Remove(name)
	}
}
