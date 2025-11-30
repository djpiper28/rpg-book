package model

import (
	"github.com/google/uuid"
)

type Note struct {
	Id       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Markdown string    `db:"markdown"`
	Created  string    `db:"created"`
}

type NoteRelations struct {
	NoteId      uuid.UUID `db:"note_id"`
	CharacterId uuid.UUID `db:"character_id"`
}
