package model

import (
	"time"
)

// User модель
type User struct {
	ID        int64      `db:"id"`
	Info      UserInfo   `db:""`
	Password  string     `db:"password"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at,omitempty"`
}

// UserInfo модель
type UserInfo struct {
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  int32  `db:"role"`
}

// UpdateUserInfo модель
type UpdateUserInfo struct {
	Name *string `db:"name,omitempty"`
	Role *int32  `db:"role,omitempty"`
}
