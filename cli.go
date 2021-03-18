package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/alecthomas/kingpin"
)

type zeroConfig struct {
	NickName      []string
	CommandPrefix string
	SuperUsers    []string
}

type Config struct {
	Host        string
	Port        string
	AccessToken string
	LogLevel    string
	ZeroBot     zeroConfig
}

type commandOptions struct {
	action string
}

//go:embed default_config.toml.tmpl
var defaultConfig string

func (c *Config) Save() error {
	tmpl, err := template.New("config").Parse(defaultConfig)
	if err != nil {
		return err
	}
	file, err := os.Create("config.toml")
	if err != nil {
		return err
	}
	err = tmpl.Execute(file, c)
	if err != nil {
		return err
	}
	return nil
}

func parseCommand() commandOptions {
	var cmd commandOptions
	app := kingpin.New("SOW", "SOW the chat game")
	app.Command("run", "start run SOW").Default()
	app.Version("SOW: soul-of-workers开发中")
	app.VersionFlag.Short('V')
	app.HelpFlag.Short('h')
	cmd.action = kingpin.MustParse(app.Parse(os.Args[1:]))
	return cmd
}

func initialConfig() {
	config := Config{
		Host:        "127.0.0.1",
		Port:        "6700",
		AccessToken: "",
		LogLevel:    "INFO",
		ZeroBot: zeroConfig{
			NickName:      []string{"机器人", "笨蛋"},
			CommandPrefix: "/",
			SuperUsers:    []string{},
		},
	}
	if err := config.Save(); err != nil {
		fmt.Println("无法生成配置文件：错误：", err)
		os.Exit(1)
	}
	fmt.Println("配置文件已生成。")
}

func readConfig() (*Config, error) {
	fileContent, err := os.ReadFile("config.toml")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("配置文件不存在，将生成默认配置文件")
			initialConfig()
			fmt.Println("请修改配置文件后再启动。")
		} else {
			fmt.Printf("无法读取配置文件：错误%s\n", err)
		}
		return nil, err
	}
	fileContent = bytes.TrimPrefix(fileContent, []byte{0xef, 0xbb, 0xbf}) // remove utf-8 BOM
	var conf Config
	if _, err = toml.Decode(string(fileContent), &conf); err != nil {
		fmt.Printf("无法解析配置文件：错误%s\n", err)
		return nil, err
	}
	return &conf, nil
}
