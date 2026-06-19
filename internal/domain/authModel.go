package domain

type AuthCredentials struct {
	UserID             string `json:"user_id"`
	PasswordHash       string `json:"-"`
	MustChangePassword bool   `json:"must_change_password"`
	IsActive           bool   `json:"is_active"`
}
