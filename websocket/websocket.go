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
	sevents "github.com/erfanshekari/go-talk/websocket/events"
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

func (w *WebSocket) StartAuthenticateTimeout() {
	time.Sleep(global.GetInstance(nil).Config.Server.WebSocket.AuthenticateTimeout)
	if w.User == nil {
		w.Connection.Close()
	}
}

func (w *WebSocket) Send(event []byte) error {
	if w.PublicKey != nil {
		encryptedChunks, err := uencrypt.Encrypt(event, w.PublicKey)
		if err != nil {
			return err
		}
		if len(encryptedChunks) > 1 {
			var chunks [][]byte
			for _, a := range encryptedChunks {
				chunks = append(chunks, a)
			}
			response := sevents.BytesWrappedJson{
				Type:    sevents.ByteArray,
				Content: chunks,
			}
			responseBytes, err := json.Marshal(response)
			if err != nil {
				return err
			}
			w.Connection.WriteMessage(1, responseBytes)
			return nil
		} else {
			response := sevents.BytesWrappedJson{
				Type:    sevents.Byte,
				Content: encryptedChunks[0],
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
		var clientPublicKey sevents.ClientPublicKey
		err := json.Unmarshal(msg, &clientPublicKey)
		if err != nil {
			return err
		}
		pubkey, err := uencrypt.ParsePublicKey(clientPublicKey.PublicKey)
		if err != nil {
			return err
		}
		w.PublicKey = pubkey
		serverPublicKey := sevents.ServerPublicKey{
			PublicKey: uencrypt.ExportPubKeyAsPEMStr(getServerPublicKey()),
		}
		response, err := json.Marshal(serverPublicKey)
		if err != nil {
			return err
		}
		err = w.Send(response)
		if err != nil {
			return err
		}
		w.StartAuthenticateTimeout()
		return nil
	} else if w.User == nil {
		var response sevents.BytesWrappedJsonType
		err := json.Unmarshal(msg, &response)
		if err != nil {
			return err
		}
		log.Println(response)
	}

	return nil
}
