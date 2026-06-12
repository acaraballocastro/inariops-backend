package auth

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token              string `json:"token"`
	UserID             string `json:"user_id"`
	Role               string `json:"role"`
	MustChangePassword bool   `json:"must_change_password"`
}

type ChangePasswordRequest struct {
	UserID      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
