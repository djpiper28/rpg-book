package projectsvc_test

import (
	"context"
	"os"
	"testing"

	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project"
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

		projectHandle, err := svc.CreateProject(context.Background(), &pb_project.CreateProjectReq{
			FileName:    filename,
			ProjectName: projectName,
		})
		require.NoError(t, err)

		character, err := svc.CreateCharacter(context.Background(), &pb_project.CreateCharacterReq{
			Name:        characterName,
			Project:     projectHandle,
			Description: characterDescription,
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

		projectHandle, err := svc.CreateProject(context.Background(), &pb_project.CreateProjectReq{
			FileName:    filename,
			ProjectName: projectName,
		})
		require.NoError(t, err)
		defer svc.CloseProject(context.Background(), projectHandle)

		character, err := svc.CreateCharacter(context.Background(), &pb_project.CreateCharacterReq{
			Name:        characterName,
			Project:     projectHandle,
			Description: characterDescription,
		})

		require.NoError(t, err)
		require.NotEmpty(t, character.Id)

		returnedCharacter, err := svc.GetCharacter(context.Background(), &pb_project.GetCharacterReq{
			Character: character,
			Project:   projectHandle,
		})
		require.NoError(t, err)
		require.Equal(t, characterName, returnedCharacter.Name)
		require.Equal(t, "", returnedCharacter.Description)
		require.Equal(t, []byte(nil), returnedCharacter.Icon)
	})
}
