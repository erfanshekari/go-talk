package cli

import (
	"fmt"
	"log"
	"strconv"

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
	fmt.Println("/////////////.......................................................................................")
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
	fmt.Println("=====> migrate                                   : Migrate all models to database")
	fmt.Println("         *Options:")
	fmt.Println("           -db [database-name]                   : Specify database name, default=go-talk")
	fmt.Println("=====> --help , -h, help                         : Print available commands")
	fmt.Println("=====> -v , --version, version                   : Getting current version of GoTalk")
	fmt.Println("////////////////////////////////////////////////////////////////////////////////////////////////////")
}

func PrintVersion(v string) {
	fmt.Println("Gotalk Server Version", v)
}

func ExportConfig(args []string) (*config.ConfigAtrs, error) {
	var ip string
	var port int32
	var debug bool
	var debugLazy bool
	var databaseName string
	for index, value := range args {
		switch value {
		case "-p":
			target := index + 1
			if (target + 1) <= len(args) {
				target, err := strconv.ParseInt(args[target], 10, 32)
				if err != nil {
					return nil, err
				} else {
					port = int32(target)
				}
			}
		case "-b":
			target := index + 1
			if (target + 1) <= len(args) {
				ip = args[index+1]
			}
		case "-d":
			debug = true
			target := index + 1
			if (target + 1) <= len(args) {
				value = args[index+1]
				if value == "lazy" || value == "L" || value == "l" {
					debugLazy = true
				}
			}
		case "-db":
			target := index + 1
			if (target + 1) <= len(args) {
				value = args[index+1]
				if value != "" {
					databaseName = value
				}
			}
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
		conf, err := ExportConfig(args)
		if err != nil {
			return "", nil
		} else {
			return RunServer, conf
		}
	case "-h", "--help", "help":
		return Help, nil
	case "migrate":
		conf, err := ExportConfig(args)
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
