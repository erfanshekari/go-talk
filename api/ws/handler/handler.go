package handler

import (
	"log"

	"github.com/gorilla/websocket"
)

func HandleWebSocketConnection(con *websocket.Conn) {
	for {
		mt, m, err := con.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		log.Println(mt, m)
	}
}
