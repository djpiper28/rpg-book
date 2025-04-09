package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	Db       *sqlx.DB
	filename string
}

const (
	DbTimeout   = time.Second * 10
	DbExtension = ".sqlite"
)

func New(filename string) (*Db, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DbTimeout)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "sqlite3", fmt.Sprintf("file:%s?cache=shared", filename))
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Cannot create database connetion for %s", filename), err)
	}

	log.Info("Connected to new database", loggertags.TagFileName, filename)
	db.SetMaxOpenConns(1)
	return &Db{
		Db:       db,
		filename: filename,
	}, nil
}

func (d *Db) Close() {
	d.Db.Close()
	log.Info("Disconnected from database", loggertags.TagFileName, d.filename)
}
