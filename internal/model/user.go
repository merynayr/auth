package model

import (
	"database/sql"
	"time"
)

// User модель пользователя сервисного слоя
type User struct {
	ID              int64
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            int32
	CreatedAt       time.Time
	UpdatedAt       sql.NullTime
}
