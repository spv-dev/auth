package model

import "github.com/dgrijalva/jwt-go"

const (
	// ExamplePath Путь для примера доступа
	ExamplePath = "/user_v1.UserV1/Get"
)

// UserClaims структура для пользовательских параметров в jwt токене
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}
