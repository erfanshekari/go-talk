package main

import (
	"context"
	"log"
	"os"

	"github.com/erfanshekari/go-talk/config"
	"github.com/erfanshekari/go-talk/internal/cli"
	"github.com/erfanshekari/go-talk/internal/global"
	"github.com/erfanshekari/go-talk/internal/mdbc"
	"github.com/erfanshekari/go-talk/internal/rdb"
	"github.com/erfanshekari/go-talk/internal/upgrader"
	"github.com/erfanshekari/go-talk/models"
	"github.com/erfanshekari/go-talk/server"
)

const (
	version string = "0.0.1"
)

func main() {

	command, conf := cli.HandleCommand(os.Args)

	switch command {

	case cli.RunServer:
		initGlobalConfig(conf)
		initRedis(conf)
		defer rdb.GetInstance(nil).Close()
		initMongoDB(conf)
		defer mdbc.GetInstance(nil).Close()
		initWebSocketUpgrader(conf)
		server := server.Server{
			Config: conf,
		}
		server.Listen()

	case cli.Migrate:
		log.Println("Migrating Models...")
		initGlobalConfig(conf)
		conf.Server.Debug = true
		conf.Server.LazyDebug = true
		initMongoDB(conf)
		defer mdbc.GetInstance(nil).Close()
		models.Migrate(context.Background(), conf)

	case cli.Help:
		cli.PrintHelp()

	case cli.Version:
		cli.PrintVersion(version)
	}

}

func initRedis(conf *config.Config) {
	// init and health check Redis
	ok, err := rdb.GetInstance(conf).Ping()
	if err != nil || !ok {
		log.Fatal(err, "On Redis Client Ping!")
	}
}

func initMongoDB(conf *config.Config) {
	// init and health check MongoDB
	ok, err := mdbc.GetInstance(conf).Ping()
	if err != nil || !ok {
		log.Fatal(err, "On MongoDB Client Ping!")
	}
}

func initWebSocketUpgrader(conf *config.Config) {
	// init gorilla websocket upgrader
	upgrader.GetInstance(conf)
}

func initGlobalConfig(conf *config.Config) {
	// init Config For Global Access
	global.GetInstance(conf)
}
