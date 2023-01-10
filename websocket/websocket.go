package websocket

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/erfanshekari/go-talk/internal/global"
	"github.com/erfanshekari/go-talk/models"
	uencrypt "github.com/erfanshekari/go-talk/utils/encrypt"
	ujwt "github.com/erfanshekari/go-talk/utils/jwt"
	sevents "github.com/erfanshekari/go-talk/websocket/events"

	ws "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func getServerPublicKey() *rsa.PublicKey {
	return &global.GetInstance(nil).PrivateKey.PublicKey
}
func NewConnection(con *ws.Conn, c echo.Context) *WebSocket {
	w := WebSocket{
		Connection: con,
		Context:    c,
	}
	w.StartKeyExchangeTimeout()

	return &w
}

type WebSocket struct {
	Connection         *ws.Conn
	Context            echo.Context
	User               *models.User
	PublicKey          *rsa.PublicKey
	KeyExchangeTimer   *time.Timer
	AuhthenticateTimer *time.Timer
}

func (w *WebSocket) KeyExchanged() bool {
	return w.PublicKey != nil
}

func (w *WebSocket) Authenticated() bool {
	return w.User != nil
}

func (w *WebSocket) Authenticate(t sevents.ClientJWTToken) (*models.User, error) {
	token, err := ujwt.GetToken(t.AccessToken)
	if err != nil {
		return nil, err
	}
	user, err := ujwt.GetUserFromJWT(token)
	if err != nil {
		return nil, err
	}
	w.User = user
	w.AuhthenticateTimer.Stop()
	return w.User, nil
}

func (w *WebSocket) StartKeyExchangeTimeout() {
	w.KeyExchangeTimer = time.AfterFunc(global.GetInstance(nil).Config.Server.WebSocket.RSAExchangeTimeout, func() {
		w.Connection.Close()
	})
}

func (w *WebSocket) StartAuthenticateTimeout() {
	w.AuhthenticateTimer = time.AfterFunc(global.GetInstance(nil).Config.Server.WebSocket.AuthenticateTimeout, func() {
		w.Connection.Close()
	})
}

func (w *WebSocket) Send(event []byte) error {
	if w.KeyExchanged() {
		encryptedJson, err := uencrypt.Encrypt(event, w.PublicKey)
		if err != nil {
			return err
		}
		responseBytes, err := json.Marshal(encryptedJson)
		if err != nil {
			return err
		}
		err = w.Connection.WriteMessage(ws.TextMessage, responseBytes)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Can't send data before key exchange...")
}

func (w *WebSocket) HandleConnection() error {

	mt, msg, err := w.Connection.ReadMessage()
	if err != nil {
		return err
	}
	log.Println(mt, string(msg))
	switch mt {
	case ws.TextMessage:
		log.Println("message type is safe")
		break
	default:
		log.Println("Unhandler message type we sould close connection")
		return errors.New(fmt.Sprintf("Can't Handle MessageType: %d", mt))
	}

	if w.KeyExchanged() && w.Authenticated() {
		log.Println("Lets handle commiunication somewhere else...")
	}

	if !w.KeyExchanged() {
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
	} else if !w.Authenticated() {
		authMsg, err := uencrypt.Decrypt(msg, global.GetInstance(nil).PrivateKey)
		if err != nil {
			return err
		}
		var token sevents.ClientJWTToken
		err = json.Unmarshal(*authMsg, &token)
		if err != nil {
			return err
		}
		user, err := w.Authenticate(token)
		if err != nil {
			return w.Connection.Close()
		}
		userACK := sevents.ServerUserACK{UserID: user.UserID}
		userACKByte, err := json.Marshal(userACK)
		if err != nil {
			return err
		}
		return w.Send(userACKByte)
	}

	return nil
}
