package contacts

import (
	"github.com/labstack/echo/v4"
)

type empity struct{}

func Contacts(e *echo.Echo) {
	e.GET("/contacts", func(c echo.Context) error {
		return c.JSON(200, empity{})
	})
}
