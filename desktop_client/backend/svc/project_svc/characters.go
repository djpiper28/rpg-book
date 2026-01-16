package projectsvc

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"os"

	"github.com/charmbracelet/log"
	"github.com/djpiper28/rpg-book/common/image"
	imagecompression "github.com/djpiper28/rpg-book/common/image/image_compression"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project_character"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project_note"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project/model"
	"github.com/google/uuid"
)

func (p *ProjectSvc) readAndCompressImageIfNeeded(path string) ([]byte, error) {
	if path == "" {
		return []byte{}, nil
	}

	compressError := errors.New("Cannot compress iamge")

	imgBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Join(compressError, err)
	}

	settings, err := p.GetSettings()
	if err != nil {
		return nil, errors.Join(compressError, err)
	}

	if settings.CompressImages {
		img, err := image.CustomDecode(bytes.NewBuffer(imgBytes))
		if err != nil {
			return nil, errors.Join(compressError, err)
		}

		compressedBytes, err := imagecompression.CompressIcon(img)
		if err != nil {
			return nil, errors.Join(compressError, err)
		}

		return compressedBytes, nil
	} else {
		return imgBytes, nil
	}
}

func (p *ProjectSvc) CreateCharacter(ctx context.Context, in *pb_project.CreateCharacterReq) (*pb_project_character.CharacterHandle, error) {
	project, err := p.getProject(in.Project)
	if err != nil {
		log.Error("Cannot get projecct id", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Cannot get proejct"), err)
	}

	var img []byte
	if in.IconPath != "" {
		img, err = p.readAndCompressImageIfNeeded(in.IconPath)
		if err != nil {
			log.Error("Cannot compress image", loggertags.TagError, err)
			return nil, errors.Join(errors.New("Cannot create character"), err)
		}
	}

	character, err := project.CreateCharacter(in.Details.Name, in.Details.Description, img)
	if err != nil {
		log.Error("Cannot create character", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Cannot create character"), err)
	}

	return &pb_project_character.CharacterHandle{
		Id: character.Id.String(),
	}, nil
}

func (p *ProjectSvc) UpdateCharacter(ctx context.Context, in *pb_project.UpdateCharacterReq) (*pb_common.Empty, error) {
	updateError := errors.New("Cannot update character")

	project, err := p.getProject(in.Project)
	if err != nil {
		log.Error("Cannot get projecct id", loggertags.TagError, err)
		return nil, updateError
	}

	characterId, err := uuid.Parse(in.Handle.Id)
	if err != nil {
		log.Error("Cannot get character id", loggertags.TagError, err)
		return nil, errors.Join(updateError, err)
	}

	img, err := p.readAndCompressImageIfNeeded(in.IconPath)
	if err != nil {
		log.Error("Cannot compress image", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Cannot update character"), err)
	}

	err = project.UpdateCharacter(&model.Character{
		Id:          characterId,
		Description: in.Details.Description,
		Name:        in.Details.Name,
		Icon:        img,
	}, in.SetIcon)
	if err != nil {
		log.Error("Cannot update character", loggertags.TagError, err)
		return nil, errors.Join(updateError, err)
	}

	return &pb_common.Empty{}, nil
}

func (p *ProjectSvc) GetCharacter(ctx context.Context, in *pb_project.GetCharacterReq) (*pb_project.GetCharacterResp, error) {
	project, err := p.getProject(in.Project)
	if err != nil {
		log.Error("Cannot get character", loggertags.TagError, err)
		return nil, errors.New("Cannot find project")
	}

	characterId, err := uuid.Parse(in.Character.Id)
	if err != nil {
		log.Error("Cannot get character", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Invalid character ID"), err)
	}

	character, err := project.GetCharacter(characterId)
	if err != nil {
		log.Error("Cannot get character", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Cannot get character"), err)
	}

	notes := make([]*pb_project_note.Note, 0)
	for _, note := range character.Notes {
		notes = append(notes, note.ToPb())
	}

	joinedUrl, err := url.JoinPath(p.baseUrl, "image", "character", in.Project.Id, in.Character.Id)
	if err != nil {
		log.Error("Cannot create icon url", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Cannot create icon url"), err)
	}

	return &pb_project.GetCharacterResp{
		Details: character.ToPb(),
		IconUrl: joinedUrl,
		Notes:   notes,
	}, nil
}

func (p *ProjectSvc) GetCharacterImage(projectID, characterID uuid.UUID) ([]byte, error) {
	project, err := p.getProject(&pb_project.ProjectHandle{Id: projectID.String()})
	if err != nil {
		return nil, err
	}

	character, err := project.GetCharacter(characterID)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get character"), err)
	}

	return character.Icon, nil
}

func (p *ProjectSvc) DeleteCharacter(ctx context.Context, in *pb_project.DeleteCharacterReq) (*pb_common.Empty, error) {
	project, err := p.getProject(in.Project)
	if err != nil {
		log.Error("Cannot get character", loggertags.TagError, err)
		return nil, errors.New("Cannot find project")
	}

	characterId, err := uuid.Parse(in.Handle.Id)
	if err != nil {
		log.Error("Cannot get character", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Invalid character ID"), err)
	}

	err = project.DeleteCharacter(characterId)
	if err != nil {
		log.Error("Cannot delete character", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Cannotdelete character"), err)
	}

	return &pb_common.Empty{}, nil
}

func (p *ProjectSvc) SearchCharacter(ctx context.Context, in *pb_project.QueryReq) (*pb_project.SearchCharacterResp, error) {
	project, err := p.getProject(in.Project)
	if err != nil {
		log.Error("Cannot search character", loggertags.TagError, err)
		return nil, errors.New("Cannot find project")
	}

	characterIds, err := project.SearchCharacter(in.Query)
	if err != nil {
		log.Error("Cannot search character", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Cannot perform search"), err)
	}

	resp := &pb_project.SearchCharacterResp{
		Details: make([]*pb_project_character.CharacterHandle, len(characterIds)),
	}

	for i, characterId := range characterIds {
		resp.Details[i] = &pb_project_character.CharacterHandle{
			Id: characterId.String(),
		}
	}

	return resp, nil
}
