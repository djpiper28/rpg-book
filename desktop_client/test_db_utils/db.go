package testdbutils

import (
	"os"

	"github.com/djpiper28/rpg-book/common/database/sqlite3"
	"github.com/djpiper28/rpg-book/desktop_client/backend"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project"
	"github.com/google/uuid"
)

func GetPrimaryDb() (*sqlite3.Db, func()) {
	name := uuid.New().String() + sqlite3.DbExtension
	db, err := sqlite3.New(name)
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

func GetProjectDb() (*sqlite3.Db, func()) {
	name := uuid.New().String() + sqlite3.DbExtension
	db, err := sqlite3.New(name)
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
