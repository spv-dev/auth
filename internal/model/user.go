package model

import (
	"database/sql"
	"time"
)

// User модель пользователя
type User struct {
	ID        int64
	Info      UserInfo
	Password  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

// UserInfo модель информации о пользователе
type UserInfo struct {
	Name  string
	Email string
	Role  int32
}

// UpdateUserInfo модель для изменения информации о пользователе
type UpdateUserInfo struct {
	Name sql.NullString
	Role sql.NullInt32
}
