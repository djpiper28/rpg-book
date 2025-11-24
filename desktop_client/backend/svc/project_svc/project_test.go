package projectsvc_test

import (
	"bytes"
	"context"
	"image/jpeg"
	"os"
	"testing"

	"github.com/djpiper28/rpg-book/common/database/sqlite3"
	testutils "github.com/djpiper28/rpg-book/common/test_utils"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project_character"
	projectsvc "github.com/djpiper28/rpg-book/desktop_client/backend/svc/project_svc"
	testdbutils "github.com/djpiper28/rpg-book/desktop_client/test_db_utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func newTestImageFile(t *testing.T) string {
	t.Helper()

	img := testutils.NewTestImage(100, 100)
	buff := bytes.NewBuffer([]byte{})
	err := jpeg.Encode(buff, img, nil)
	require.NoError(t, err)

	file, err := os.CreateTemp("", "test-image-*.jpg")
	require.NoError(t, err)
	defer file.Close()

	_, err = file.Write(buff.Bytes())
	require.NoError(t, err)

	return file.Name()
}

func TestProjectSvc(t *testing.T) {
	t.Parallel()
	db, closeDb := testdbutils.GetPrimaryDb()
	defer closeDb()

	svc := projectsvc.New(db)
	defer svc.Close()

	closeProjectWithoutDelete := func(t *testing.T, handle *pb_project.ProjectHandle) {
		_, err := svc.CloseProject(context.Background(), handle)
		require.NoError(t, err)
	}

	closeProject := func(t *testing.T, filename string, handle *pb_project.ProjectHandle) {
		defer os.Remove(filename)

		closeProjectWithoutDelete(t, handle)
	}

	t.Run("Test no recent projects", func(t *testing.T) {
		projects, err := svc.RecentProjects(context.Background(), &pb_common.Empty{})
		require.NoError(t, err)
		require.Len(t, projects.Projects, 0)
	})

	t.Run("Test open project that does not exist", func(t *testing.T) {
		const fileDoesNotExist = "does-not-exist"
		_, err := svc.OpenProject(context.Background(), &pb_project.OpenProjectReq{
			FileName: fileDoesNotExist,
		})

		require.Error(t, err)
		var count int
		rows := db.Db.QueryRow("SELECT COUNT(file_name) FROM recently_opened WHERE file_name=?;", fileDoesNotExist)
		err = rows.Scan(&count)

		require.NoError(t, err)
		require.Equal(t, 0, count)
	})

	t.Run("Test create new project", func(t *testing.T) {
		filename := uuid.New().String() + sqlite3.DbExtension
		name := uuid.New().String()

		handle, err := svc.CreateProject(context.Background(), &pb_project.CreateProjectReq{
			FileName:    filename,
			ProjectName: name,
		})
		require.NoError(t, err)
		defer closeProject(t, filename, handle)

		require.NotEmpty(t, handle)

		var count int
		var dbName string
		rows := db.Db.QueryRow("SELECT COUNT(file_name), project_name FROM recently_opened WHERE file_name=?;", filename)
		err = rows.Scan(&count, &dbName)

		require.NoError(t, err)
		require.Equal(t, 1, count)
		require.Equal(t, name, dbName)
	})

	t.Run("Open project after create", func(t *testing.T) {
		filename := uuid.New().String() + sqlite3.DbExtension
		name := uuid.New().String()

		handle, err := svc.CreateProject(context.Background(), &pb_project.CreateProjectReq{
			FileName:    filename,
			ProjectName: name,
		})
		require.NoError(t, err)
		closeProjectWithoutDelete(t, handle)

		resp, err := svc.OpenProject(context.Background(), &pb_project.OpenProjectReq{
			FileName: filename,
		})
		require.NoError(t, err)
		defer closeProject(t, filename, resp.Handle)

		var count int
		var dbName string
		rows := db.Db.QueryRow("SELECT COUNT(file_name), project_name FROM recently_opened WHERE file_name=?;", filename)
		err = rows.Scan(&count, &dbName)

		require.NoError(t, err)
		require.Equal(t, 1, count)
		require.Equal(t, name, dbName)
	})

	t.Run("Test open project that is not tracked", func(t *testing.T) {
		filename := uuid.New().String() + sqlite3.DbExtension
		name := uuid.New().String()

		handle, err := svc.CreateProject(context.Background(), &pb_project.CreateProjectReq{
			FileName:    filename,
			ProjectName: name,
		})
		require.NoError(t, err)
		closeProjectWithoutDelete(t, handle)

		require.NotEmpty(t, handle)

		var count int
		var dbName string
		rows := db.Db.QueryRow("SELECT COUNT(file_name), project_name FROM recently_opened WHERE file_name=?;", filename)
		err = rows.Scan(&count, &dbName)

		require.NoError(t, err)
		require.Equal(t, 1, count)
		require.Equal(t, name, dbName)

		os.Rename(filename, "new_"+filename)
		filename = "new_" + filename

		resp, err := svc.OpenProject(context.Background(), &pb_project.OpenProjectReq{
			FileName: filename,
		})
		require.NoError(t, err)
		defer closeProject(t, filename, resp.Handle)

		rows = db.Db.QueryRow("SELECT COUNT(file_name), project_name FROM recently_opened WHERE file_name=?;", filename)
		err = rows.Scan(&count, &dbName)

		require.NoError(t, err)
		require.Equal(t, 1, count)
		require.Equal(t, name, dbName)
	})

	t.Run("Test close projet that does not exist", func(t *testing.T) {
		_, err := svc.CloseProject(context.Background(), &pb_project.ProjectHandle{
			Id: uuid.NewString(),
		})

		require.Error(t, err)
	})

	t.Run("Test created project appears in recently openeed projects", func(t *testing.T) {
		filename := uuid.New().String() + sqlite3.DbExtension
		name := uuid.New().String()

		handle, err := svc.CreateProject(context.Background(), &pb_project.CreateProjectReq{
			FileName:    filename,
			ProjectName: name,
		})
		require.NoError(t, err)
		defer closeProject(t, filename, handle)

		recents, err := svc.RecentProjects(context.Background(), &pb_common.Empty{})
		require.NoError(t, err)

		for _, recent := range recents.Projects {
			if recent.FileName == filename {
				require.Equal(t, name, recent.ProjectName)
				return
			}
		}

		t.Log("Cannot find the project in the recent list")
		t.Fail()
	})

	t.Run("Test create character in project", func(t *testing.T) {
		filename := uuid.New().String() + sqlite3.DbExtension
		projectName := uuid.New().String()
		characterName := uuid.New().String()
		characterDescription := uuid.New().String()
		iconPath := newTestImageFile(t)
		defer os.Remove(iconPath)

		projectHandle, err := svc.CreateProject(context.Background(), &pb_project.CreateProjectReq{
			FileName:    filename,
			ProjectName: projectName,
		})
		require.NoError(t, err)

		character, err := svc.CreateCharacter(context.Background(), &pb_project.CreateCharacterReq{
			Project: projectHandle,
			Details: &pb_project_character.BasicCharacterDetails{
				Name:        characterName,
				Description: characterDescription,
			},
			IconPath: iconPath,
		})

		require.NoError(t, err)
		require.NotEmpty(t, character.Id)

		_, err = svc.CloseProject(context.Background(), projectHandle)
		require.NoError(t, err)

		project, err := svc.OpenProject(context.Background(), &pb_project.OpenProjectReq{
			FileName: filename,
		})
		defer closeProject(t, filename, project.Handle)
		require.Len(t, project.Characters, 1)
	})

	t.Run("Test get character after create", func(t *testing.T) {
		filename := uuid.New().String() + sqlite3.DbExtension
		projectName := uuid.New().String()
		characterName := uuid.New().String()
		characterDescription := uuid.New().String()
		iconPath := newTestImageFile(t)
		defer os.Remove(iconPath)

		projectHandle, err := svc.CreateProject(context.Background(), &pb_project.CreateProjectReq{
			FileName:    filename,
			ProjectName: projectName,
		})
		require.NoError(t, err)
		defer svc.CloseProject(context.Background(), projectHandle)

		character, err := svc.CreateCharacter(context.Background(), &pb_project.CreateCharacterReq{
			Project: projectHandle,
			Details: &pb_project_character.BasicCharacterDetails{
				Name:        characterName,
				Description: characterDescription,
			},
			IconPath: iconPath,
		})

		require.NoError(t, err)
		require.NotEmpty(t, character.Id)

		returnedCharacter, err := svc.GetCharacter(context.Background(), &pb_project.GetCharacterReq{
			Character: character,
			Project:   projectHandle,
		})
		require.NoError(t, err)
		require.Equal(t, characterName, returnedCharacter.Details.Name)
		require.Equal(t, characterDescription, returnedCharacter.Details.Description)
		require.NotNil(t, returnedCharacter.Icon)
		require.True(t, len(returnedCharacter.Icon) > 0)
	})

	t.Run("Test update character with icon set", func(t *testing.T) {
		filename := uuid.New().String() + sqlite3.DbExtension
		projectName := uuid.New().String()
		characterName := uuid.New().String()
		characterDescription := uuid.New().String()
		iconPath := newTestImageFile(t)
		defer os.Remove(iconPath)

		projectHandle, err := svc.CreateProject(context.Background(), &pb_project.CreateProjectReq{
			FileName:    filename,
			ProjectName: projectName,
		})
		require.NoError(t, err)
		defer closeProject(t, filename, projectHandle)

		characterHandle, err := svc.CreateCharacter(context.Background(), &pb_project.CreateCharacterReq{
			Project: projectHandle,
			Details: &pb_project_character.BasicCharacterDetails{
				Name:        characterName,
				Description: characterDescription,
			},
			IconPath: iconPath,
		})
		require.NoError(t, err)
		require.NotEmpty(t, characterHandle.Id)

		updatedCharacterName := uuid.New().String()
		updatedCharacterDescription := uuid.New().String()
		updatedIconPath := newTestImageFile(t)
		defer os.Remove(updatedIconPath)

		_, err = svc.UpdateCharacter(context.Background(), &pb_project.UpdateCharacterReq{
			Handle:  characterHandle,
			Project: projectHandle,
			Details: &pb_project_character.BasicCharacterDetails{
				Name:        updatedCharacterName,
				Description: updatedCharacterDescription,
			},
			IconPath: updatedIconPath,
			SetIcon:  true,
		})
		require.NoError(t, err)

		returnedCharacter, err := svc.GetCharacter(context.Background(), &pb_project.GetCharacterReq{
			Character: characterHandle,
			Project:   projectHandle,
		})
		require.NoError(t, err)
		require.Equal(t, updatedCharacterName, returnedCharacter.Details.Name)
		require.Equal(t, updatedCharacterDescription, returnedCharacter.Details.Description)
		require.NotNil(t, returnedCharacter.Icon)
		require.True(t, len(returnedCharacter.Icon) > 0)
	})

	t.Run("Test update character without icon set", func(t *testing.T) {
		filename := uuid.New().String() + sqlite3.DbExtension
		projectName := uuid.New().String()
		characterName := uuid.New().String()
		characterDescription := uuid.New().String()
		iconPath := newTestImageFile(t)
		defer os.Remove(iconPath)

		projectHandle, err := svc.CreateProject(context.Background(), &pb_project.CreateProjectReq{
			FileName:    filename,
			ProjectName: projectName,
		})
		require.NoError(t, err)
		defer closeProject(t, filename, projectHandle)

		characterHandle, err := svc.CreateCharacter(context.Background(), &pb_project.CreateCharacterReq{
			Project: projectHandle,
			Details: &pb_project_character.BasicCharacterDetails{
				Name:        characterName,
				Description: characterDescription,
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, characterHandle.Id)

		updatedCharacterName := uuid.New().String()
		updatedCharacterDescription := uuid.New().String()
		updatedIconPath := newTestImageFile(t)
		defer os.Remove(updatedIconPath)

		_, err = svc.UpdateCharacter(context.Background(), &pb_project.UpdateCharacterReq{
			Handle:  characterHandle,
			Project: projectHandle,
			Details: &pb_project_character.BasicCharacterDetails{
				Name:        updatedCharacterName,
				Description: updatedCharacterDescription,
			},
			IconPath: updatedIconPath,
			SetIcon:  false,
		})
		require.NoError(t, err)

		returnedCharacter, err := svc.GetCharacter(context.Background(), &pb_project.GetCharacterReq{
			Character: characterHandle,
			Project:   projectHandle,
		})
		require.NoError(t, err)
		require.Equal(t, updatedCharacterName, returnedCharacter.Details.Name)
		require.Equal(t, updatedCharacterDescription, returnedCharacter.Details.Description)
		require.Empty(t, returnedCharacter.Icon)
	})

	t.Run("Test delete character", func(t *testing.T) {
		filename := uuid.New().String() + sqlite3.DbExtension
		projectName := uuid.New().String()
		characterName := uuid.New().String()
		characterDescription := uuid.New().String()
		iconPath := newTestImageFile(t)
		defer os.Remove(iconPath)

		projectHandle, err := svc.CreateProject(context.Background(), &pb_project.CreateProjectReq{
			FileName:    filename,
			ProjectName: projectName,
		})
		require.NoError(t, err)
		defer closeProject(t, filename, projectHandle)

		characterHandle, err := svc.CreateCharacter(context.Background(), &pb_project.CreateCharacterReq{
			Project: projectHandle,
			Details: &pb_project_character.BasicCharacterDetails{
				Name:        characterName,
				Description: characterDescription,
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, characterHandle.Id)

		resp, err := svc.DeleteCharacter(context.Background(), &pb_project.DeleteCharacterReq{
			Project: projectHandle,
			Handle:  characterHandle,
		})

		require.NotNil(t, resp)
		require.NoError(t, err)

		_, err = svc.GetCharacter(context.Background(), &pb_project.GetCharacterReq{
			Character: characterHandle,
			Project:   projectHandle,
		})
		require.Error(t, err)
	})

	t.Run("Test parse nil projectId", func(t *testing.T) {
		_, err := svc.GetCharacter(context.Background(), &pb_project.GetCharacterReq{
			Character: &pb_project_character.CharacterHandle{
				Id: uuid.New().String(),
			},
			Project: nil,
		})
		require.Error(t, err)
	})

	t.Run("Test parse nil characterId", func(t *testing.T) {
		_, err := svc.GetCharacter(context.Background(), &pb_project.GetCharacterReq{
			Character: &pb_project_character.CharacterHandle{
				Id: uuid.New().String(),
			},
			Project: &pb_project.ProjectHandle{
				Id: uuid.New().String(),
			},
		})
		require.Error(t, err)
	})

	t.Run("Test search character", func(t *testing.T) {
		filename := uuid.New().String() + sqlite3.DbExtension
		projectName := uuid.New().String()
		characterName := uuid.New().String()
		characterDescription := uuid.New().String()
		iconPath := newTestImageFile(t)
		defer os.Remove(iconPath)

		projectHandle, err := svc.CreateProject(context.Background(), &pb_project.CreateProjectReq{
			FileName:    filename,
			ProjectName: projectName,
		})
		require.NoError(t, err)
		defer closeProject(t, filename, projectHandle)

		characterHandle, err := svc.CreateCharacter(context.Background(), &pb_project.CreateCharacterReq{
			Project: projectHandle,
			Details: &pb_project_character.BasicCharacterDetails{
				Name:        characterName,
				Description: characterDescription,
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, characterHandle.Id)

		searchRes, err := svc.SearchCharacter(context.Background(), &pb_project.SearchCharacterReq{
			Project: projectHandle,
			Query:   characterName,
		})

		require.NoError(t, err)
		require.Len(t, searchRes.Details, 1)
		require.Equal(t, characterHandle, searchRes.Details[0].Handle)
	})
}
