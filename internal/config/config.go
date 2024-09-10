package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type DbConfig struct {
	Driver  string `yaml:"driver"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type BotConfig struct {
	Token     string `yaml:"token"`
	ChannelId int64  `yaml:"channel_id"`
}

type OpenAiConfig struct {
	Prompt string `yaml:"prompt"`
	Key    string `yaml:"key"`
}

type ServerConfig struct {
	FetchInterval        time.Duration `yaml:"fetch_interval"`
	NotificationInterval time.Duration `yaml:"notification_interval"`
	FilterKeywords       []string      `yaml:"filter_keywords"`
}

type Config struct {
	DbConfig     DbConfig     `yaml:"db"`
	BotConfig    BotConfig    `yaml:"bot"`
	OpenAiConfig OpenAiConfig `yaml:"openai"`
	ServerConfig ServerConfig `yaml:"server"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		file, err := os.Open("config.yml")
		if err != nil {
			panic(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Println(err)
			}
		}(file)

		data, err := io.ReadAll(file)
		if err != nil {
			panic("failed to read config file: " + err.Error())
		}
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			panic("failed to parse config file: " + err.Error())
		}
	})

	return cfg
}
