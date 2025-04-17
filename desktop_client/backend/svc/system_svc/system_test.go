package systemsvc_test

import (
	"context"
	"testing"

	buildinfo "github.com/djpiper28/rpg-book/common/build_info"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_system"
	systemsvc "github.com/djpiper28/rpg-book/desktop_client/backend/svc/system_svc"
	testdbutils "github.com/djpiper28/rpg-book/desktop_client/test_db_utils"
	"github.com/stretchr/testify/require"
)

func TestSettings(t *testing.T) {
	db, cleanup := testdbutils.GetPrimaryDb()
	defer cleanup()

	system := systemsvc.New(db)

	t.Run("Init", func(t *testing.T) {
		require.NotNil(t, system)
	})

	t.Run("Get Settings", func(t *testing.T) {
		settings, err := system.GetSettings(context.Background(), &pb_common.Empty{})
		require.NoError(t, err)
		require.NotNil(t, settings)

		require.True(t, settings.DarkMode)
		require.False(t, settings.DevMode)
	})

	t.Run("Set Settings", func(t *testing.T) {
		settings, err := system.GetSettings(context.Background(), &pb_common.Empty{})
		require.NoError(t, err)

		settings.DevMode = false
		settings.DarkMode = false
		_, err = system.SetSettings(context.Background(), settings)
		require.NoError(t, err)

		settings2, err := system.GetSettings(context.Background(), &pb_common.Empty{})
		require.NoError(t, err)

		require.Equal(t, settings, settings2)
		require.False(t, settings.DarkMode)
		require.False(t, settings.DevMode)
	})

	t.Run("Set Settings (no change)", func(t *testing.T) {
		settings, err := system.GetSettings(context.Background(), &pb_common.Empty{})
		require.NoError(t, err)

		_, err = system.SetSettings(context.Background(), settings)
		require.NoError(t, err)

		settings2, err := system.GetSettings(context.Background(), &pb_common.Empty{})
		require.NoError(t, err)

		require.Equal(t, settings, settings2)
	})

	t.Run("Call Logger", func(t *testing.T) {
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
	})

	t.Run("Call Version", func(t *testing.T) {
		version, err := system.GetVersion(context.Background(), &pb_common.Empty{})
		require.NoError(t, err)
		require.Equal(t, version.Version, buildinfo.Version)
	})
}
