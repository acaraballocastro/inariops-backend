package auth

import (
	"database/sql"
	"inariops/internal/domain"
	"inariops/internal/shared/errors"
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
		FROM user_credentials
		WHERE user_id = $1
	`, userID).Scan(&c.UserID, &c.PasswordHash, &c.MustChangePassword, &c.IsActive)

	return c, err
}

func (r *Repository) CreateCredentials(credentials Credentials) error {
	_, err := r.db.Exec(`
		INSERT INTO user_credentials (id, user_id, password_hash, must_change_password, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, credentials.ID, credentials.UserID, credentials.PasswordHash, credentials.MustChangePassword, credentials.IsActive, credentials.CreatedAt, credentials.UpdatedAt)
	return err
}

func (r *Repository) UpdatePassword(userID, hash string) error {
	result, err := r.db.Exec(`
		UPDATE user_credentials
		SET password_hash = $1,
		    must_change_password = false,
		    updated_at = NOW()
		WHERE user_id = $2
	`, hash, userID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}
