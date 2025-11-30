package projectsvc

import (
	"context"
	"errors"

	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project_note"
)

func (p *ProjectSvc) CreateNote(ctx context.Context, in *pb_project_note.CreateNoteReq) (*pb_project_note.NoteHandle, error) {
	return nil, errors.ErrUnsupported
}
