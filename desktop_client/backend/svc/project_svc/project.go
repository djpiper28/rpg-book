package projectsvc

import (
	"context"
	"errors"
	"sync"
	"time"

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
	lock      sync.Mutex
}

func New(db *database.Db) *ProjectSvc {
	return &ProjectSvc{
		primaryDb: db,
		projects:  make(map[uuid.UUID]*project.Project),
	}
}

func (p *ProjectSvc) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, project := range p.projects {
		project.Close()
	}
}

func (p *ProjectSvc) trackProject(proj *project.Project) (uuid.UUID, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, otherProject := range p.projects {
		if otherProject.Filename == proj.Filename {
			return uuid.Nil, errors.New("The project is already open")
		}
	}

	id := uuid.New()
	p.projects[id] = proj
	return id, nil
}

// Called by CreateProject
func (p *ProjectSvc) insertProjectOpened(ctx context.Context, filename, projectname string) error {
	openedAt := time.Now()

	_, err := p.primaryDb.Db.ExecContext(ctx, "INSERT INTO recently_opened(file_name, project_name, last_opened) VALUES (?, ?, ?);", filename, projectname, openedAt)
	if err != nil {
		return errors.Join(errors.New("Cannot insert into recently opened projects"), err)
	}

	return nil
}

// Called by OpenProject
func (p *ProjectSvc) updateProjectOpened(ctx context.Context, filename, projectname string) error {
	openedAt := time.Now()

	_, err := p.primaryDb.Db.ExecContext(ctx, "UPDATE recently_opened SET file_name=?, project_name=?, last_opened=?;", filename, projectname, openedAt)
	if err != nil {
		return errors.Join(errors.New("Cannot update recently opened projects"), err)
	}

	return nil
}

func (p *ProjectSvc) CreateProject(ctx context.Context, in *pb_project.CreateProjectReq) (*pb_project.ProjectHandle, error) {
	proj, err := project.Create(in.FileName, in.ProjectName)
	if err != nil {
		return nil, err
	}

	err = p.insertProjectOpened(ctx, proj.Filename, proj.Settings.Name)
	if err != nil {
		return nil, err
	}

	id, err := p.trackProject(proj)
	if err != nil {
		return nil, err
	}

	return &pb_project.ProjectHandle{
		Id: id.String(),
	}, nil
}

func (p *ProjectSvc) OpenProject(ctx context.Context, in *pb_project.OpenProjectReq) (*pb_project.ProjectHandle, error) {
	proj, err := project.Open(in.FileName)
	if err != nil {
		return nil, err
	}

	err = p.insertProjectOpened(ctx, proj.Filename, proj.Settings.Name)
	if err != nil {
		return nil, err
	}

	id, err := p.trackProject(proj)
	if err != nil {
		return nil, err
	}

	return &pb_project.ProjectHandle{
		Id: id.String(),
	}, nil
}

func (p *ProjectSvc) RecentProjects(ctx context.Context, in *pb_common.Empty) (*pb_project.RecentProjectsResp, error) {
	return nil, errors.New("Unimplemented")
}
