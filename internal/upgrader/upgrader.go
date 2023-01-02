package upgrader

import (
	"log"
	"sync"

	"github.com/erfanshekari/go-talk/config"
	"github.com/gorilla/websocket"
)

var lock = &sync.Mutex{}

type Upgrader struct {
	Upgrader *websocket.Upgrader
}

var singleInstance *Upgrader

func GetInstance(conf *config.ConfigAtrs) *Upgrader {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			log.Println("**  Initializing Gorilla Websocket Upgrader...")
			singleInstance = &Upgrader{
				Upgrader: &websocket.Upgrader{
					HandshakeTimeout: conf.WebSocketHandshakeTimeout,
					ReadBufferSize:   conf.WebSocketReadBufferSize,
					WriteBufferSize:  conf.WebSocketWriteBufferSize,
				},
			}
		}
	}
	return singleInstance
}
