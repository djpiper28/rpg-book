package testutils_test

import (
	"image/jpeg"
	"os"
	"testing"

	testutils "github.com/djpiper28/rpg-book/common/test_utils"
	"github.com/stretchr/testify/require"
)

func TestImageGenerator(t *testing.T) {
	img := testutils.NewTestImage(1000, 1000)
	f, err := os.OpenFile("test.jpg", os.O_WRONLY|os.O_RDONLY|os.O_CREATE, os.FileMode(0o666))
	require.NoError(t, err)
	defer f.Close()

	err = jpeg.Encode(f, img, nil)
	require.NoError(t, err)
}
