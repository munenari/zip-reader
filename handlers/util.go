package handlers

import (
	"io/fs"
	"path"

	"github.com/munenari/read-zip/util"
)

type (
	DirInfo struct {
		Name       string `json:"name"`
		HashedName string `json:"hashed_name"`
		IsDir      bool   `json:"is_dir"`
	}
)

func newDirInfo(baseDir, name string, isDir bool) (*DirInfo, error) {
	hashedName, err := util.GzipPath(path.Join(baseDir, name))
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
