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
	App      AppConfig      `yaml:"app"`
	Grpc     GrpcConfig     `yaml:"grpc"`
	Storage  StorageConfig  `yaml:"storage"`
	Telegram TelegramConfig `yaml:"telegram"`
	Pooler   PoolerConfig   `yaml:"pooler"`
	Database DatabaseConfig `yaml:"database"`
}

type AppConfig struct {
	Storage string `yaml:"storage" env:"APP_STORAGE"`
}

type GrpcConfig struct {
	Port        int `yaml:"port" env:"GRPC_PORT"`
	GatewayPort int `yaml:"gateway_port" env:"GRPC_GATEWAY_PORT"`
}

type StorageConfig struct {
	PoolSize  int `yaml:"pool_size" env:"STORAGE_POOL_SIZE"`
	TimeoutMs int `yaml:"timeout_ms" env:"STORAGE_TIMEOUT_MS"`
}

type TelegramConfig struct {
	ApiKey  string `yaml:"api_key" env:"TELEGRAM_API_KEY"`
	Timeout int    `yaml:"timeout" env:"TELEGRAM_TIMEOUT"`
	Offset  int    `yaml:"offset" env:"TELEGRAM_OFFSET"`
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
