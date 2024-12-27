package config

import (
	"github.com/joho/godotenv"
)

// Load читает .env файл по указанному пути
// и загружает переменные в проект
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

// GRPCConfig is interface of a grpc config
type GRPCConfig interface {
	Address() string
}

// PGConfig is interface of a postgres config
type PGConfig interface {
	DSN() string
}
