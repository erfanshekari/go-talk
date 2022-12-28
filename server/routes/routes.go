package routes

import (
	"github.com/erfanshekari/go-talk/server/routes/message"
	"github.com/labstack/echo/v4"
)

type RouteGroupRegistering func(e *echo.Echo)

var (
	Routes []RouteGroupRegistering = []RouteGroupRegistering{
		message.Message,
	}
)

func RegisterRoutes(e *echo.Echo) {
	for _, route := range Routes {
		route(e)
	}
}
