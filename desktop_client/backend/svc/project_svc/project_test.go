package projectsvc_test

import (
	"bytes"
	"context"
	"image/jpeg"
	"os"
	"testing"

	testutils "github.com/djpiper28/rpg-book/common/test_utils"
	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project_character"
	projectsvc "github.com/djpiper28/rpg-book/desktop_client/backend/svc/project_svc"
	testdbutils "github.com/djpiper28/rpg-book/desktop_client/test_db_utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestProjectSvc(t *testing.T) {
	db, closeDb := testdbutils.GetPrimaryDb()
	defer closeDb()

	svc := projectsvc.New(db)
	defer svc.Close()

	closeProjectWithoutDelete := func(t *testing.T, filename string, handle *pb_project.ProjectHandle) {
		_, err := svc.CloseProject(context.Background(), handle)
		require.NoError(t, err)
	}

	closeProject := func(t *testing.T, filename string, handle *pb_project.ProjectHandle) {
		defer os.Remove(filename)

		closeProjectWithoutDelete(t, filename, handle)
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
		filename := uuid.New().String() + database.DbExtension
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
		filename := uuid.New().String() + database.DbExtension
		name := uuid.New().String()

		handle, err := svc.CreateProject(context.Background(), &pb_project.CreateProjectReq{
			FileName:    filename,
			ProjectName: name,
		})
		require.NoError(t, err)
		closeProjectWithoutDelete(t, filename, handle)

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
		filename := uuid.New().String() + database.DbExtension
		name := uuid.New().String()

		handle, err := svc.CreateProject(context.Background(), &pb_project.CreateProjectReq{
			FileName:    filename,
			ProjectName: name,
		})
		require.NoError(t, err)
		closeProjectWithoutDelete(t, filename, handle)

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
		filename := uuid.New().String() + database.DbExtension
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
		filename := uuid.New().String() + database.DbExtension
		projectName := uuid.New().String()
		characterName := uuid.New().String()
		characterDescription := uuid.New().String()
		icon := bytes.NewBuffer([]byte{})

		err := jpeg.Encode(icon, testutils.NewTestImage(100, 100), nil)
		require.NoError(t, err)

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
				Icon:        icon.Bytes(),
			},
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
		filename := uuid.New().String() + database.DbExtension
		projectName := uuid.New().String()
		characterName := uuid.New().String()
		characterDescription := uuid.New().String()
		icon := bytes.NewBuffer([]byte{})

		err := jpeg.Encode(icon, testutils.NewTestImage(100, 100), nil)
		require.NoError(t, err)

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
				Icon:        icon.Bytes(),
			},
		})

		require.NoError(t, err)
		require.NotEmpty(t, character.Id)

		returnedCharacter, err := svc.GetCharacter(context.Background(), &pb_project.GetCharacterReq{
			Character: character,
			Project:   projectHandle,
		})
		require.NoError(t, err)
		require.Equal(t, characterName, returnedCharacter.Name)
		require.Equal(t, characterDescription, returnedCharacter.Description)
		require.NotNil(t, returnedCharacter.Icon)
		require.True(t, len(returnedCharacter.Icon) > 0)
	})

	t.Run("Test update character with icon set", func(t *testing.T) {
		filename := uuid.New().String() + database.DbExtension
		projectName := uuid.New().String()
		characterName := uuid.New().String()
		characterDescription := uuid.New().String()
		icon := bytes.NewBuffer([]byte{})

		err := jpeg.Encode(icon, testutils.NewTestImage(100, 100), nil)
		require.NoError(t, err)

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
				Icon:        icon.Bytes(),
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, characterHandle.Id)

		updatedCharacterName := uuid.New().String()
		updatedCharacterDescription := uuid.New().String()
		updatedIcon := bytes.NewBuffer([]byte{})
		err = jpeg.Encode(updatedIcon, testutils.NewTestImage(50, 50), nil)
		require.NoError(t, err)

		_, err = svc.UpdateCharacter(context.Background(), &pb_project.UpdateCharacterReq{
			Handle:  characterHandle,
			Project: projectHandle,
			Details: &pb_project_character.BasicCharacterDetails{
				Name:        updatedCharacterName,
				Description: updatedCharacterDescription,
				Icon:        updatedIcon.Bytes(),
			},
			SetImage: true,
		})
		require.NoError(t, err)

		returnedCharacter, err := svc.GetCharacter(context.Background(), &pb_project.GetCharacterReq{
			Character: characterHandle,
			Project:   projectHandle,
		})
		require.NoError(t, err)
		require.Equal(t, updatedCharacterName, returnedCharacter.Name)
		require.Equal(t, updatedCharacterDescription, returnedCharacter.Description)
		require.NotNil(t, returnedCharacter.Icon)
		require.True(t, len(returnedCharacter.Icon) > 0)
    // Image is compressed by default
		require.NotEqual(t, updatedIcon.Bytes(), returnedCharacter.Icon)
	})

	t.Run("Test update character without icon set", func(t *testing.T) {
		filename := uuid.New().String() + database.DbExtension
		projectName := uuid.New().String()
		characterName := uuid.New().String()
		characterDescription := uuid.New().String()
		icon := bytes.NewBuffer([]byte{})

		err := jpeg.Encode(icon, testutils.NewTestImage(100, 100), nil)
		require.NoError(t, err)

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
				Icon:        icon.Bytes(),
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, characterHandle.Id)

		updatedCharacterName := uuid.New().String()
		updatedCharacterDescription := uuid.New().String()

		_, err = svc.UpdateCharacter(context.Background(), &pb_project.UpdateCharacterReq{
			Handle:  characterHandle,
			Project: projectHandle,
			Details: &pb_project_character.BasicCharacterDetails{
				Name:        updatedCharacterName,
				Description: updatedCharacterDescription,
				Icon:        []byte{},
			},
			SetImage: false,
		})
		require.NoError(t, err)

		returnedCharacter, err := svc.GetCharacter(context.Background(), &pb_project.GetCharacterReq{
			Character: characterHandle,
			Project:   projectHandle,
		})
		require.NoError(t, err)
		require.Equal(t, updatedCharacterName, returnedCharacter.Name)
		require.Equal(t, updatedCharacterDescription, returnedCharacter.Description)
		require.NotNil(t, returnedCharacter.Icon)
		require.True(t, len(returnedCharacter.Icon) > 0)
		require.Equal(t, icon.Bytes(), returnedCharacter.Icon)
	})
}
