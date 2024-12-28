package user

import (
	"github.com/merynayr/auth/internal/repository"
	"github.com/merynayr/auth/internal/service"
)

type srv struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) service.UserService {
	return &srv{
		userRepository: userRepository,
	}
}
