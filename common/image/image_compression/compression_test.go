package imagecompression_test

import (
	"bytes"
	"image"
	"testing"

	imagecompression "github.com/djpiper28/rpg-book/common/image/image_compression"
	testutils "github.com/djpiper28/rpg-book/common/test_utils"
	"github.com/stretchr/testify/require"
)

func TestCompresion(t *testing.T) {
	img := testutils.NewTestImage(1000, 1000)
	imgBytes, err := imagecompression.Compress(img)
	require.NoError(t, err)

	readImage, format, err := image.Decode(bytes.NewBuffer(imgBytes))
	require.NoError(t, err)
	require.Equal(t, format, "jpeg")
	require.NotNil(t, readImage)
}
