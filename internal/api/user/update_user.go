package user

import (
	"context"

	"github.com/merynayr/auth/internal/converter"
	desc "github.com/merynayr/auth/pkg/user_v1"
	"github.com/pkg/errors"

	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateUser - отправляет запрос в сервисный слой на обновление данных пользователя
func (i *API) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	convertedReq := converter.ToUserFromDescUpdate(req)
	if convertedReq == nil {
		return nil, errors.New("failed to convert user from desc")
	}

	_, err := i.userService.UpdateUser(ctx, convertedReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update user")
	}

	return &emptypb.Empty{}, nil
}
