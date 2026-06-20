package domain

type AuthCredentials struct {
	UserID             string
	PasswordHash       string
	MustChangePassword bool
	IsActive           bool
}
