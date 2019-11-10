package main

import (
	"github.com/ega-forever/otus-image-service/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {

	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.SetDefault("LOG_LEVEL", 4)
	viper.SetDefault("REST_PORT", "8080")
	viper.SetDefault("LRU_CACHE", 10)
	viper.SetDefault("STORE_DIR", "temp")

	viper.ReadInConfig()
	viper.AutomaticEnv()

	logLevel := viper.GetInt("LOG_LEVEL")
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.Level(logLevel))
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
