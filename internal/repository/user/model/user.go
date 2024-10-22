package model

import (
	"database/sql"
	"time"
)

// User модель
type User struct {
	ID        int64        `db:"id"`
	Info      UserInfo     `db:""`
	Password  string       `db:"password"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

// UserInfo модель
type UserInfo struct {
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  int32  `db:"role"`
}

// UpdateUserInfo модель
type UpdateUserInfo struct {
	Name sql.NullString `db:"name"`
	Role sql.NullInt32  `db:"role"`
}
