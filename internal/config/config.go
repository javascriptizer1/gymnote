package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env             string        `env:"ENV" env-default:"local"`
	GracefulTimeout time.Duration `env:"GRACEFUL_TIMEOUT" env-default:"10s"`
	DB              DBConfig
	Redis           CacheConfig
	Telegram        TelegramConfig
}

type DBConfig struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     string `env:"DB_PORT" env-required:"false"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	Name     string `env:"DB_NAME" env-required:"true"`
}

func (c *DBConfig) ConnectionString() string {
	var uri string
	appName := "Gymnote"

	if c.Port == "" {
		uri = fmt.Sprintf("mongodb+srv://%s:%s@%s/?appName=%s", c.User, c.Password, c.Host, appName)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/?appName=%s", c.User, c.Password, c.Host, c.Port, appName)
	}

	return uri
}

type CacheConfig struct {
	Address  string `env:"REDIS_ADDRESS" env-required:"true"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
	DB       int    `env:"REDIS_DB" env-required:"true"`
}

type TelegramConfig struct {
	BotToken          string `env:"TELEGRAM_BOT_TOKEN" env-required:"true"`
	GraphicsPath      string `env:"TELEGRAM_BOT_GRAPHICS_PATH" env-required:"true"`
	GreetingStickerID string `env:"TELEGRAM_BOT_GREETING_STICKER_ID" env-required:"false"`
	AuthorName        string `env:"TELEGRAM_BOT_AUTHOR_NAME" env-required:"false"`
	Timeout           int    `env:"TELEGRAM_BOT_TIMEOUT" env-default:"60"`
	Debug             bool   `env:"TELEGRAM_BOT_DEBUG" env-default:"false"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("No loading env variables: %v", err)
	}

	return &cfg
}
