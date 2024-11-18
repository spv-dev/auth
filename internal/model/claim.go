package model

import "github.com/dgrijalva/jwt-go"

// UserClaims структура для пользовательских параметров в jwt токене
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}
