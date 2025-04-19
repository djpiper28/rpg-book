package projectsvc_test

import (
	"context"
	"testing"

	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
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
		require.NotEmpty(t, handle)

		var count int
		var dbName string
		rows := db.Db.QueryRow("SELECT COUNT(file_name), project_name FROM recently_opened WHERE file_name=?;", filename)
		err = rows.Scan(&count, &dbName)

		require.NoError(t, err)
		require.Equal(t, 1, count)
		require.Equal(t, name, dbName)
	})
}
