package config

type ConfigAtrs struct {
	IP           string
	Port         int32
	Debug        bool
	DebugLazy    bool
	DatabaseName string
}

var Config = ConfigAtrs{
	IP:           "0.0.0.0",
	Port:         8080,
	Debug:        false,
	DebugLazy:    false,
	DatabaseName: "go-talk",
}
