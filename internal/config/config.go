package config

import (
	"errors"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Telegram TelegramConfig `yaml:"telegram"`
}

type TelegramConfig struct {
	ApiKey string `yaml:"api_key" env:"TELEGRAM_API_KEY"`
}

var cfg *Config

func Init() error {
	cfg = &Config{}

	err := cleanenv.ReadConfig("./config.yaml", cfg)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return err
	}

	return nil
}

func Get() *Config {
	return cfg
}
