package imagecompression

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"

	"github.com/charmbracelet/log"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
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

	log.Info("Compressed image", loggertags.TagLength, w.Len())
	return w.Bytes(), nil
}
