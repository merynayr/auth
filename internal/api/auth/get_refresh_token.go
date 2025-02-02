package auth

import (
	"context"
	"fmt"

	desc "github.com/merynayr/auth/pkg/auth_v1"
)

// GetRefreshToken отправляет запрос в сервисный слой на получение refresh токена
func (a *API) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	token, err := a.authService.GetRefreshToken(ctx, req.GetOldRefreshToken())
	if err != nil {
		return nil, fmt.Errorf("codes.Unauthenticated")
	}

	return &desc.GetRefreshTokenResponse{
		RefreshToken: token,
	}, nil
}
