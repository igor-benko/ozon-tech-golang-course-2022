package config

import (
	"errors"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	StoragePostgres = "postgres"
	StorageMemory   = "memory"
)

type Config struct {
	PersonService PersonServiceConfig `yaml:"person_service"`
	Storage       StorageConfig       `yaml:"storage"`
	Telegram      TelegramConfig      `yaml:"telegram"`
	Pooler        PoolerConfig        `yaml:"pooler"`
	Database      DatabaseConfig      `yaml:"database"`
}

type PersonServiceConfig struct {
	Port        int    `yaml:"port" env:"PERSON_SERVICE_PORT"`
	GatewayPort int    `yaml:"gateway_port" env:"PERSON_SERVICE_GATEWAY_PORT"`
	Storage     string `yaml:"storage" env:"PERSON_SERVICE_STORAGE"`
}

type StorageConfig struct {
	PoolSize  int `yaml:"pool_size" env:"STORAGE_POOL_SIZE"`
	TimeoutMs int `yaml:"timeout_ms" env:"STORAGE_TIMEOUT_MS"`
}

type TelegramConfig struct {
	ApiKey          string `yaml:"api_key" env:"TELEGRAM_API_KEY"`
	Timeout         int    `yaml:"timeout" env:"TELEGRAM_TIMEOUT"`
	Offset          int    `yaml:"offset" env:"TELEGRAM_OFFSET"`
	PersonService   string `yaml:"person_service" env:"TELEGRAM_PERSON_SERVICE"`
	RetryMax        uint   `yaml:"retry_max" env:"TELEGRAM_RETRY_MAX"`
	RetryIntervalMs uint   `yaml:"retry_interval_ms" env:"TELEGRAM_RETRY_INTERVAL_MS"`
}

type PoolerConfig struct {
	Host     string `yaml:"host" env:"POOLER_HOST"`
	Port     int    `yaml:"port" env:"POOLER_PORT"`
	User     string `yaml:"user" env:"POOLER_USER"`
	Password string `yaml:"password" env:"POOLER_PASSWORD"`
	Name     string `yaml:"name" env:"POOLER_NAME"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DATABASE_HOST"`
	Port     int    `yaml:"port" env:"DATABASE_PORT"`
	User     string `yaml:"user" env:"DATABASE_USER"`
	Password string `yaml:"password" env:"DATABASE_PASSWORD"`
	Name     string `yaml:"name" env:"DATABASE_NAME"`
}

func Init() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config.yaml", cfg)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
