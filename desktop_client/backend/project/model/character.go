package model

import (
	"time"

	"github.com/djpiper28/rpg-book/common/normalisation"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project_character"
	"github.com/google/uuid"
)

type Character struct {
	Id                    uuid.UUID `db:"id"`
	Name                  string    `db:"name"`
	NameNormalised        string    `db:"name_normalised"`
	Created               string    `db:"created"`
	Icon                  []byte    `db:"icon"`
	Description           string    `db:"description"`
	DescriptionNormalised string    `db:"description_normalised"`
}

func NewCharacter(name string) *Character {
  c := &Character{
		Id:      uuid.New(),
		Name:    name,
		Created: time.Now().String(),
	}
  c.Normalise()
  return c
}

func (c *Character) Normalise() {
	c.NameNormalised = normalisation.Normalise(c.Name)
	c.DescriptionNormalised = normalisation.Normalise(c.Description)
}

// The icon is not set here, if you want an icon then you need to set it to another proto
// this is to avoid sending the icon as much
func (c *Character) ToPb() *pb_project_character.BasicCharacterDetails {
	return &pb_project_character.BasicCharacterDetails{
		Handle: &pb_project_character.CharacterHandle{
			Id: c.Id.String(),
		},
		Name:        c.Name,
		Description: c.Description,
	}
}
