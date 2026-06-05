package optimize

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
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
	_             = draw.NearestNeighbor
	_             = draw.ApproxBiLinear
	_             = draw.BiLinear
	defaultScaler = draw.CatmullRom
	jpegEncodeOpt = jpeg.Options{Quality: 90}
)

func elapsed(name string) func() {
	t := time.Now()
	return func() {
		fmt.Println(name, "elapsed:", time.Since(t))
	}
}

type customReader struct {
	ctx context.Context
	r   io.Reader
}

func (x customReader) Read(p []byte) (int, error) {
	select {
	case <-x.ctx.Done():
		return 0, x.ctx.Err()
	default:
		return x.r.Read(p)
	}
}

func ResizeToMax(ctx context.Context, r io.Reader, w io.Writer) error {
	defer elapsed("image decode")()
	rr := customReader{ctx: ctx, r: r}
	i, _, err := image.Decode(rr)
	if err != nil {
		return err
	}
	width, height := i.Bounds().Dx(), i.Bounds().Dy()
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
	return jpeg.Encode(w, newImg, &jpegEncodeOpt)
}

func getLimitSize(width, height, limit int) (newWidth, newHeight int) {
	limitEdge := min(height, width)
	f := float64((limitEdge * limit))
	newW := math.Round(f / float64(height))
	newH := math.Round(f / float64(width))
	return int(newW), int(newH)
}
