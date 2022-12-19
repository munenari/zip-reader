package handlers

import (
	"bytes"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/munenari/read-zip/optimize"
)

var (
	readBaseDir string
)

func handler(c echo.Context) error {
	filepath := c.Param("filepath")
	fp, err := ungzipPath(filepath)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "filepath param was invalid").SetInternal(err)
	}
	pageStr := c.QueryParam("page")
	if pageStr == "" {
		pageStr = "0"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "page was invalid").SetInternal(err)
	}
	b, err := readPage(path.Join(readBaseDir, fp), page)
	if err != nil {
		dirEntries, err := os.ReadDir(path.Join(readBaseDir, fp))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
		}
		b, err = readPage(path.Join(readBaseDir, fp, dirEntries[0].Name()), page)
		if err != nil {
			b, err = readPage(path.Join(readBaseDir, fp, dirEntries[1].Name()), page)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "failed to read page").SetInternal(err)
			}
		}
	}
	c.Response().Header().Add("Cache-Control", "max-age=86400")
	buf := &bytes.Buffer{}
	if err := optimize.ResizeToMax(bytes.NewReader(b), buf); err != nil {
		c.Logger().Info(err)
		return c.Blob(http.StatusOK, http.DetectContentType(b), b)
	}
	return c.Stream(http.StatusOK, http.DetectContentType(buf.Bytes()), buf)
}

func listHandler(c echo.Context) error {
	dirname, err := ungzipPath(c.Param("dirname"))
	if err != nil {
		c.Logger().Error(err)
	}
	dirEntries, err := os.ReadDir(path.Join(readBaseDir, dirname))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	sort.Slice(dirEntries, func(i, j int) bool {
		nameI := dirAndFileName(dirEntries[i])
		nameJ := dirAndFileName(dirEntries[j])
		return nameI < nameJ
	})
	pa, err := newDirInfo(dirname, "..", true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	dirInfos := []*DirInfo{pa}
	for _, d := range dirEntries {
		if strings.Index(d.Name(), ".") == 0 {
			continue
		}
		di, err := newDirInfo(dirname, d.Name(), d.IsDir())
		if err != nil {
			continue
		}
		dirInfos = append(dirInfos, di)
	}
	return c.JSON(http.StatusOK, dirInfos)
}

func infoHandler(c echo.Context) error {
	filepath := c.Param("filepath")
	fp, err := ungzipPath(filepath)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "filepath param was invalid").SetInternal(err)
	}
	l, err := length(path.Join(readBaseDir, fp))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "page was invalid").SetInternal(err)
	}
	dirEntries, err := os.ReadDir(path.Join(readBaseDir, path.Dir(fp)))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "path was invalid").SetInternal(err)
	}
	currentIndex := 0
	for i, de := range dirEntries {
		if strings.Index(de.Name(), ".") == 0 {
			continue
		}
		if path.Join(path.Dir(fp), de.Name()) == fp {
			currentIndex = i
			break
		}
	}
	prevIndex := currentIndex - 1
	if prevIndex < 0 {
		prevIndex = len(dirEntries) - 1
	}
	nextIndex := currentIndex + 1
	if nextIndex > len(dirEntries)-1 {
		nextIndex = 0
	}
	prevHashedName, err := gzipPath(path.Join(path.Dir(fp), dirEntries[prevIndex].Name()))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to encode filename").SetInternal(err)
	}
	nextHashedName, err := gzipPath(path.Join(path.Dir(fp), dirEntries[nextIndex].Name()))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to encode filename").SetInternal(err)
	}
	parentDir, err := gzipPath(path.Dir(fp))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to encode filename").SetInternal(err)
	}
	v := map[string]interface{}{
		"name":               fp,
		"size":               l,
		"prev_hashed_name":   prevHashedName,
		"next_hashed_name":   nextHashedName,
		"parent_hashed_name": parentDir,
	}
	return c.JSON(http.StatusOK, v)
}
