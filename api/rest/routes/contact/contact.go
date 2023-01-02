package contact

import (
	"github.com/labstack/echo/v4"
)

func Contact(g *echo.Group) {
	g.GET("/contact/:id", func(c echo.Context) error {
		return c.JSON(200, nil)
	})
}
