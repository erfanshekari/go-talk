package cli

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/erfanshekari/go-talk/config"
)

type Command string

const (
	RunServer Command = "run"
	Migrate   Command = "migrate"
	Help      Command = "help"
	Version   Command = "version"
)

func PrintHelp() {
	fmt.Println("GoTalk Server")
	fmt.Println("Commands:")
	fmt.Println("=====> run                                       : Run GoTalk server")
	fmt.Println("         *Options:")
	fmt.Println("           -b [ip]                               : Bind server to given ip address, default=0.0.0.0")
	fmt.Println("           -p [port]                             : Specify the connection port, default=8080")
	fmt.Println("           -db [database-name]                   : Specify database name, default=go-talk")
	fmt.Println("           -d                                    : enable debug mode")
	fmt.Println("               *Options:")
	fmt.Println("                   1.leave empity will build tests")
	fmt.Println("                   2.\"l\"=(lazy debug run) mean's only run server without building tests")
	fmt.Println("                   3.\"lazy\"=equal to option (1)")
	fmt.Println("           -ws                                   : Customize websocket configuration")
	fmt.Println("               *Options:")
	fmt.Println("                   h=3                           : Specify the duration of the WebSocket handshake in seconds")
	fmt.Println("                   rb=1024                       : Specify read buffer size in bytes")
	fmt.Println("                   wb=1024                       : Specify write buffer size in bytes")
	fmt.Println("=====> migrate                                   : Migrate all models to database")
	fmt.Println("         *Options:")
	fmt.Println("           -db [database-name]                   : Specify database name, default=go-talk")
	fmt.Println("=====> --help , -h, help                         : Print available commands")
	fmt.Println("=====> -v , --version, version                   : Getting current version of GoTalk")
}

func PrintVersion(v string) {
	fmt.Println("Gotalk Server Version", v)
}

func parseInt(val string) (int64, bool) {
	value, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0, false
	} else {
		return value, true
	}
}

func exportWebSocketConfig(args []string, conf *config.Config) {
	for _, val := range args {
		split := strings.Split(val, "=")
		if len(split) == 2 {
			switch split[0] {
			case "h":
				val, ok := parseInt(split[1])
				if !ok {
					continue
				}
				conf.Server.WebSocket.HandshakeTimeout = time.Duration(val) * time.Second
			case "rb":
				val, ok := parseInt(split[1])
				if !ok {
					continue
				}
				conf.Server.WebSocket.ReadBufferSize = int(val)
			case "wb":
				val, ok := parseInt(split[1])
				if !ok {
					continue
				}
				conf.Server.WebSocket.WriteBufferSize = int(val)
			}
		}
	}
}

func exportConfig(args []string) (*config.Config, error) {
	defaultConfig := config.GetDefaultConfig()
	lenOfArgs := len(args)
	for index, value := range args {
		switch value {
		case "-p":
			target := index + 1
			if (target + 1) <= lenOfArgs {
				target, err := strconv.ParseInt(args[target], 10, 32)
				if err != nil {
					return nil, err
				} else {
					defaultConfig.Server.Port = int(target)
				}
			}
		case "-b":
			target := index + 1
			if (target + 1) <= lenOfArgs {
				defaultConfig.Server.Host = args[index+1]
			}
		case "-d":
			defaultConfig.Server.Debug = true
			target := index + 1
			if (target + 1) <= lenOfArgs {
				value = args[index+1]
				if value == "lazy" || value == "L" || value == "l" {
					defaultConfig.Server.LazyDebug = true
				}
			}
		case "-db":
			target := index + 1
			if (target + 1) <= lenOfArgs {
				value = args[index+1]
				if value != "" {
					defaultConfig.Database.Name = value
				}
			}
		case "-ws":
			target := index + 1
			var possibleTargets []string
			for i := target; i < (target + 3); i++ {
				if (i + 1) <= lenOfArgs {
					possibleTargets = append(possibleTargets, args[i])
				}
			}
			exportWebSocketConfig(possibleTargets, defaultConfig)
		}
	}

	return defaultConfig, nil
}

func HandleCommand(args []string) (Command, *config.Config) {
	var command string
	if len(args) < 2 {
		command = "help"
	} else {
		command = args[1]
	}
	switch command {
	case "run", "Run":
		conf, err := exportConfig(args)
		if err != nil {
			return "", nil
		} else {
			return RunServer, conf
		}
	case "-h", "--help", "help":
		return Help, nil
	case "migrate":
		conf, err := exportConfig(args)
		if err != nil {
			log.Fatal(err)
		}
		return Migrate, conf
	case "-v", "--version", "version":
		return Version, nil
	default:
		return Help, nil
	}
}
