package user

import (
	"context"

	"github.com/merynayr/auth/internal/model"
)

func (s *srv) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	id, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}
