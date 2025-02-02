package access

import (
	"context"
	"fmt"

	desc "github.com/merynayr/auth/pkg/access_v1"
)

// Check отправляет запрос в сервисный слой на проверку доступа пользователя к эндпоинту
func (a *API) Check(ctx context.Context, req *desc.CheckRequest) (*desc.CheckResponse, error) {
	username, err := a.accessService.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		return nil, fmt.Errorf("permission denied")
	}

	return &desc.CheckResponse{
		Username: username,
	}, nil
}
