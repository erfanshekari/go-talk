package cli

import (
	"fmt"
	"strconv"

	"github.com/erfanshekari/go-talk/config"
)

func PrintHelp() {
	fmt.Println("")
	fmt.Println("GoChat Server Commands")
	fmt.Println("")
	fmt.Println("-> run                        :   Run GoChat server")
	fmt.Println("-> --help , -h, help          :   Print available commands")
	fmt.Println("-> -v , --version, version    :   Getting current version of GoChat")
	fmt.Println("-------------------------------------------------------------------")
	fmt.Println("Flags")
	fmt.Println("-b [ip]                       :    Bind server to given ip address")
	fmt.Println("-p [port]                     :    set port")
	fmt.Println("-d options(\"\", \"lazy\")    :    enable debug mode")
	fmt.Println("-------------------------------------------------------------------")
	fmt.Println("")
}

func ExportConfig(args []string) (*config.ConfigAtrs, error) {
	var ip string
	var port int32
	var debug bool
	var debugLazy bool
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
		}
	}
	conf := config.ConfigAtrs{
		IP:        ip,
		Port:      port,
		Debug:     debug,
		DebugLazy: debugLazy,
	}
	if conf.IP == "" {
		conf.IP = config.Config.IP
	}
	if conf.Port == 0 {
		conf.Port = config.Config.Port
	}
	return &conf, nil
}

func HandleCommand(args []string) (bool, *config.ConfigAtrs) {
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
			return false, nil
		} else {
			return true, conf
		}
	case "-h", "--help", "help":
		PrintHelp()
		return false, nil
	default:
		PrintHelp()
		return false, nil
	}
}
