package systemsvc

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_system"
)

type SystemSvc struct {
	pb_system.UnimplementedSystemSvcServer
}

func New() *SystemSvc {
	return &SystemSvc{}
}

func (s *SystemSvc) GetSettings(ctx context.Context, in *pb_common.Empty) (*pb_system.Settings, error) {
	return &pb_system.Settings{}, nil
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
