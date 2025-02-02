package repository

import (
	"context"

	"github.com/merynayr/auth/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserRepository - интерфейс репо слоя user
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.UserUpdate) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, userID int64) (*emptypb.Empty, error)
	GetUserByName(ctx context.Context, name string) (*model.UserInfo, error)
	IsNameExist(ctx context.Context, name string) (bool, error)
}
