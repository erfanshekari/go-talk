package message

import (
	"github.com/labstack/echo/v4"
)

type empity struct{}

func Message(e *echo.Echo) {
	e.POST("/message", func(c echo.Context) error {
		e.Logger.Print(c.Get("user"))
		return c.JSON(200, empity{})
	})
}
