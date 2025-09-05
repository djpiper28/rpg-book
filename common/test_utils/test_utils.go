package testutils

import (
	"context"
	"image"
	"image/color"
	"os"
	"testing"
	"time"

	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/djpiper28/rpg-book/desktop_client/backend/project"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func NewTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Minute/2)
}

func NewTestImage(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := range h {
		for x := range w {
			val := x + y
			img.SetRGBA(x, y, color.RGBA{
				R: uint8(val % 255),
				G: uint8(val % 255),
				B: uint8(val % 255),
			})
		}
	}

	return img
}

func NewProject(t *testing.T) (*project.Project, func()) {
	filename := uuid.New().String() + database.DbExtension

	projectName := uuid.New().String()

	project, err := project.Create(filename, projectName)
	require.NoError(t, err)
	require.Equal(t, project.Settings.Name, projectName)
	require.Equal(t, project.Filename, filename)

	return project, func() {
		os.Remove(filename)
	}
}
