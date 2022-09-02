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
	PersonService    PersonServiceConfig    `yaml:"person_service"`
	PersonConsumer   PersonConsumerConfig   `yaml:"person_consumer"`
	VerifyConsumer   VerifyConsumerConfig   `yaml:"verify_consumer"`
	RollbackConsumer RollbackConsumerConfig `yaml:"rollback_consumer"`
	Storage          StorageConfig          `yaml:"storage"`
	Telegram         TelegramConfig         `yaml:"telegram"`
	Pooler           PoolerConfig           `yaml:"pooler"`
	Database         DatabaseConfig         `yaml:"database"`
	Kafka            KafkaConfig            `yaml:"kafka"`
}

type PersonServiceConfig struct {
	AppName     string `yaml:"app_name" env:"PERSON_SERVICE_APP_NAME"`
	Port        int    `yaml:"port" env:"PERSON_SERVICE_PORT"`
	GatewayPort int    `yaml:"gateway_port" env:"PERSON_SERVICE_GATEWAY_PORT"`
	Storage     string `yaml:"storage" env:"PERSON_SERVICE_STORAGE"`
}

type PersonConsumerConfig struct {
	AppName    string `yaml:"app_name" env:"PERSON_CONSUMER_APP_NAME"`
	GroupName  string `yaml:"group_name" env:"PERSON_CONSUMER_GROUP_NAME"`
	ExpvarPort int    `yaml:"expvar_port" env:"PERSON_CONSUMER_EXPVAR_PORT"`
}

type VerifyConsumerConfig struct {
	GroupName string `yaml:"group_name" env:"VERIFY_CONSUMER_GROUP_NAME"`
}

type RollbackConsumerConfig struct {
	GroupName string `yaml:"group_name" env:"ROLLBACK_CONSUMER_GROUP_NAME"`
}

type StorageConfig struct {
	PoolSize  int `yaml:"pool_size" env:"STORAGE_POOL_SIZE"`
	TimeoutMs int `yaml:"timeout_ms" env:"STORAGE_TIMEOUT_MS"`
}

type TelegramConfig struct {
	AppName         string `yaml:"app_name" env:"TELEGRAM_APP_NAME"`
	ApiKey          string `yaml:"api_key" env:"TELEGRAM_API_KEY"`
	Timeout         int    `yaml:"timeout" env:"TELEGRAM_TIMEOUT"`
	Offset          int    `yaml:"offset" env:"TELEGRAM_OFFSET"`
	PersonService   string `yaml:"person_service" env:"TELEGRAM_PERSON_SERVICE"`
	RetryMax        uint   `yaml:"retry_max" env:"TELEGRAM_RETRY_MAX"`
	RetryIntervalMs uint   `yaml:"retry_interval_ms" env:"TELEGRAM_RETRY_INTERVAL_MS"`
	ExpvarPort      int    `yaml:"expvar_port" env:"TELEGRAM_EXPVAR_PORT"`
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

type KafkaConfig struct {
	Brokers     []string `yaml:"brokers" env:"KAFKA_BROKERS"`
	IncomeTopic string   `yaml:"income_topic" env:"KAFKA_INCOME_TOPIC"`
	VerifyTopic string   `yaml:"verify_topic" env:"KAFKA_VERIFY_TOPIC"`
	ErrorTopic  string   `yaml:"error_topic" env:"KAFKA_ERROR_TOPIC"`
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
