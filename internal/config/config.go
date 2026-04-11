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
	VK              VKConfig
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
	Enabled           bool   `env:"TELEGRAM_ENABLED" env-default:"false"`
	BotToken          string `env:"TELEGRAM_BOT_TOKEN"`
	GraphicsPath      string `env:"TELEGRAM_BOT_GRAPHICS_PATH"`
	GreetingStickerID string `env:"TELEGRAM_BOT_GREETING_STICKER_ID"`
	AuthorName        string `env:"TELEGRAM_BOT_AUTHOR_NAME"`
	Timeout           int    `env:"TELEGRAM_BOT_TIMEOUT" env-default:"60"`
	Debug             bool   `env:"TELEGRAM_BOT_DEBUG" env-default:"false"`
}

type VKConfig struct {
	Enabled      bool   `env:"VK_ENABLED" env-default:"false"`
	AccessToken  string `env:"VK_ACCESS_TOKEN"`
	GroupID      int    `env:"VK_GROUP_ID"`
	GraphicsPath string `env:"VK_GRAPHICS_PATH"`
	AuthorID     int    `env:"VK_AUTHOR_ID"`
	AuthorName   string `env:"VK_AUTHOR_NAME"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("No loading env variables: %v", err)
	}

	if cfg.Telegram.Enabled {
		if cfg.Telegram.BotToken == "" {
			log.Fatal("TELEGRAM_BOT_TOKEN is required when TELEGRAM_ENABLED=true")
		}
		if cfg.Telegram.GraphicsPath == "" {
			log.Fatal("TELEGRAM_BOT_GRAPHICS_PATH is required when TELEGRAM_ENABLED=true")
		}
	}

	if cfg.VK.Enabled {
		if cfg.VK.AccessToken == "" {
			log.Fatal("VK_ACCESS_TOKEN is required when VK_ENABLED=true")
		}
		if cfg.VK.GroupID == 0 {
			log.Fatal("VK_GROUP_ID is required when VK_ENABLED=true")
		}
		if cfg.VK.GraphicsPath == "" {
			log.Fatal("VK_GRAPHICS_PATH is required when VK_ENABLED=true")
		}
	}

	if !cfg.Telegram.Enabled && !cfg.VK.Enabled {
		log.Fatal("At least one bot must be enabled (TELEGRAM_ENABLED or VK_ENABLED)")
	}

	return &cfg
}
