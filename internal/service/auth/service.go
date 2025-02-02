package auth

import (
	"github.com/merynayr/auth/internal/config"
	"github.com/merynayr/auth/internal/repository"
	"github.com/merynayr/auth/internal/service"
)

type srv struct {
	userRepository repository.UserRepository
	authCfg        config.AuthConfig
}

// NewService возвращает новый объект сервисного слоя auth
func NewService(userRepo repository.UserRepository, authCfg config.AuthConfig) service.AuthService {
	return &srv{
		userRepository: userRepo,
		authCfg:        authCfg,
	}
}
