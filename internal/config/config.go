package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `env:"ENV" env-default:"local"`
	DB       DBConfig
	Telegram TelegramConfig
}

type DBConfig struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     string `env:"DB_PORT" env-required:"true"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	Name     string `env:"DB_NAME" env-required:"true"`
}

type TelegramConfig struct {
	BotToken string `env:"TELEGRAM_BOT_TOKEN" env-required:"true"`
	Timeout  int    `env:"TELEGRAM_BOT_TIMEOUT" env-default:"60"`
	Debug    bool   `env:"TELEGRAM_BOT_DEBUG" env-default:"false"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("No loading env variables: %v", err)
	}

	return &cfg
}
