package user

import (
	"context"
	"fmt"

	"github.com/merynayr/auth/internal/model"
)

func (s *srv) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	if user.Password != user.PasswordConfirm {
		return 0, fmt.Errorf("password and passwordConfirm do not match")
	}
	var userID int64
	userID, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
