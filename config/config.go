package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	ConnectionUri string
	Database      string
}

// Represents database server and credentials
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// Read and parse the Config file
func (c *Config) Read() {
	viper.SetConfigType("yml")
	viper.SetConfigName("config")                                            // name of config file (without extension)
	viper.AddConfigPath("/go/src/github.com/wallacebenevides/star-wars-api") // path to look for the config file in
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}
