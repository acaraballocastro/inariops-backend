package auth

import (
	"database/sql"
	"inariops/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByEmail(email string) (domain.User, error) {
	var u domain.User

	err := r.db.QueryRow(`
		SELECT id, email, role, is_active
		FROM users
		WHERE email = $1
	`, email).Scan(&u.ID, &u.Email, &u.Role, &u.IsActive)

	return u, err
}

func (r *Repository) GetCredentials(userID string) (domain.AuthCredentials, error) {
	var c domain.AuthCredentials

	err := r.db.QueryRow(`
		SELECT user_id, password_hash, must_change_password, is_active
		FROM auth_credentials
		WHERE user_id = $1
	`, userID).Scan(&c.UserID, &c.PasswordHash, &c.MustChangePassword, &c.IsActive)

	return c, err
}

func (r *Repository) UpdatePassword(userID, hash string) error {
	_, err := r.db.Exec(`
		UPDATE auth_credentials
		SET password_hash = $1,
		    must_change_password = false,
		    updated_at = NOW()
		WHERE user_id = $2
	`, hash, userID)

	return err
}
