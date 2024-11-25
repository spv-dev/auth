package constants

// Roles тип ролей
type Roles int32

const (
	// RolesUNKNOWN неизвестный
	RolesUNKNOWN Roles = 0
	// RolesADMIN администратор
	RolesADMIN Roles = 1
	// RolesUSER пользователь
	RolesUSER Roles = 2

	// ExampleRole Роль для примера
	ExampleRole = "ADMIN"
)
