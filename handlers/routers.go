package handlers

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Routes(htmls embed.FS, baseDir string) (*echo.Echo, error) {
	readBaseDir = baseDir

	e := echo.New()
	e.Use(middleware.Logger())
	e.HideBanner = true
	public, err := fs.Sub(htmls, "static")
	if err != nil {
		return nil, err
	}
	// _ = public
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(public))))
	e.GET("/d", listHandler)
	e.GET("/d/", listHandler)
	e.GET("/d/:dirname", listHandler)
	e.GET("/c/:filepath", handler)
	e.GET("/i/:filepath", infoHandler)
	return e, nil
}
