package channels

import (
	"github.com/labstack/echo/v4"
)

func Channels(g *echo.Group) {
	g.GET("/channels", func(c echo.Context) error {
		return c.JSON(200, nil)
	})
}
