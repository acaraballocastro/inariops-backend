package domain

import "time"

type UserRole string

const (
	RoleAdmin UserRole = "ADMIN"
	RoleGuide UserRole = "GUIDE"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      UserRole  `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type UserCredentials struct {
	User            User
	AuthCredentials AuthCredentials
}
