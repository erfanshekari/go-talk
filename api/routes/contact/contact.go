package contact

import (
	"github.com/labstack/echo/v4"
)

type empity struct{}

func Contact(e *echo.Echo) {
	e.GET("/contact/:id", func(c echo.Context) error {
		return c.JSON(200, empity{})
	})
}
