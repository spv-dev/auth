package model

type UserCache struct {
	ID        string `redis:"id"`
	Name      string `redis:"name"`
	Email     string `redis:"email"`
	Role      int32  `redis:"role"`
	CreatedAt int64  `redis:"created_at"`
	UpdatedAt *int64 `redis:"updated_at"`
}
