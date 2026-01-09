package projectsvc

import (
	"context"
	"errors"

	"github.com/charmbracelet/log"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project_character"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project_note"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project/model"
	"github.com/google/uuid"
)

func (p *ProjectSvc) CreateNote(ctx context.Context, in *pb_project.CreateNoteReq) (*pb_project_note.NoteHandle, error) {
	project, err := p.getProject(in.Project)
	if err != nil {
		log.Error("Cannot get projecct id", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Cannot get proejct"), err)
	}

	characterIds := make([]uuid.UUID, 0)
	for _, characterHandle := range in.Characters {
		id, err := uuid.Parse(characterHandle.Id)
		if err != nil {
			log.Error("Cannot parse uuid", loggertags.TagError, err)
			return nil, errors.Join(errors.New("Cannot parse uuid"), err)
		}

		characterIds = append(characterIds, id)
	}

	note, err := project.CreateNote(in.Details.Name, in.Details.Markdown, characterIds)
	if err != nil {
		log.Error("Cannot create note", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Cannot get proejct"), err)
	}

	return &pb_project_note.NoteHandle{
		Id: note.Id.String(),
	}, nil
}

func (p *ProjectSvc) GetNote(ctx context.Context, in *pb_project.GetNoteReq) (*pb_project.GetNoteResp, error) {
	project, err := p.getProject(in.Project)
	if err != nil {
		log.Error("Cannot get project from id", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Cannot get project"), err)
	}

	id, err := uuid.Parse(in.Note.Id)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot parse note ID"), err)
	}

	completeNote, err := project.GetNote(id)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get note"), err)
	}

	result := &pb_project.GetNoteResp{
		Details:    completeNote.Note.ToPb(),
		Characters: make([]*pb_project_character.BasicCharacterDetails, 0),
	}

	for _, character := range completeNote.Characters {
		result.Characters = append(result.Characters, character.ToPb())
	}

	return result, nil
}

func (p *ProjectSvc) SearchNote(ctx context.Context, in *pb_project.QueryReq) (*pb_project.SearchNoteResp, error) {
	project, err := p.getProject(in.Project)
	if err != nil {
		log.Error("Cannot search note", loggertags.TagError, err)
		return nil, errors.New("Cannot find project")
	}

	noteIds, err := project.SearchNote(in.Query)
	if err != nil {
		log.Error("Cannot search note", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Cannot perform search"), err)
	}

	resp := &pb_project.SearchNoteResp{
		Details: make([]*pb_project_note.NoteHandle, len(noteIds)),
	}

	for i, noteId := range noteIds {
		resp.Details[i] = &pb_project_note.NoteHandle{
			Id: noteId.String(),
		}
	}

	return resp, nil
}

func (p *ProjectSvc) UpdateNote(ctx context.Context, in *pb_project.UpdateNoteReq) (*pb_common.Empty, error) {
	project, err := p.getProject(in.Project)
	if err != nil {
		log.Error("Cannot get project from id", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Cannot get project"), err)
	}

	noteId, err := uuid.Parse(in.Handle.Id)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot parse note ID"), err)
	}

	note := &model.Note{
		Id:       noteId,
		Name:     in.Details.Name,
		Markdown: in.Details.Markdown,
	}

	characterIds := make([]uuid.UUID, 0)
	for _, character := range in.Characters {
		characterId, err := uuid.Parse(character.Id)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot parse character ID"), err)
		}

		characterIds = append(characterIds, characterId)
	}

	err = project.UpdateNote(note, characterIds)
	if err != nil {
		log.Error("Cannot update note", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Cannot update note"), err)
	}

	return nil, nil
}
