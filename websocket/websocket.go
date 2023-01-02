package websocket

import (
	"log"

	ws "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocket struct {
	Connection *ws.Conn
	Context    echo.Context
}

func NewConnection(con *ws.Conn, c echo.Context) *WebSocket {
	ws := WebSocket{
		Connection: con,
		Context:    c,
	}
	return &ws
}

func (w *WebSocket) HandleConnection() error {
	mt, msg, err := w.Connection.ReadMessage()
	if err != nil {
		return err
	}
	log.Println(mt, msg)
	return nil
}
