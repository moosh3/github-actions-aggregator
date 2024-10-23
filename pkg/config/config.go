package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string
	LogLevel   string
	// Add other configuration fields
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs/")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	return &Config{
		ServerPort: viper.GetString("server.port"),
		LogLevel:   viper.GetString("log.level"),
		// Initialize other fields
	}
}
