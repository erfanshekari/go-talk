package routes

import (
	"github.com/erfanshekari/go-talk/api/message"
	"github.com/labstack/echo/v4"
)

type RouteGroupRegistering func(e *echo.Echo)

var (
	Routes []RouteGroupRegistering = []RouteGroupRegistering{
		message.Message,
	}
)