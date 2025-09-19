package model

import (
	"time"

	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project_character"
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

func (c *Character) ToPb() *pb_project_character.BasicCharacterDetails {
	return &pb_project_character.BasicCharacterDetails{
		Handle: &pb_project_character.CharacterHandle{
			Id: c.Id.String(),
		},
		Name:        c.Name,
		Description: c.Description,
		Icon:        c.Icon,
	}
}
