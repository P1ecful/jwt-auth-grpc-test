package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
	"os"
	"time"
)

type Config struct {
	Service ServiceConfig `yaml:"service" env-required:"true"`
	Storage StorageConfig `yaml:"storage" env-required:"true"`
	GRPC    GRPCConfig    `yaml:"grpc" env-required:"true"`
}

type ServiceConfig struct {
	SecretKey      string        `yaml:"secret_key"`
	AccessTokenTTL time.Duration `yaml:"access_token_ttl"`
}

type GRPCConfig struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func LoadConfig(path string, logger *zap.Logger) *Config {
	var cfg = &Config{
		Service: ServiceConfig{
			SecretKey:      "",
			AccessTokenTTL: time.Minute,
		},
		Storage: StorageConfig{
			Host:     "0.0.0.0",
			Port:     "0000",
			Username: "postgres",
			Password: "postgres",
			Database: "postgres",
		},
		GRPC: GRPCConfig{
			Port: "0000",
		},
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Debug("config file does not exist",
			zap.Any("path", path),
		)
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		logger.Debug("failed to read config", zap.Error(err))
	}

	logger.Info("successfully loaded config",
		zap.Any("config", cfg))
	return cfg
}
