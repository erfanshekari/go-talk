package channel

import (
	"github.com/labstack/echo/v4"
)

func Channel(g *echo.Group) {
	g.GET("/channel/:id", func(c echo.Context) error {

		return c.JSON(200, nil)
	})
}
