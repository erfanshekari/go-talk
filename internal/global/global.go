package global

import (
	"sync"

	"github.com/erfanshekari/go-talk/config"
)

var lock = &sync.Mutex{}

type Global struct {
	Config *config.Config
}

var singleInstance *Global

func GetInstance(conf *config.Config) *Global {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			singleInstance = &Global{
				Config: conf,
			}
		}
	}
	return singleInstance
}
