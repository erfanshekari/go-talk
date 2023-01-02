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

type webSocketConfig struct {
	HandshakeTimeout                *time.Duration
	ReadBufferSize, WriteBufferSize *int
}

func parseInt(val string) (int64, bool) {
	value, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0, false
	} else {
		return value, true
	}
}

func exportWebSocketConfig(args []string) webSocketConfig {
	conf := webSocketConfig{}
	for _, val := range args {
		split := strings.Split(val, "=")
		if len(split) == 2 {
			switch split[0] {
			case "h":
				val, ok := parseInt(split[1])
				if !ok {
					continue
				}
				duration := time.Duration(val) * time.Second
				conf.HandshakeTimeout = &duration
			case "rb":
				val, ok := parseInt(split[1])
				if !ok {
					continue
				}
				rb := int(val)
				conf.ReadBufferSize = &rb
			case "wb":
				val, ok := parseInt(split[1])
				if !ok {
					continue
				}
				wb := int(val)
				conf.WriteBufferSize = &wb
			}
		}
	}
	return conf
}

func exportConfig(args []string) (*config.ConfigAtrs, error) {
	var ip string
	var port int32
	var debug bool
	var debugLazy bool
	var databaseName string
	var webSocketConf webSocketConfig
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
					port = int32(target)
				}
			}
		case "-b":
			target := index + 1
			if (target + 1) <= lenOfArgs {
				ip = args[index+1]
			}
		case "-d":
			debug = true
			target := index + 1
			if (target + 1) <= lenOfArgs {
				value = args[index+1]
				if value == "lazy" || value == "L" || value == "l" {
					debugLazy = true
				}
			}
		case "-db":
			target := index + 1
			if (target + 1) <= lenOfArgs {
				value = args[index+1]
				if value != "" {
					databaseName = value
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
			webSocketConf = exportWebSocketConfig(possibleTargets)
		}
	}
	conf := config.ConfigAtrs{
		IP:           ip,
		Port:         port,
		Debug:        debug,
		DebugLazy:    debugLazy,
		DatabaseName: config.Config.DatabaseName,
	}
	if conf.IP == "" {
		conf.IP = config.Config.IP
	}
	if conf.Port == 0 {
		conf.Port = config.Config.Port
	}
	if databaseName != "" {
		conf.DatabaseName = databaseName
	}

	if webSocketConf.HandshakeTimeout != nil {
		conf.WebSocketHandshakeTimeout = *webSocketConf.HandshakeTimeout
	}
	if webSocketConf.ReadBufferSize != nil {
		conf.WebSocketReadBufferSize = *webSocketConf.ReadBufferSize
	}
	if webSocketConf.WriteBufferSize != nil {
		conf.WebSocketWriteBufferSize = *webSocketConf.WriteBufferSize
	}

	return &conf, nil
}

func HandleCommand(args []string) (Command, *config.ConfigAtrs) {
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
