package model

import (
	"github.com/djpiper28/rpg-book/common/normalisation"
	"github.com/google/uuid"
)

type Note struct {
	Id                 uuid.UUID `db:"id"`
	Name               string    `db:"name"`
	NameNormalised     string    `db:"name_normalised"`
	Markdown           string    `db:"markdown"`
	MarkdownNormalised string    `db:"markdown_normalised"`
	Created            string    `db:"created"`
}

func (n *Note) Normalise() {
	n.NameNormalised = normalisation.Normalise(n.Name)
	n.MarkdownNormalised = normalisation.Normalise(n.Markdown)
}

type NoteRelations struct {
	NoteId      uuid.UUID `db:"note_id"`
	CharacterId uuid.UUID `db:"character_id"`
}
