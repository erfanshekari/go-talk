package routes

import (
	"github.com/erfanshekari/go-talk/api/rest/routes/message"
	"github.com/erfanshekari/go-talk/api/ws"
	"github.com/labstack/echo/v4"
)

type RouteRegistering func(e *echo.Echo)

type RouteGroupRegistering func(e *echo.Group)

// All routes require JWT token in header

var (
	Routes []RouteRegistering = []RouteRegistering{
		ws.WebSocketRoute, // ["/"]
	}
	ProtectedRoutes []RouteGroupRegistering = []RouteGroupRegistering{
		message.Message, // ["/message"]
	}
)
