package imagecompression

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
)

const (
	JpegQuality = 80
)

func Compress(in image.Image) ([]byte, error) {
	w := bytes.NewBuffer([]byte{})
	err := jpeg.Encode(w, in, &jpeg.Options{
		Quality: JpegQuality,
	})
	if err != nil {
		return nil, errors.Join(errors.New("Cannot compress image"), err)
	}

	return w.Bytes(), nil
}
