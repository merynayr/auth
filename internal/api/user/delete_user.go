package user

import (
	"context"

	desc "github.com/merynayr/auth/pkg/user_v1"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser - отправляет запрос в сервисный слой на удаление пользователя
func (i *API) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	_, err := i.userService.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, errors.Wrap(err, "failed to delete user")
	}

	return &emptypb.Empty{}, nil
}
