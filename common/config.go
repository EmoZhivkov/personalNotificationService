package common

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HostPort        string `env:"API_PORT"`
	JwtSecret       string `env:"JWT_SECRET"`
	DbConnectionURL string `env:"DB_CONNECTION_URL"`
	KafkaBrokerURL  string `env:"KAFKA_BROKER_URL"`
	SlackBotToken   string `env:"SLACK_BOT_TOKEN"`
	SlackChannelID  string `env:"SLACK_CHANNEL_ID"`
	SmtpHost        string `env:"SMTP_HOST"`
	SmtpPort        string `env:"SMTP_PORT"`
	SmtpUsername    string `env:"SMTP_USERNAME"`
	SmtpPassword    string `env:"SMTP_PASSWORD"`
	SmscHost        string `env:"SMSC_HOST"`
	SmscPort        string `env:"SMSC_PORT"`
	SmscUsername    string `env:"SMSC_USERNAME"`
	SmscPassword    string `env:"SMSC_PASSWORD"`
}

func NewConfig() Config {
	var config Config
	if err := env.Parse(&config); err != nil {
		log.Fatalf("error parsing the environment variables: %v", err)
	}

	if config.HostPort == "" {
		log.Fatal("API_PORT is not set")
	}
	if config.JwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	if config.DbConnectionURL == "" {
		log.Fatal("DB_CONNECTION_URL is not set")
	}
	if config.KafkaBrokerURL == "" {
		log.Fatal("KAFKA_BROKER_URL is not set")
	}
	if config.SlackBotToken == "" {
		log.Fatal("SLACK_BOT_TOKEN is not set")
	}
	if config.SmtpHost == "" {
		log.Fatal("SMTP_HOST is not set")
	}
	if config.SmtpPort == "" {
		log.Fatal("SMTP_PORT is not set")
	}
	if config.SmscHost == "" {
		log.Fatal("SMSC_HOST is not set")
	}
	if config.SmscPort == "" {
		log.Fatal("SMSC_PORT is not set")
	}
	if config.SmscUsername == "" {
		log.Fatal("SMSC_USERNAME is not set")
	}
	if config.SmscPassword == "" {
		log.Fatal("SMSC_PASSWORD is not set")
	}
	// the smtp username and pass can be left blank for local development
	//if config.SmtpUsername == "" {
	//	log.Fatal("SMTP_USERNAME is not set")
	//}
	//if config.SmtpPassword == "" {
	//	log.Fatal("SMTP_PASSWORD is not set")
	//}
	return config
}
