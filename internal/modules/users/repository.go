package users

import (
	"database/sql"
	"inariops/internal/db"
	"inariops/internal/domain"
)

type Repository struct {
	db db.DBTX
}

func NewRepository(database db.DBTX) *Repository {
	return &Repository{
		db: database,
	}
}

func (r *Repository) WithTx(tx *sql.Tx) *Repository {
	return &Repository{
		db: tx,
	}
}

func (r *Repository) GetAllUsers() ([]domain.User, error) {
	rows, err := r.db.Query("SELECT id, name, email, phone, role, is_active, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.IsActive, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *Repository) CreateUser(user domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (id, name, email, phone, role, is_active, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.ID, user.Name, user.Email, user.Phone, user.Role, user.IsActive, user.CreatedAt)
	return err
}

func (r *Repository) GetUserByID(id string) (*domain.User, error) {
	row := r.db.QueryRow("SELECT id, name, email, phone, role, is_active, created_at FROM users WHERE id = $1", id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.IsActive, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserByEmail(email string) (*domain.User, error) {
	var u domain.User

	err := r.db.QueryRow(`
		SELECT id, name, email, phone, role, is_active, created_at 
		FROM users
		WHERE email = $1
	`, email).Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Role, &u.IsActive, &u.CreatedAt)

	return &u, err
}

func (r *Repository) UpdateUser(user domain.User) error {
	_, err := r.db.Exec("UPDATE users SET name = $1, email = $2, phone = $3, role = $4, is_active = $5 WHERE id = $6",
		user.Name, user.Email, user.Phone, user.Role, user.IsActive, user.ID)
	return err
}

func (r *Repository) DeactivateUser(id string) error {
	_, err := r.db.Exec("UPDATE users SET is_active = false WHERE id = $1", id)
	return err
}

func (r *Repository) DeleteUser(id string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
