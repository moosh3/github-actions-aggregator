package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort      string
	LogLevel        string
	EnablePolling   bool
	PollingInterval time.Duration
	Database        DatabaseConfig
	GitHub          GitHubConfig
}

type GitHubConfig struct {
	ClientID      string
	ClientSecret  string
	AccessToken   string
	WebhookSecret string
}

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs/")
	viper.AutomaticEnv()

	// Bind environment variables to specific keys
	viper.BindEnv("github.client_id", "GITHUB_CLIENT_ID")
	viper.BindEnv("github.client_secret", "GITHUB_CLIENT_SECRET")
	viper.BindEnv("github.access_token", "GITHUB_ACCESS_TOKEN")
	viper.BindEnv("github.webhook_secret", "GITHUB_WEBHOOK_SECRET")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	return &Config{
		ServerPort:      viper.GetString("server.port"),
		LogLevel:        viper.GetString("log.level"),
		EnablePolling:   viper.GetBool("polling.enabled"),
		PollingInterval: viper.GetDuration("polling.interval"),
		GitHub: GitHubConfig{
			ClientID:      viper.GetString("github.client_id"),
			ClientSecret:  viper.GetString("github.client_secret"),
			AccessToken:   viper.GetString("github.access_token"),
			WebhookSecret: viper.GetString("github.webhook_secret"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("database.host"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			DBName:   viper.GetString("database.db_name"),
			Port:     viper.GetString("database.port"),
		},
	}
}
