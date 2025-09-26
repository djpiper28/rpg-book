package projectsvc

import (
	"bytes"
	"context"
	"errors"
	"os"

	"github.com/charmbracelet/log"
	"github.com/djpiper28/rpg-book/common/image"
	imagecompression "github.com/djpiper28/rpg-book/common/image/image_compression"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project_character"
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
		log.Error("Cannot get projecct id")
		return nil, updateError
	}

	characterId, err := uuid.Parse(in.Handle.Id)
	if err != nil {
		log.Error("Cannot get character id")
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
		log.Error("Cannot get character")
		return nil, errors.New("Cannot find project")
	}

	characterId, err := uuid.Parse(in.Character.Id)
	if err != nil {
		log.Error("Cannot get character")
		return nil, errors.Join(errors.New("Invalid character ID"), err)
	}

	character, err := project.GetCharacter(characterId)
	if err != nil {
		log.Error("Cannot get character")
		return nil, errors.Join(errors.New("Cannot get character"), err)
	}

	return &pb_project.GetCharacterResp{
		Details: character.ToPb(),
		Icon:    character.Icon,
	}, nil
}
