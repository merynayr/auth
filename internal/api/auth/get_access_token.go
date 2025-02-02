package auth

import (
	"context"
	"fmt"

	desc "github.com/merynayr/auth/pkg/auth_v1"
)

// GetAccessToken отправляет запрос в сервисный слой на получение access токена
func (a *API) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	token, err := a.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, fmt.Errorf("codes.Unauthenticated")
	}

	return &desc.GetAccessTokenResponse{
		AccessToken: token,
	}, nil
}
