package constants

type Roles int32

const (
	Roles_UNKNOWN Roles = 0
	Roles_ADMIN   Roles = 1
	Roles_USER    Roles = 2
)

var (
	Roles_name = map[int32]string{
		0: "UNKNOWN",
		1: "ADMIN",
		2: "USER",
	}
	Roles_value = map[string]int32{
		"UNKNOWN": 0,
		"ADMIN":   1,
		"USER":    2,
	}
)
