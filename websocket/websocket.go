package websocket

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/erfanshekari/go-talk/internal/global"
	"github.com/erfanshekari/go-talk/models"
	uencrypt "github.com/erfanshekari/go-talk/utils/encrypt"
	ws "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func getServerPublicKey() *rsa.PublicKey {
	return &global.GetInstance(nil).PrivateKey.PublicKey
}

type WebSocket struct {
	Connection *ws.Conn
	Context    echo.Context
	User       *models.User
	PublicKey  *rsa.PublicKey
}

type ClientPublicKeyEvent struct {
	PublicKey string `json:"publicKey"`
}

type ServerPublicKeyEvent struct {
	PublicKey string `json:"publicKey"`
}

type EventType string

const (
	Byte      EventType = "byte"
	ByteArray EventType = "byteArray"
)

type BytesWrappedJson struct {
	Type    EventType `json:"type"`
	Content any       `json:"content"`
}

func NewConnection(con *ws.Conn, c echo.Context) *WebSocket {
	ws := WebSocket{
		Connection: con,
		Context:    c,
	}

	go ws.StartRSAExchangeTimeout()

	return &ws
}

func (w *WebSocket) StartRSAExchangeTimeout() {
	time.Sleep(global.GetInstance(nil).Config.Server.WebSocket.RSAExchangeTimeout)
	if w.PublicKey == nil {
		w.Connection.Close()
	}
}

func (w *WebSocket) Send(event []byte) error {
	log.Println("len", len(event))
	if w.PublicKey != nil {
		encryptedChunks, err := uencrypt.Encrypt(event, w.PublicKey)
		if err != nil {
			return err
		}
		if len(encryptedChunks) > 1 {
			var chunks []string
			for _, a := range encryptedChunks {
				chunks = append(chunks, string(a))
			}
			response := BytesWrappedJson{
				Type:    ByteArray,
				Content: chunks,
			}
			responseBytes, err := json.Marshal(response)
			if err != nil {
				return err
			}
			w.Connection.WriteMessage(1, responseBytes)
			return nil
		} else {
			response := BytesWrappedJson{
				Type:    Byte,
				Content: string(encryptedChunks[0]),
			}
			responseBytes, err := json.Marshal(response)
			if err != nil {
				return err
			}
			w.Connection.WriteMessage(1, responseBytes)
			return nil
		}
	} else {
		return errors.New("PublicKey not exist")
	}
}

func (w *WebSocket) HandleConnection() error {

	mt, msg, err := w.Connection.ReadMessage()
	log.Println(mt, string(msg))
	if err != nil {
		return err
	}

	if w.PublicKey == nil {
		var clientPublicKeyEvent ClientPublicKeyEvent
		err := json.Unmarshal(msg, &clientPublicKeyEvent)
		if err != nil {
			return err
		}
		pubkey, err := uencrypt.ParsePublicKey(clientPublicKeyEvent.PublicKey)
		if err != nil {
			return err
		}
		w.PublicKey = pubkey
		serverPublicKeyEvent := ServerPublicKeyEvent{
			PublicKey: uencrypt.ExportPubKeyAsPEMStr(getServerPublicKey()),
		}
		response, err := json.Marshal(serverPublicKeyEvent)
		if err != nil {
			return err
		}
		err = w.Send(response)
		if err != nil {
			return err
		}
	} else {

	}

	return nil
}
