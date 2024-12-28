package user

import (
	"context"

	"github.com/merynayr/auth/internal/converter"
	desc "github.com/merynayr/auth/pkg/user_v1"
)

// GetUser - отправляет запрос в сервисный слой на получение данных пользователя
func (a *API) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	userObj, err := a.userService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converter.ToDescUserFromService(userObj), nil
}
