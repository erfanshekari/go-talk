package message

import (
	"log"

	"github.com/erfanshekari/go-talk/api/rest/types"
	"github.com/erfanshekari/go-talk/context"
	"github.com/labstack/echo/v4"
)

func Message(g *echo.Group) {
	g.POST("/message", func(c echo.Context) error {
		cc := c.(*context.Context)
		log.Println(cc.Get("user"))
		message := new(types.Message)
		if err := c.Bind(message); err != nil {
			log.Println("bind error:", err)
			return c.JSON(200, types.Empity{})
		}
		if err := c.Validate(message); err != nil {
			log.Println("err:", err)
			return c.JSON(200, types.Empity{})
		}
		log.Println(message)
		return c.JSON(200, types.Empity{})
	})
}
