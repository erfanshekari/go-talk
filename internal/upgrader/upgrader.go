package upgrader

import (
	"log"
	"net/http"
	"sync"

	"github.com/erfanshekari/go-talk/config"
	"github.com/gorilla/websocket"
)

var lock = &sync.Mutex{}

type Upgrader struct {
	Upgrader *websocket.Upgrader
}

var singleInstance *Upgrader

func GetInstance(conf *config.Config) *Upgrader {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			log.Println("**  Initializing gorilla Websocket Upgrader...")
			singleInstance = &Upgrader{
				Upgrader: &websocket.Upgrader{
					HandshakeTimeout: conf.Server.WebSocket.HandshakeTimeout,
					ReadBufferSize:   conf.Server.WebSocket.ReadBufferSize,
					WriteBufferSize:  conf.Server.WebSocket.WriteBufferSize,
					CheckOrigin:      func(r *http.Request) bool { return true },
				},
			}
		}
	}
	return singleInstance
}
