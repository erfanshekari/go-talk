package channels

import (
	"github.com/labstack/echo/v4"
)

type empity struct{}

func Channels(e *echo.Echo) {
	e.GET("/channels", func(c echo.Context) error {
		return c.JSON(200, empity{})
	})
}
