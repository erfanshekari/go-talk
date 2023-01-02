package config

import "time"

type ConfigAtrs struct {
	IP                                                string
	Port                                              int32
	Debug                                             bool
	DebugLazy                                         bool
	DatabaseName                                      string
	WebSocketHandshakeTimeout                         time.Duration
	WebSocketReadBufferSize, WebSocketWriteBufferSize int
}

var Config = ConfigAtrs{
	IP:                        "0.0.0.0",
	Port:                      8080,
	Debug:                     false,
	DebugLazy:                 false,
	DatabaseName:              "go-talk",
	WebSocketHandshakeTimeout: time.Second * 3,
	WebSocketReadBufferSize:   1024,
	WebSocketWriteBufferSize:  1024,
}
