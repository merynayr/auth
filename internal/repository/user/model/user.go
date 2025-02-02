package model

import (
	"database/sql"
	"time"
)

// User модель пользователя в репо слое
type User struct {
	ID        int64        `db:"id"`
	Name      string       `db:"username"`
	Email     string       `db:"email"`
	Password  string       `db:"password"`
	Role      int32        `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

// UserInfo модель пользователя для Авторизации
type UserInfo struct {
	Name     string `db:"username"`
	Password string `db:"password"`
	Role     int64  `db:"role"`
}
