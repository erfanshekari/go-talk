package websocket

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"log"
	"time"

	"github.com/erfanshekari/go-talk/internal/global"
	"github.com/erfanshekari/go-talk/models"
	ws "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func getServerPublicKey() *rsa.PublicKey {
	return &global.GetInstance(nil).PrivateKey.PublicKey
}

func exportPubKeyAsPEMStr(pubkey *rsa.PublicKey) string {
	pubKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pubkey),
		},
	))
	return pubKeyPem
}

func parseClientPublicKey(pk string) (*rsa.PublicKey, error) {
	log.Println(pk)
	block, _ := pem.Decode([]byte(pk))
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubkey := publicKey.(*rsa.PublicKey)
	return pubkey, nil
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
	PKX   EventType = "pkx"
	Bytes EventType = "bytes"
)

type BytesWrappedJson struct {
	Type    EventType `json:"type"`
	Content string    `json:"content"`
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
	log.Println(string(event))
	if w.PublicKey != nil {
		encryptedBytes, err := rsa.EncryptOAEP(
			sha256.New(),
			rand.Reader,
			w.PublicKey,
			event,
			nil)
		if err != nil {
			return err
		}
		println(string(encryptedBytes))
		event := BytesWrappedJson{
			Type:    Bytes,
			Content: string(encryptedBytes),
		}
		eventAsBytes, err := json.Marshal(event)
		if err != nil {
			return err
		}
		w.Connection.WriteMessage(1, eventAsBytes)
	} else {
		return errors.New("PublicKey not exist")
	}

	return nil
}

func (w *WebSocket) HandleConnection() error {

	_, msg, err := w.Connection.ReadMessage()
	// log.Println(mt, string(msg))
	if err != nil {
		return err
	}

	if w.PublicKey == nil {
		var clientPublicKeyEvent ClientPublicKeyEvent
		err := json.Unmarshal(msg, &clientPublicKeyEvent)
		if err != nil {
			return err
		}
		pubkey, err := parseClientPublicKey(clientPublicKeyEvent.PublicKey)
		log.Println(pubkey)
		if err != nil {
			return err
		}
		w.PublicKey = pubkey
		// serverPublicKeyEvent := ServerPublicKeyEvent{
		// 	PublicKey: exportPubKeyAsPEMStr(getServerPublicKey()),
		// }
		// _, err := json.Marshal(serverPublicKeyEvent)
		// if err != nil {
		// 	return err
		// }
		err = w.Send([]byte("hello bitch"))
		if err != nil {
			return err
		}
	} else {

	}

	return nil
}
