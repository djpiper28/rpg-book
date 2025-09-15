package systemsvc

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	buildinfo "github.com/djpiper28/rpg-book/common/build_info"
	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/djpiper28/rpg-book/desktop_client/backend/model"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_system"
)

type SystemSvc struct {
	pb_system.UnimplementedSystemSvcServer
	db *database.Db
}

func New(db *database.Db) *SystemSvc {
	return &SystemSvc{
		db: db,
	}
}

func (s *SystemSvc) GetSettings(ctx context.Context, in *pb_common.Empty) (*pb_system.Settings, error) {
	rows := s.db.Db.QueryRowx("SELECT * FROM settings;")

	var settings model.Settings
	err := rows.StructScan(&settings)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get settings"), err)
	}

	return settings.ToProto(), nil
}

func (s *SystemSvc) SetSettings(ctx context.Context, in *pb_system.Settings) (*pb_common.Empty, error) {
	_, err := s.db.Db.Exec("UPDATE settings SET dev_mode=?, dark_mode=?, compress_images=?;", in.DevMode, in.DarkMode, in.CompressImages)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot set settings"), err)
	}

	return &pb_common.Empty{}, nil
}

func (s *SystemSvc) Log(ctx context.Context, req *pb_system.LogRequest) (*pb_common.Empty, error) {
	const callerTag = "js_caller"

	switch req.GetLevel() {
	case pb_system.LogLevel_INFO:
		log.Info(req.GetMessage(), callerTag, req.GetCaller(), "properties", req.GetProperties())
	case pb_system.LogLevel_WARNING:
		log.Warn(req.GetMessage(), callerTag, req.GetCaller(), "properties", req.GetProperties())
	case pb_system.LogLevel_ERROR:
		log.Error(req.GetMessage(), callerTag, req.GetCaller(), "properties", req.GetProperties())
	case pb_system.LogLevel_FATAL:
		log.Fatal(req.GetMessage(), callerTag, req.GetCaller(), "properties", req.GetProperties())
	default:
		log.Warn(req.GetMessage(), callerTag, req.GetCaller(),
			"unknownLevel", req.GetLevel(), "properties", req.GetProperties())
	}

	return &pb_common.Empty{}, nil
}

func (s *SystemSvc) GetVersion(ctx context.Context, req *pb_common.Empty) (*pb_system.Version, error) {
	return &pb_system.Version{Version: buildinfo.Version}, nil
}

func (s *SystemSvc) ReadFile(ctx context.Context, req *pb_system.ReadFileReq) (*pb_system.ReadFileRes, error) {
	data, err := os.ReadFile(req.Filepath)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Cannot read file %s", req.Filepath), err)
	}

	return &pb_system.ReadFileRes{
		Data: data,
	}, nil
}
