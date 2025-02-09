package config

import (
	"fmt"
	"go.uber.org/zap"
)

type StorageConfig struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	Database        string `yaml:"database"`
	Password        string `yaml:"password"`
	Username        string `yaml:"username"`
	URI             string `yaml:"uri"`
	sessionPoolSize string `yaml:"session_pool_size"`
}

func (s *StorageConfig) SetURI(logger *zap.Logger) {
	s.URI = fmt.Sprintf(
		s.URI,
		s.Username,
		s.Password,
		s.Host,
		s.Port,
		s.Database,
	)

	logger.Info("successfully set uri", zap.Any("uri", s.URI))
}

func (s *StorageConfig) GetURI() string {
	return s.URI
}
