package testutils

import (
	"context"
	"image"
	"image/color"
	"time"
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
