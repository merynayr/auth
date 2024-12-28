package user

import (
	"context"

	"github.com/merynayr/auth/internal/converter"
	desc "github.com/merynayr/auth/pkg/user_v1"
)

// CreateUser - отправляет запрос в сервисный слой на создание пользователя
func (a *API) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	userID, err := a.userService.CreateUser(ctx, converter.ToUserFromDescUser(req))
	if err != nil {
		return nil, err
	}

	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}
