package projectsvc

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/djpiper28/rpg-book/desktop_client/backend/model"
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
  absPath, err := filepath.Abs(filename)
  if err != nil {
    filename = absPath
  }

	_, err = p.primaryDb.Db.ExecContext(ctx,
		"INSERT INTO recently_opened(file_name, project_name, last_opened) VALUES (?, ?, ?);",
		filename,
		projectname,
		openedAt)
	if err != nil {
		return errors.Join(errors.New("Cannot insert into recently opened projects"), err)
	}

	return nil
}

// Called by OpenProject
func (p *ProjectSvc) updateProjectOpened(ctx context.Context, filename, projectname string) error {
	openedAt := time.Now()
  absPath, err := filepath.Abs(filename)
  if err != nil {
    filename = absPath
  }

	_, err = p.primaryDb.Db.ExecContext(ctx,
		"UPDATE recently_opened SET project_name=?, last_opened=? WHERE file_name=?;",
		filename,
		projectname,
		openedAt)
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

	err = p.updateProjectOpened(ctx, proj.Filename, proj.Settings.Name)
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

func (p *ProjectSvc) CloseProject(ctx context.Context, in *pb_project.ProjectHandle) (*pb_common.Empty, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, errors.Join(errors.New("Invalid project handle"), err)
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	proj, found := p.projects[id]
	if !found {
		return nil, errors.New("Cannot find proejct")
	}

	delete(p.projects, id)
	proj.Close()
	return &pb_common.Empty{}, nil
}

func (p *ProjectSvc) RecentProjects(ctx context.Context, in *pb_common.Empty) (*pb_project.RecentProjectsResp, error) {
	rows, err := p.primaryDb.Db.QueryxContext(ctx, "SELECT * FROM recently_opened ORDER BY last_opened DESC LIMIT 10;")
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get recently opened projects"), err)
	}

	projects := make([]*pb_project.RecentProject, 0)
	for rows.Next() {
		var recent model.RecentlyOpened
		err = rows.StructScan(&recent)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot scan recently opened projects"), err)
		}
		var size int64 = 0
		stat, err := os.Stat(recent.FileName)
		if err != nil {
			log.Warn("Recently opened project does not exist", loggertags.TagFileName, recent.FileName)
			size = -1
		} else {
			size = stat.Size()
		}

		projects = append(projects, &pb_project.RecentProject{
			ProjectName:   recent.ProjectName,
			FileName:      recent.FileName,
			LastOpened:    recent.LastOpened,
			FileSizeBytes: size,
		})
	}

	return &pb_project.RecentProjectsResp{
		Projects: projects,
	}, nil
}
