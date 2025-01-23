package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/merynayr/auth/internal/client/db"
	"github.com/merynayr/auth/internal/model"
	"github.com/merynayr/auth/internal/repository"
	repositoryMocks "github.com/merynayr/auth/internal/repository/mocks"

	"github.com/merynayr/auth/internal/service/user"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.User
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		userID          = gofakeit.Int64()
		Name            = gofakeit.Name()
		Email           = gofakeit.Email()
		Password        = gofakeit.Password(true, false, true, true, false, 10)
		PasswordConfirm = Password
		Role            = 0

		repoErr = fmt.Errorf("repository err")

		req = &model.User{
			Name:            Name,
			Email:           Email,
			Password:        Password,
			PasswordConfirm: PasswordConfirm,
			Role:            int32(Role),
		}
	)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success test",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: userID,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, req).Return(userID, nil)
				return mock
			},
		},
		{
			name: "error in CreateUser function",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, req).Return(0, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userRepositoryMock := tt.userRepositoryMock(mc)

			service := user.NewService(userRepositoryMock, nil)
			res, err := service.CreateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
