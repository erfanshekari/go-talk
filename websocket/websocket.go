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
	Connection         *ws.Conn
	Context            echo.Context
	User               *models.User
	PublicKey          *rsa.PublicKey
	KeyExchangeTimer   *time.Timer
	AuhthenticateTimer *time.Timer
}

func NewConnection(con *ws.Conn, c echo.Context) *WebSocket {
	w := WebSocket{
		Connection: con,
		Context:    c,
	}
	w.StartKeyExchangeTimeout()

	return &w
}

func (w *WebSocket) StartKeyExchangeTimeout() {
	w.KeyExchangeTimer = time.AfterFunc(global.GetInstance(nil).Config.Server.WebSocket.RSAExchangeTimeout, func() {
		log.Println("closing connection key exchange timeout")
		w.Connection.Close()
	})
	log.Println("end of StartKeyExchangeTimeout function")
}

func (w *WebSocket) StartAuthenticateTimeout() {
	w.AuhthenticateTimer = time.AfterFunc(global.GetInstance(nil).Config.Server.WebSocket.AuthenticateTimeout, func() {
		log.Println("closing connection Auth timeout")
		w.Connection.Close()
	})
	log.Println("end of StartAuthenticateTimeout function")
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
			err = w.Connection.WriteMessage(1, responseBytes)
			if err != nil {
				return err
			}
			return nil
		}
	} else {
		return errors.New("PublicKey not exist")
	}
}

func (w *WebSocket) HandleConnection() error {

	mt, msg, err := w.Connection.ReadMessage()
	log.Println(string(msg))
	log.Println("Message Type:", mt)
	if err != nil {
		return err
	}

	// switch mt {
	// case ws.MsgTypeBinary:
	// 	log.Println("message is MsgTypeBinary")
	// case ws.MsgTypeText:
	// 	log.Println("message is MsgTypeText")
	// }

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
		w.KeyExchangeTimer.Stop()
		w.StartAuthenticateTimeout()
		return nil
	} else if w.User == nil {
		log.Println("auth block")
		var response sevents.BytesWrappedJson
		err := json.Unmarshal(msg, &response)
		if err != nil {
			return err
		}
		log.Println(response)
		return nil
	}

	return nil
}
