package imagecompression

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/ccitt"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/vector"
	_ "golang.org/x/image/vp8"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"

	"github.com/charmbracelet/log"
	"github.com/disintegration/imaging"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
)

const (
	JpegQuality  = 80
	MaxDimension = 1000.0
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

func CompressIcon(in image.Image) ([]byte, error) {
	if in.Bounds().Max.X > MaxDimension || in.Bounds().Max.Y > MaxDimension {
		scaleX := MaxDimension / float32(in.Bounds().Dx())
		scaleY := MaxDimension / float32(in.Bounds().Dy())

		scale := max(scaleX, scaleY)
		in = imaging.Resize(in,
			int(float32(in.Bounds().Dx())*scale),
			int(float32(in.Bounds().Dy())*scale),
			imaging.Lanczos)
		log.Info("Resized image",
			loggertags.TagScale, scale,
			loggertags.TagWidth, in.Bounds().Dx(),
			loggertags.TagHeight, in.Bounds().Dy())
	}

	return Compress(in)
}
