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
	"time"

	"golang.org/x/image/draw"
)

const (
	limitSize    = 2000
	parallelSize = 4
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
	locker        = make(chan struct{}, parallelSize)
)

func elapsed(name string) func() {
	t := time.Now()
	return func() {
		fmt.Println(name, "elapsed:", time.Since(t))
	}
}

func ResizeToMax(r io.Reader, w io.Writer) error {
	defer elapsed("image decode")()
	locker <- struct{}{}
	defer func() {
		<-locker
	}()
	i, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	width := i.Bounds().Dx()
	height := i.Bounds().Dy()
	log.Printf("image size: %dx%d\n", width, height)
	if width < limitSize && height < limitSize {
		return encoder.Encode(w, i)
	}
	var newImgData *image.RGBA
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
