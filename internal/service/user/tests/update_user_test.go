package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/merynayr/auth/internal/model"
	"github.com/merynayr/auth/internal/repository"
	repoMocks "github.com/merynayr/auth/internal/repository/mocks"
	"github.com/merynayr/auth/internal/service/user"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	type UserRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		req *model.UserUpdate
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Username()
		email = gofakeit.Email()

		repoErr = fmt.Errorf("repo error")

		req = &model.UserUpdate{
			ID:    id,
			Name:  &name,
			Email: &email,
		}

		res = &emptypb.Empty{}
	)

	tests := []struct {
		name               string
		args               args
		want               *emptypb.Empty
		err                error
		UserRepositoryMock UserRepositoryMockFunc
	}{
		{
			name: "success from repo",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			UserRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, req).Return(res, nil)
				return mock
			},
		},
		{
			name: "error from repo",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  repoErr,
			UserRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, req).Return(nil, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepositoryMock := tt.UserRepositoryMock(mc)

			service := user.NewService(userRepositoryMock, nil)

			res, err := service.UpdateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
