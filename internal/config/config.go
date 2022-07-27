package config

import (
	"errors"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Grpc     GrpcConfig     `yaml:"grpc"`
	Storage  StorageConfig  `yaml:"storage"`
	Telegram TelegramConfig `yaml:"telegram"`
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
