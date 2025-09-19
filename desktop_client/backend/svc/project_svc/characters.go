package projectsvc

import (
	"bytes"
	"context"
	"errors"

	"github.com/charmbracelet/log"
	"github.com/disintegration/imaging"
	imagecompression "github.com/djpiper28/rpg-book/common/image/image_compression"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project_character"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project/model"
	"github.com/google/uuid"
)

func (p *ProjectSvc) compressImageIfNeeded(img []byte) ([]byte, error) {
	compressError := errors.New("Cannot compress iamge")

	settings, err := p.GetSettings()
	if err != nil {
		return nil, errors.Join(compressError, err)
	}

	if settings.CompressImages {
		img, err := imaging.Decode(bytes.NewBuffer(img), imaging.AutoOrientation(true))
		if err != nil {
			return nil, errors.Join(compressError, err)
		}

		compressedBytes, err := imagecompression.Compress(img)
		if err != nil {
			return nil, errors.Join(compressError, err)
		}

		return compressedBytes, nil
	} else {
		return img, nil
	}
}

func (p *ProjectSvc) CreateCharacter(ctx context.Context, in *pb_project.CreateCharacterReq) (*pb_project_character.CharacterHandle, error) {
	project, err := p.getProject(in.Project)

	img, err := p.compressImageIfNeeded(in.Details.Icon)
	if err != nil {
		log.Error("Cannot compress image", loggertags.TagError, err)
		return nil, errors.Join(errors.New("Cannot create character"), err)
	}
	in.Details.Icon = img

	character, err := project.CreateCharacter(in.Details.Name, in.Details.Description, in.Details.Icon)
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

	if in.SetImage {
		img, err := p.compressImageIfNeeded(in.Details.Icon)
		if err != nil {
			log.Error("Cannot compress image", loggertags.TagError, err)
			return nil, errors.Join(errors.New("Cannot create character"), err)
		}
		in.Details.Icon = img
	}

	err = project.UpdateCharacter(&model.Character{
		Id:          characterId,
		Description: in.Details.Description,
		Name:        in.Details.Name,
		Icon:        in.Details.Icon,
	}, in.SetImage)
	if err != nil {
		log.Error("Cannot update character", loggertags.TagError, err)
		return nil, errors.Join(updateError, err)
	}

	return &pb_common.Empty{}, nil
}

func (p *ProjectSvc) GetCharacter(ctx context.Context, in *pb_project.GetCharacterReq) (*pb_project_character.BasicCharacterDetails, error) {
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

	return character.ToPb(), nil
}
