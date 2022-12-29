package channels

import (
	"log"

	ctx "github.com/erfanshekari/go-talk/context"
	"github.com/labstack/echo/v4"
)

type empity struct{}

func Channels(e *echo.Echo) {
	e.GET("/channels", func(c echo.Context) error {
		cc := c.(*ctx.Context)
		log.Println(cc.User)
		return c.JSON(200, empity{})
	})
}
