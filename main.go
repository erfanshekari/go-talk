package main

import (
	"context"
	"log"
	"os"

	"github.com/erfanshekari/go-talk/config"
	"github.com/erfanshekari/go-talk/internal/cli"
	"github.com/erfanshekari/go-talk/internal/mdbc"
	"github.com/erfanshekari/go-talk/internal/rdb"
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
		initRedis(conf)
		defer rdb.GetInstance(nil).Close()
		initMongoDB(conf)
		defer mdbc.GetInstance(nil).Close()
		server := server.Server{
			Config: conf,
		}
		server.Listen()
	case cli.Migrate:
		log.Println("Migrating Models...")
		conf.Debug = true
		conf.DebugLazy = true
		initMongoDB(conf)
		defer mdbc.GetInstance(nil).Close()
		models.Migrate(context.Background(), conf)
	case cli.Help:
		cli.PrintHelp()
	case cli.Version:
		cli.PrintVersion(version)
	}

}

func initRedis(conf *config.ConfigAtrs) {
	// init and health check Redis
	ok, err := rdb.GetInstance(conf).Ping()
	if err != nil || !ok {
		log.Fatal(err, "On Redis Client Ping!")
	}
}

func initMongoDB(conf *config.ConfigAtrs) {
	// init and health check MongoDB
	ok, err := mdbc.GetInstance(conf).Ping()
	if err != nil || !ok {
		log.Fatal(err, "On MongoDB Client Ping!")
	}
}
