package model

import (
	"database/sql"
	"time"
)

// User модель
type User struct {
	ID        int64        `db:"id"`
	Info      Info         `db:""`
	Password  string       `db:"password"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

// Info модель
type Info struct {
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  int32  `db:"role"`
}
