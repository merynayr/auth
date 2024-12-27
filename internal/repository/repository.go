package repository

import (
	"context"

	"github.com/merynayr/auth/internal/model"
)

// UserRepository - интерфейс репо слоя user
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUser(ctx context.Context, userID int64) (*model.User, error)
}
