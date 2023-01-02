package contacts

import (
	"github.com/labstack/echo/v4"
)

func Contacts(g *echo.Group) {
	g.GET("/contacts", func(c echo.Context) error {
		return c.JSON(200, nil)
	})
}
