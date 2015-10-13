package utils

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	ConsumerKey       string `toml:"consumer_key"`
	ConsumerSecret    string `toml:"consumer_secret"`
	AccessToken       string `toml:"access_token"`
	AccessTokenSecret string `toml:"access_token_secret"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load() {
	_, err := toml.DecodeFile("config.tml", c)
	if err != nil {
		log.Fatal(err)
	}
}
