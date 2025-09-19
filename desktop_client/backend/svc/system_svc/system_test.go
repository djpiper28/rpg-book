package systemsvc_test

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"os"
	"testing"

	buildinfo "github.com/djpiper28/rpg-book/common/build_info"
	imagecompression "github.com/djpiper28/rpg-book/common/image/image_compression"
	testutils "github.com/djpiper28/rpg-book/common/test_utils"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_system"
	systemsvc "github.com/djpiper28/rpg-book/desktop_client/backend/svc/system_svc"
	testdbutils "github.com/djpiper28/rpg-book/desktop_client/test_db_utils"
	"github.com/google/uuid"
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
				{
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

	t.Run("Test read file (failure)", func(t *testing.T) {
		_, err := system.ReadImageFile(context.Background(), &pb_system.ReadImageFileReq{
			Filepath: "./not_a_real_file",
		})

		require.Error(t, err)
	})

	t.Run("Test read non image file", func(t *testing.T) {
		const FileName = "./testing_file.txt"
		data := []byte(uuid.New().String())

		err := os.WriteFile(FileName, data, 0o555)
		defer os.Remove(FileName)
		require.NoError(t, err)

		_, err = system.ReadImageFile(context.Background(), &pb_system.ReadImageFileReq{
			Filepath: FileName,
		})

		require.Error(t, err)
	})

	t.Run("Test read image file (not icon)", func(t *testing.T) {
		const FileName = "./testing_file.jpeg"
		const imgSize = 2500

		img := testutils.NewTestImage(imgSize, imgSize)
		buffer := bytes.NewBuffer([]byte{})

		err := png.Encode(buffer, img)
		require.NoError(t, err)

		data := buffer.Bytes()

		err = os.WriteFile(FileName, data, 0o555)
		defer os.Remove(FileName)
		require.NoError(t, err)

		resp, err := system.ReadImageFile(context.Background(), &pb_system.ReadImageFileReq{
			Filepath: FileName,
			IsIcon:   false,
		})

		require.NoError(t, err)
		require.NotNil(t, data)
		require.True(t, len(data) > 0)

		img, format, err := image.Decode(bytes.NewBuffer(resp.Data))
		require.Equal(t, "jpeg", format)
		require.Equal(t, imgSize, img.Bounds().Dx())
		require.Equal(t, imgSize, img.Bounds().Dy())
		require.Greater(t, imgSize, int(imagecompression.MaxDimension), "this checks that the icon compression code was never actually called")
	})

	t.Run("Test read image file (icon)", func(t *testing.T) {
		const FileName = "./testing_file.jpeg"
		const imgSize = 2000

		img := testutils.NewTestImage(imgSize, imgSize)
		buffer := bytes.NewBuffer([]byte{})

		err := png.Encode(buffer, img)
		require.NoError(t, err)

		data := buffer.Bytes()

		err = os.WriteFile(FileName, data, 0o555)
		defer os.Remove(FileName)
		require.NoError(t, err)

		resp, err := system.ReadImageFile(context.Background(), &pb_system.ReadImageFileReq{
			Filepath: FileName,
			IsIcon:   true,
		})

		require.NoError(t, err)
		require.NotNil(t, data)
		require.True(t, len(data) > 0)

		img, format, err := image.Decode(bytes.NewBuffer(resp.Data))
		require.Equal(t, "jpeg", format)
		require.Equal(t, int(imagecompression.MaxDimension), img.Bounds().Dx())
		require.Equal(t, int(imagecompression.MaxDimension), img.Bounds().Dy())
	})
}
