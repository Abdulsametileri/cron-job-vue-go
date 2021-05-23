package config

import (
	"github.com/spf13/viper"
	"log"
)

func Setup() {
	if IsDebug {
		viper.SetConfigFile("config_dev.json")
	} else {
		viper.SetConfigFile(".env")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
}
