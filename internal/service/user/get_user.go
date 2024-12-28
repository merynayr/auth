package user

import (
	"context"

	"github.com/merynayr/auth/internal/model"
)

func (s *srv) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.userRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
