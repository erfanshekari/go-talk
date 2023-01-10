package ws

import (
	"log"
	"net/http"

	"github.com/erfanshekari/go-talk/internal/upgrader"
	"github.com/erfanshekari/go-talk/websocket"
	"github.com/labstack/echo/v4"
)

func WebSocketRoute(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {

		upgrader := *upgrader.GetInstance(nil).Upgrader

		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		defer ws.Close()

		w := websocket.NewConnection(ws, c)

		for {
			err := w.HandleConnection()
			if err != nil {
				log.Println(err)
				break
			}
			log.Println("looping...")
		}

		return nil
	})
}
