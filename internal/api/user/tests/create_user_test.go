package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/merynayr/auth/internal/api/user"
	"github.com/merynayr/auth/internal/model"
	"github.com/merynayr/auth/internal/service"
	serviceMocks "github.com/merynayr/auth/internal/service/mocks"
	desc "github.com/merynayr/auth/pkg/user_v1"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		userID          = gofakeit.Int64()
		Name            = gofakeit.Name()
		Email           = gofakeit.Email()
		Password        = gofakeit.Password(true, false, true, true, false, 10)
		PasswordConfirm = Password
		Role            = desc.Role_USER
		serviceErr      = fmt.Errorf("service err")

		req = &desc.CreateUserRequest{
			Info: &desc.CreateUserInfo{
				Name:            Name,
				Email:           Email,
				Password:        Password,
				PasswordConfirm: PasswordConfirm,
				Role:            Role,
			},
		}

		info = &model.User{
			Name:            Name,
			Email:           Email,
			Password:        Password,
			PasswordConfirm: PasswordConfirm,
			Role:            0,
		}

		res = &desc.CreateUserResponse{
			Id: userID,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateUserResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success test",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateUserMock.Expect(ctx, info).Return(userID, nil)
				return mock
			},
		},
		{
			name: "error from service",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateUserMock.Expect(ctx, info).Return(int64(0), serviceErr)
				return mock
			},
		},
		{
			name: "Request is nil",
			args: args{
				ctx: ctx,
				req: nil,
			},
			want: nil,
			err:  errors.New("failed to create user: Request id bad"),
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				return mock
			},
		},
		{
			name: "empty Name",
			args: args{
				ctx: ctx,
				req: &desc.CreateUserRequest{
					Info: &desc.CreateUserInfo{
						Email:           Email,
						Password:        Password,
						PasswordConfirm: PasswordConfirm,
						Role:            Role,
					},
				},
			},
			want: nil,
			err:  errors.New("failed to create user: Request id bad"),
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceMock := tt.userServiceMock(mc)
			api := user.NewAPI(userServiceMock)

			res, err := api.CreateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
