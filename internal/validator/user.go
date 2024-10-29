package validator

import (
	"fmt"
	"regexp"
)

const (
	errUserNameIsEmpty   = "Пустое имя пользователя"
	errUserNameIsTooLong = "Имя не должно превышать 255 символов"
	errEmailIsEmpty      = "Пустой email пользователя"
	errEmailIsTooLong    = "Email не должен превышать 255 символов"
	errEmailIsNotValid   = "Указан неверный Email"
	errPIsEmpty          = "Пустой пароль"
	errPIsTooSmall       = "Длина пароля должна быть не меньше 8 символов"
)

// CheckName проверки имени пользователя
func CheckName(name *string) error {
	if name == nil || len(*name) == 0 {
		return fmt.Errorf(errUserNameIsEmpty)
	}
	if len(*name) > 255 {
		return fmt.Errorf(errUserNameIsTooLong)
	}

	return nil
}

func isValidEmail(email string) bool {
	regEx := "^[a-z0-9._%+-]+@[a-z0-9.-]+.[a-z]{2,}$"
	regex := regexp.MustCompile(regEx)
	return regex.MatchString(email)
}

// CheckEmail проверки email пользователя
func CheckEmail(email *string) error {
	if email == nil || len(*email) == 0 {
		return fmt.Errorf(errEmailIsEmpty)
	}
	if len(*email) > 255 {
		return fmt.Errorf(errEmailIsTooLong)
	}
	if !isValidEmail(*email) {
		return fmt.Errorf(errEmailIsNotValid)
	}
	return nil
}

// CheckPassword проверки пароля пользователя
func CheckPassword(password string) error {
	if len(password) == 0 {
		return fmt.Errorf(errPIsEmpty)
	}
	if len(password) < 8 {
		return fmt.Errorf(errPIsTooSmall)
	}
	return nil
}
