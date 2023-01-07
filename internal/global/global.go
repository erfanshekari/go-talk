package global

import (
	"crypto/rsa"
	"log"
	"sync"

	"github.com/erfanshekari/go-talk/config"
	uencrypt "github.com/erfanshekari/go-talk/utils/encrypt"
)

var lock = &sync.Mutex{}

type Global struct {
	Config     *config.Config
	PrivateKey *rsa.PrivateKey
}

var singleInstance *Global

func GetInstance(conf *config.Config) *Global {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			privateKey, err := uencrypt.GenerateKey()
			if err != nil {
				log.Fatal(err)
			}
			singleInstance = &Global{
				Config:     conf,
				PrivateKey: privateKey,
			}
		}
	}
	return singleInstance
}
