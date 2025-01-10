package converter

import (
	"github.com/merynayr/auth/internal/model"
	desc "github.com/merynayr/auth/pkg/user_v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUserFromDescUser конвертирует модель пользователя API слоя в
// модель сервисного слоя
func ToUserFromDescUser(user *desc.CreateUserRequest) *model.User {
	if user == nil {
		return nil
	}

	return &model.User{
		Name:            user.Info.Name,
		Email:           user.Info.Email,
		Password:        user.Info.Password,
		PasswordConfirm: user.Info.PasswordConfirm,
		Role:            int32(user.Info.Role), // по умолчанию роль 1 (USER)
	}
}

// ToDescUserFromService конвертирует сервисную модель пользователя в
// в gRPC модель
func ToDescUserFromService(user *model.User) *desc.GetUserResponse {
	if user == nil {
		return nil
	}

	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.GetUserResponse{
		User: &desc.User{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      desc.Role(user.Role),
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: updatedAt,
		},
	}
}

// ToUserFromDescUpdate конвертирует модель обновления пользователя API слов в
// модель сервисного слоя
func ToUserFromDescUpdate(user *desc.UpdateUserRequest) *model.UserUpdate {
	if user == nil {
		return nil
	}

	u := &model.UserUpdate{}
	u.ID = user.Info.Id.Value

	var name, email string

	if user.Info.Name != nil {
		name = user.Info.Name.GetValue()
		u.Name = &name
	}

	if user.Info.Email != nil {
		email = user.Info.Email.GetValue()
		u.Email = &email
	}

	return u
}
