package user

import (
	"context"

	"github.com/merynayr/auth/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *srv) UpdateUser(ctx context.Context, user *model.UserUpdate) (*emptypb.Empty, error) {
	_, err := s.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
