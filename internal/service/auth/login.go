package auth

import (
	"context"
	"fmt"

	"github.com/merynayr/auth/internal/utils/hash"
	"github.com/merynayr/auth/internal/utils/jwt"
)

// Login валидирует данные пользователя, и если все ок, возвращает refresh token
func (s *srv) Login(ctx context.Context, name string, password string) (string, error) {
	userInfo, err := s.userRepository.GetUserByName(ctx, name)
	if err != nil {
		return "", err
	}

	err = hash.CompareHashAndPass(password, userInfo.Password)
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	token, err := jwt.GenerateToken(userInfo, s.authCfg.RefreshTokenSecretKey(), s.authCfg.RefreshTokenExp())
	if err != nil {
		return "", err
	}

	return token, nil
}
