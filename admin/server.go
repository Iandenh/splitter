package admin

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"splitter/listener"
	"splitter/upstream"
	"strconv"
)

func Start(adminPort int) {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hallo")
	})
	e.GET("/requests", func(c echo.Context) error {
		return c.JSON(http.StatusOK, listener.GetRequests())
	})
	e.GET("/upstreams", func(c echo.Context) error {
		return c.JSON(http.StatusOK, upstream.GetUpstreams())
	})

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(adminPort)))
}
