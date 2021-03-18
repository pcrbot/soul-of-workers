package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"

	"github.com/pcrbot/soul-of-workers/game"
)

func main() {
	cmd := parseCommand()
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableQuote:    true,
		FullTimestamp:   true,
		TimestampFormat: "01/02 15:04:05",
	})
	switch cmd.action {
	case "run":
		run()
	default:
		fmt.Println("unknown command " + cmd.action)
		os.Exit(1)
	}
}

func run() {
	conf, err := readConfig()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if err = game.InitialDB(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if err = game.RegisterCommands(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	if err=game.InitialCron();err!=nil{
		log.Error(err)
		os.Exit(1)
	}

	zero.Run(zero.Config{
		NickName:      conf.ZeroBot.NickName,
		CommandPrefix: conf.ZeroBot.CommandPrefix,
		SuperUsers:    conf.ZeroBot.SuperUsers,
		Driver: []zero.Driver{
			driver.NewWebSocketClient(conf.Host, conf.Port, conf.AccessToken),
		},
	})
	select {}
}
