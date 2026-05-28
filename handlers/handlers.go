package handlers

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/munenari/read-zip/optimize"
	"github.com/munenari/read-zip/util"
)

var (
	readBaseDir string
)

func handler(c echo.Context) error {
	filepath := c.Param("filepath")
	fp, err := util.UnGzipPath(filepath)
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
	rc, closer, err := recursizeReadPage(c.Request().Context(), fp, page)
	if err != nil {
		return err
	}
	defer closer.Close()
	defer rc.Close()
	if err := c.Request().Context().Err(); err != nil {
		return err
	}
	c.Response().Header().Set("Cache-Control", "max-age=86400")
	c.Response().Header().Set("Content-Disposition", "inline")
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().WriteHeader(http.StatusOK)
	return optimize.ResizeToMax(c.Request().Context(), rc, c.Response())
}

func recursizeReadPage(ctx context.Context, fp string, page int) (rc io.ReadCloser, closer io.Closer, err error) {
	zr, err := util.NewZipReaderWithOpen(path.Join(readBaseDir, fp))
	if err == nil {
		rc, err = zr.At(page)
		if err == nil {
			return rc, zr, nil
		}
		zr.Close()
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	dirEntries, err := os.ReadDir(path.Join(readBaseDir, fp))
	if err != nil {
		return nil, nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to read dir").SetInternal(err)
	}
	digSize := min(len(dirEntries), 4)
	for i := range digSize {
		if err := ctx.Err(); err != nil {
			return nil, nil, err
		}
		zr, err = util.NewZipReaderWithOpen(path.Join(readBaseDir, fp, dirEntries[i].Name()))
		if err != nil {
			continue
		}
		rc, err = zr.At(page)
		if err == nil {
			return rc, zr, nil
		}
		zr.Close()
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	return nil, nil, echo.NewHTTPError(http.StatusBadRequest, "failed to read page").SetInternal(err)
}

func listHandler(c echo.Context) error {
	dirname, err := util.UnGzipPath(c.Param("dirname"))
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
	fp, err := util.UnGzipPath(filepath)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "filepath param was invalid").SetInternal(err)
	}
	zr, err := util.NewZipReaderWithOpen(path.Join(readBaseDir, fp))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "page was invalid").SetInternal(err)
	}
	defer zr.Close()
	dirEntries, err := os.ReadDir(path.Join(readBaseDir, path.Dir(fp)))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "path was invalid").SetInternal(err)
	}
	dirEntries = slices.DeleteFunc(dirEntries, func(v fs.DirEntry) bool {
		return strings.HasPrefix(v.Name(), ".")
	})
	currentIndex := slices.IndexFunc(dirEntries, func(v os.DirEntry) bool {
		if strings.Index(v.Name(), ".") == 0 {
			return false
		}
		return path.Join(path.Dir(fp), v.Name()) == fp
	})
	prevIndex := currentIndex - 1
	if prevIndex < 0 {
		prevIndex = len(dirEntries) - 1
	}
	nextIndex := currentIndex + 1
	if nextIndex > len(dirEntries)-1 {
		nextIndex = 0
	}
	prevHashedName, err := util.GzipPath(path.Join(path.Dir(fp), dirEntries[prevIndex].Name()))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to encode filename").SetInternal(err)
	}
	nextHashedName, err := util.GzipPath(path.Join(path.Dir(fp), dirEntries[nextIndex].Name()))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to encode filename").SetInternal(err)
	}
	parentDir, err := util.GzipPath(path.Dir(fp))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to encode filename").SetInternal(err)
	}
	v := map[string]interface{}{
		"name":               fp,
		"size":               zr.Len(),
		"prev_hashed_name":   prevHashedName,
		"next_hashed_name":   nextHashedName,
		"parent_hashed_name": parentDir,
	}
	return c.JSON(http.StatusOK, v)
}
