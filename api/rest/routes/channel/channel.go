package channel

import (
	"github.com/labstack/echo/v4"
)

type empity struct{}

func Channel(e *echo.Echo) {
	e.GET("/channel/:id", func(c echo.Context) error {

		return c.JSON(200, empity{})
	})
}
