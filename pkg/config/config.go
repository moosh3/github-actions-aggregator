package config

import (
	"log"

	"github.com/spf13/viper"
)

type GitHubConfig struct {
	ClientID      string
	ClientSecret  string
	AccessToken   string
	WebhookSecret string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type Config struct {
	ServerPort string
	LogLevel   string
	GitHub     GitHubConfig
	Database   DatabaseConfig
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
		GitHub: GitHubConfig{
			ClientID:      viper.GetString("github.client_id"),
			ClientSecret:  viper.GetString("github.client_secret"),
			AccessToken:   viper.GetString("github.access_token"),
			WebhookSecret: viper.GetString("github.webhook_secret"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetString("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
		},
	}
}
