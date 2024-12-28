package service

import (
	"context"

	"github.com/merynayr/auth/internal/model"
)

// UserService интерфейс сервисного слоя user
type UserService interface {
	GetUser(ctx context.Context, userID int64) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (int64, error)
}
