package systemsvc_test

import (
	"context"
	"testing"

	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_system"
	systemsvc "github.com/djpiper28/rpg-book/desktop_client/backend/svc/system_svc"
	"github.com/stretchr/testify/require"
)

var system pb_system.SystemSvcServer = systemsvc.New()

func TestNew(t *testing.T) {
	require.NotNil(t, system)
}

func TestGetSettings(t *testing.T) {
	settings, err := system.GetSettings(context.Background(), &pb_common.Empty{})
	require.NoError(t, err)
	require.NotNil(t, settings)
}

func TestLogger(t *testing.T) {
	_, err := system.Log(context.Background(), &pb_system.LogRequest{
		Caller:  "test_caller.js:123",
		Level:   pb_system.LogLevel_WARNING,
		Message: "Hello world",
		Properties: []*pb_system.LogProperty{
			&pb_system.LogProperty{
				Key:   "key",
				Value: "Value",
			},
		},
	})

  require.NoError(t, err)
}
