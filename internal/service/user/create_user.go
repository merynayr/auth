package user

import (
	"context"

	"github.com/merynayr/auth/internal/model"
)

func (s *srv) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	var userID int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		userID, errTx = s.userRepository.CreateUser(ctx, user)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.userRepository.GetUser(ctx, userID)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return userID, nil
}
