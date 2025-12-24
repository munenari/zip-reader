package optimize

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"time"

	"golang.org/x/image/draw"
)

const (
	limitSize = 2000
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
	jpegEncodeOpt = jpeg.Options{Quality: 85}
)

func elapsed(name string) func() {
	t := time.Now()
	return func() {
		fmt.Println(name, "elapsed:", time.Since(t))
	}
}

func ResizeToMax(ctx context.Context, r io.Reader, w io.Writer) error {
	defer elapsed("image decode")()
	i, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	width, height := i.Bounds().Dx(), i.Bounds().Dy()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	log.Printf("image size: %dx%d\n", width, height)
	var newImg image.Image
	if width < limitSize && height < limitSize {
		newImg = i
	} else {
		newW, newH := getLimitSize(width, height, limitSize)
		newImgData := image.NewRGBA(image.Rect(0, 0, newW, newH))
		defaultScaler.Scale(newImgData, newImgData.Bounds(), i, i.Bounds(), draw.Over, nil)
		newImg = newImgData
	}
	// return encoder.Encode(w, newImg)
	return jpeg.Encode(w, newImg, &jpegEncodeOpt)
}

func getLimitSize(width, height, limit int) (newWidth, newHeight int) {
	limitEdge := min(height, width)
	f := float64((limitEdge * limit))
	newW := math.Round(f / float64(height))
	newH := math.Round(f / float64(width))
	return int(newW), int(newH)
}
