package global

import (
	"crypto/rsa"
	"log"
	"os"
	"sync"

	"github.com/erfanshekari/go-talk/config"
	uencrypt "github.com/erfanshekari/go-talk/utils/encrypt"
	"github.com/joho/godotenv"
)

var lock = &sync.Mutex{}

type Global struct {
	Config     *config.Config
	PrivateKey *rsa.PrivateKey
	SecretKey  string
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
			if conf != nil && conf.Server.Debug {
				err := godotenv.Load()
				if err != nil {
					log.Fatal(err)
				}
			}
			secretKey := os.Getenv("SECRET_KEY")
			if secretKey == "" {
				log.Fatal("SECRET_KEY is empity, you must provide a SECRET_KEY env var...")
			}
			singleInstance = &Global{
				Config:     conf,
				PrivateKey: privateKey,
				SecretKey:  secretKey,
			}
		}
	}
	return singleInstance
}
