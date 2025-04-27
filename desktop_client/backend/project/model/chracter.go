package model

import (
	"time"

	"github.com/google/uuid"
)

type Character struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Created     string    `db:"created"`
	Icon        []byte    `db:"icon"`
	Description string    `db:"description"`
}

func NewCharacter(name string) *Character {
	return &Character{
		Id:      uuid.New(),
		Name:    name,
		Created: time.Now().String(),
	}
}
