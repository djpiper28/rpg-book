package projectsvc

import (
	"context"
	"errors"

	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_project"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project"
	"github.com/google/uuid"
)

type ProjectSvc struct {
	pb_project.UnimplementedProjectSvcServer
	primaryDb *database.Db
	projects  map[uuid.UUID]*project.Project
}

func New(db *database.Db) *ProjectSvc {
	return &ProjectSvc{
		primaryDb: db,
		projects:  make(map[uuid.UUID]*project.Project),
	}
}

func (p *ProjectSvc) Close() {
	for _, project := range p.projects {
		project.Close()
	}
}

func (p *ProjectSvc) CreateProject(ctx context.Context, in *pb_project.CreateProjectReq) (*pb_project.ProjectHandle, error) {
	return nil, errors.New("Unimplemented")
}

func (p *ProjectSvc) OpenProject(ctx context.Context, in *pb_project.OpenProjectReq) (*pb_project.ProjectHandle, error) {
	return nil, errors.New("Unimplemented")
}

func (p *ProjectSvc) RecentProjects(ctx context.Context, in *pb_common.Empty) (*pb_project.RecentProjectsResp, error) {
	return nil, errors.New("Unimplemented")
}
