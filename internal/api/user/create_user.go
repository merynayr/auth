package user

import (
	"context"
	"fmt"

	"github.com/merynayr/auth/internal/converter"
	desc "github.com/merynayr/auth/pkg/user_v1"
)

// CreateUser - отправляет запрос в сервисный слой на создание пользователя
func (a *API) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	convertedUser := converter.ToUserFromDescUser(req)
	if convertedUser == nil {
		return nil, fmt.Errorf("failed to create user: Request id bad")
	}

	userID, err := a.userService.CreateUser(ctx, convertedUser)
	if err != nil {
		return nil, err
	}

	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}
