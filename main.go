package main

import (
	"log"
	"os"

	"github.com/erfanshekari/go-talk/internal/cli"
	"github.com/erfanshekari/go-talk/internal/mdbc"
	"github.com/erfanshekari/go-talk/internal/rdb"
	"github.com/erfanshekari/go-talk/server"
)

func main() {

	ok, conf := cli.HandleCommand(os.Args)

	if !ok {
		return
	}

	// init and health check Redis
	ok, err := rdb.GetInstance(conf).Ping()
	if err != nil || !ok {
		log.Fatal(err, "On Redis Client Ping!")
	}
	defer rdb.GetInstance(nil).Close()

	// init and health check MongoDB
	ok, err = mdbc.GetInstance(conf).Ping()
	if err != nil || !ok {
		log.Fatal(err, "On MongoDB Client Ping!")
	}
	defer mdbc.GetInstance(nil).Close()

	server := server.Server{
		Config: conf,
	}

	server.Listen()
}
