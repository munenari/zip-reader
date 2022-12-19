package optimize

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"math"

	"golang.org/x/image/draw"
)

const (
	limitSize = 4000
)

var (
	encoder = &png.Encoder{
		CompressionLevel: png.BestSpeed,
		BufferPool:       NewBufferPool(),
	}
	_             = draw.NearestNeighbor
	_             = draw.ApproxBiLinear
	_             = draw.BiLinear
	defaultScaler = draw.CatmullRom
)

func ResizeToMax(r io.Reader, w io.Writer) error {
	i, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	width := i.Bounds().Dx()
	height := i.Bounds().Dy()
	if width < limitSize && height < limitSize {
		return fmt.Errorf("image optimize was not needed")
	}
	log.Println("resizing...")
	newImgData := &image.RGBA{}
	_ = newImgData
	if height >= width {
		f := float64((width * limitSize))
		w := math.Round(f / float64(height))
		newImgData = image.NewRGBA(image.Rect(0, 0, int(w), limitSize))
	} else {
		f := float64((limitSize * height))
		h := math.Round(f / float64(width))
		newImgData = image.NewRGBA(image.Rect(0, 0, limitSize, int(h)))
	}
	defaultScaler.Scale(newImgData, newImgData.Bounds(), i, i.Bounds(), draw.Over, nil)
	return encoder.Encode(w, newImgData)
}
