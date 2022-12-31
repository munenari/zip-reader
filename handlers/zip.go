package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"sort"
)

func readPage(filename string, index int) (io.ReadCloser, io.Closer, error) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return nil, nil, err
	}
	// defer r.Close()
	for i, f := range r.File {
		if f.FileInfo().IsDir() {
			r.File = append(r.File[:i], r.File[i+1:]...)
		}
	}
	if index < 0 || index >= len(r.File) {
		defer r.Close()
		return nil, nil, fmt.Errorf("out of file index in zip")
	}
	sort.Slice(r.File, func(i, j int) bool {
		return r.File[i].Name < r.File[j].Name
	})
	f := r.File[index]
	// fmt.Println(f.Name)
	// fmt.Println(f.FileInfo().IsDir())
	rc, err := f.Open()
	return rc, r, err
}

func length(filename string) (int, error) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return 0, err
	}
	defer r.Close()
	for i, f := range r.File {
		if f.FileInfo().IsDir() {
			r.File = append(r.File[:i], r.File[i+1:]...)
		}
	}
	return len(r.File), nil
}
