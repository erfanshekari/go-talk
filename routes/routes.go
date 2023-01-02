package routes

import (
	"github.com/erfanshekari/go-talk/api/rest/routes/message"
	"github.com/erfanshekari/go-talk/api/ws"
	"github.com/labstack/echo/v4"
)

type RouteGroupRegistering func(e *echo.Echo)

// All routes require JWT token in header

var (
	Routes []RouteGroupRegistering = []RouteGroupRegistering{
		ws.WebSocketRoute, // ["/"]
		message.Message,   // ["/message"]
	}
)
