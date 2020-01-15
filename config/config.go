package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Server struct {
	Port string
}

type Database struct {
	Uri          string
	DatabaseName string
	Username     string
	Password     string
}

// Represents database server and credentials
type Config struct {
	Server   Server
	Database Database
}

// Read and parse the Config file
func (c *Config) Read() {
	viper.SetConfigType("yml")
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/go/src/github.com/wallacebenevides/star-wars-api") // path to look for the config file in
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}
