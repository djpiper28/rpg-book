package testdbutils

import (
	"os"

	"github.com/djpiper28/rpg-book/desktop_client/backend"
	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/google/uuid"
)

func GetDb() (*database.Db, func()) {
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
