package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/merynayr/auth/internal/api/user"
	"github.com/merynayr/auth/internal/model"
	"github.com/merynayr/auth/internal/service"
	serviceMocks "github.com/merynayr/auth/internal/service/mocks"

	desc "github.com/merynayr/auth/pkg/user_v1"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mn *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.UpdateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Username()
		email = gofakeit.Email()

		serviceErr = fmt.Errorf("failed to update user")
		converErr  = fmt.Errorf("failed to convert user from desc")

		req = &desc.UpdateUserRequest{
			Info: &desc.UpdateUserInfo{
				Id:    id,
				Name:  name,
				Email: email,
			},
		}

		info = &model.UserUpdate{
			ID:    id,
			Name:  &name,
			Email: &email,
		}

		res = &emptypb.Empty{}
	)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		UserServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			UserServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, info).Return(&emptypb.Empty{}, nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			UserServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, info).Return(nil, serviceErr)
				return mock
			},
		},
		{
			name: "convert error",
			args: args{
				ctx: ctx,
				req: nil,
			},
			want: nil,
			err:  converErr,
			UserServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			UserServiceMock := tt.UserServiceMock(mc)
			api := user.NewAPI(UserServiceMock)

			res, err := api.UpdateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
