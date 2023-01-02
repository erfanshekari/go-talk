package ws

import (
	"net/http"

	"github.com/erfanshekari/go-talk/api/ws/handler"
	"github.com/erfanshekari/go-talk/internal/upgrader"
	"github.com/labstack/echo/v4"
)

func WebSocketRoute(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		upgrader := *upgrader.GetInstance(nil).Upgrader
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		handler.HandleWebSocketConnection(ws)
		return nil
	})
}
