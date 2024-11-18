package utils

import "golang.org/x/crypto/bcrypt"

// VerifyPassword метод для проверки пароля
func VerifyPassword(hashedPassword string, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}
