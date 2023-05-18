package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	KeepassFile     string `yaml:"keepassFile"`
	KeepassPassword string `yaml:"keepassPassword"`
	ServerPort      string `yaml:"serverPort"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Println("read configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			log.Println(err)
		}
	})
	return instance
}
