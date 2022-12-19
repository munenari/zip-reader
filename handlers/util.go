package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"io/fs"
	"path"
)

type (
	DirInfo struct {
		Name       string `json:"name"`
		HashedName string `json:"hashed_name"`
		IsDir      bool   `json:"is_dir"`
	}
)

func gzipPath(p string) (string, error) {
	buf := &bytes.Buffer{}
	gw, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if err != nil {
		return "", nil
	}
	defer gw.Close()
	if _, err := gw.Write([]byte(p)); err != nil {
		return "", nil
	}
	gw.Close()
	return base64.URLEncoding.EncodeToString(buf.Bytes()), nil
}

func ungzipPath(str string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	buf := bytes.NewReader(b)
	gr, err := gzip.NewReader(buf)
	if err != nil {
		return "", nil
	}
	defer gr.Close()
	res, err := io.ReadAll(gr)
	return string(res), err
}

func newDirInfo(baseDir, name string, isDir bool) (*DirInfo, error) {
	hashedName, err := gzipPath(path.Join(baseDir, name))
	if err != nil {
		return nil, err
	}
	return &DirInfo{
		Name:       name,
		HashedName: hashedName,
		IsDir:      isDir,
	}, nil
}

func dirAndFileName(d fs.DirEntry) string {
	dir := "__1__"
	if d.IsDir() {
		dir = "__0__"
	}
	return dir + d.Name()
}
