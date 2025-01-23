package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/merynayr/auth/internal/api/user"
	"github.com/merynayr/auth/internal/model"
	"github.com/merynayr/auth/internal/service"
	serviceMocks "github.com/merynayr/auth/internal/service/mocks"

	desc "github.com/merynayr/auth/pkg/user_v1"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mn *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.GetUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Username()
		email     = gofakeit.Email()
		role      = desc.Role(0)
		createdAt = time.Now()

		serviceErr   = fmt.Errorf("service err")
		errConverter = fmt.Errorf("failed to convert desc from user")
		req          = &desc.GetUserRequest{
			Id: id,
		}

		info = &model.User{
			ID:        id,
			Name:      name,
			Email:     email,
			Role:      int32(role),
			CreatedAt: createdAt,
		}

		res = &desc.GetUserResponse{
			User: &desc.User{
				Id:        id,
				Name:      name,
				Email:     email,
				Role:      role,
				CreatedAt: timestamppb.New(createdAt),
			},
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetUserResponse
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
				mock.GetUserMock.Expect(ctx, id).Return(info, nil)
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
				mock.GetUserMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
		{
			name: "error user is nil",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  errConverter,
			UserServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(nil, nil)
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

			res, err := api.GetUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
