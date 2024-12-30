package user

import (
	"github.com/merynayr/auth/internal/client/db"
	"github.com/merynayr/auth/internal/repository"
	"github.com/merynayr/auth/internal/service"
)

// Структура сервисного слоя с объектами репо слоя
// и транзакционного менеджера
type srv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

// NewService возвращает объект сервисного слоя
func NewService(
	userRepository repository.UserRepository,
	txManager db.TxManager,
) service.UserService {
	return &srv{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
