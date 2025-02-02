package converter

import (
	"github.com/merynayr/auth/internal/model"
	modelRepo "github.com/merynayr/auth/internal/repository/user/model"
)

// ToUserFromRepo конвертирует модель пользователя репо слоя в
// модель сервисного слоя
func ToUserFromRepo(user *modelRepo.User) *model.User {
	if user == nil {
		return nil
	}

	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserInfoFromRepo конвертирует модель пользователя репо слоя в
// модель сервисного слоя для авторизации
func ToUserInfoFromRepo(user *modelRepo.UserInfo) *model.UserInfo {
	if user == nil {
		return nil
	}

	return &model.UserInfo{
		Name:     user.Name,
		Password: user.Password,
		Role:     user.Role,
	}
}
