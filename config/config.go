package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		Debug     bool   `yaml:"debug"`
		LazyDebug bool   `yaml:"lazy-debug"`
		WebSocket struct {
			HandshakeTimeout    time.Duration `yaml:"handshake-timeout"`
			ReadBufferSize      int           `yaml:"read-buffer-size"`
			WriteBufferSize     int           `yaml:"write-buffer-size"`
			RSAExchangeTimeout  time.Duration `yaml:"rsa-exchange-timeout"`
			AuthenticateTimeout time.Duration `yaml:"authenticate-timeout"`
		} `yaml:"websocket"`
	} `yaml:"server"`
	Database struct {
		Name string `yaml:"name"`
	} `yaml:"database"`
}

func GetDefaultConfig() *Config {
	f, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var conf Config
	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		log.Fatal(err)
	}
	return &conf
}
