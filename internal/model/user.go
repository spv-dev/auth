package model

import (
	"time"

	"github.com/spv-dev/auth/internal/constants"
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
	Name  string          `db:"name"`
	Email string          `db:"email"`
	Role  constants.Roles `db:"role"`
}

// UpdateUserInfo модель
type UpdateUserInfo struct {
	Name *string          `db:"name,omitempty"`
	Role *constants.Roles `db:"role,omitempty"`
}

type TokenUserInfo struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
