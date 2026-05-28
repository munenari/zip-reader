package util

import (
	"archive/zip"
	"fmt"
	"io"
	"sort"
)

type (
	ZipReader struct {
		filename string
		zr       *zip.ReadCloser
	}
)

func NewZipReader(filename string) *ZipReader {
	return &ZipReader{filename: filename}
}
func NewZipReaderWithOpen(filename string) (*ZipReader, error) {
	x := NewZipReader(filename)
	err := x.Open()
	return x, err
}

func (x *ZipReader) Open() error {
	zr, err := zip.OpenReader(x.filename)
	if err != nil {
		return err
	}
	shrinked := make([]*zip.File, 0, len(zr.File))
	for _, f := range zr.File {
		// shrink index if dir
		if !f.FileInfo().IsDir() {
			shrinked = append(shrinked, f)
		}
	}
	sort.Slice(shrinked, func(i, j int) bool {
		return shrinked[i].Name < shrinked[j].Name
	})
	zr.File = shrinked
	x.zr = zr
	return err
}
func (x *ZipReader) Close() error {
	if x.zr == nil {
		return nil
	}
	err := x.zr.Close()
	x.zr = nil
	return err
}
func (x *ZipReader) Len() int {
	if x.zr == nil {
		return -1
	}
	return len(x.zr.File)
}
func (x *ZipReader) At(i int) (io.ReadCloser, error) {
	if x.zr == nil {
		return nil, fmt.Errorf("zipfile not was opened")
	}
	if i < 0 || i >= len(x.zr.File) {
		return nil, fmt.Errorf("out of file index in zip, (%d/%d)", i, len(x.zr.File))
	}
	f := x.zr.File[i]
	return f.Open()
}
