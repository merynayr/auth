package service

import (
	"context"

	"github.com/merynayr/auth/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserService интерфейс сервисного слоя user
type UserService interface {
	GetUser(ctx context.Context, userID int64) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	UpdateUser(ctx context.Context, user *model.UserUpdate) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, userID int64) (*emptypb.Empty, error)
}
