package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectSettings struct {
	Name string `db:"name"`
}

type Project struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Created     string    `db:"created"`
	Icon        []byte    `db:"icon"`
	Description string    `db:"Description"`
}

func NewProject(name string) *Project {
	return &Project{
		Id:      uuid.New(),
		Name:    name,
		Created: time.Now().String(),
	}
}
