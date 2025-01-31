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

	u := &model.User{}

	if user.Info.Name != "" {
		u.Name = user.Info.Name
	} else {
		return nil
	}

	if user.Info.Email != "" {
		u.Email = user.Info.Email
	} else {
		return nil
	}

	if user.Info.Password != "" {
		u.Password = user.Info.Password
	} else {
		return nil
	}

	if user.Info.PasswordConfirm != "" {
		u.PasswordConfirm = user.Info.PasswordConfirm
	} else {
		return nil
	}

	if user.Info.Role == 0 || user.Info.Role == 1 {
		u.Role = int32(user.Info.Role)
	} else {
		return nil
	}
	return u
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
	if user.Info.Id > 0 {
		u.ID = user.Info.Id
	} else {
		return nil
	}

	var name, email string

	if user.Info.Name != "" {
		name = user.Info.Name
		u.Name = &name
	} else {
		return nil
	}

	if user.Info.Email != "" {
		email = user.Info.Email
		u.Email = &email
	} else {
		return nil
	}

	return u
}
