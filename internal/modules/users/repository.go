package users

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllUsers() ([]User, error) {
	rows, err := r.db.Query("SELECT id, name, email, phone, role, is_active, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.IsActive, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *Repository) CreateUser(user User) error {
	_, err := r.db.Exec("INSERT INTO users (id, name, email, phone, role, is_active, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.ID, user.Name, user.Email, user.Phone, user.Role, user.IsActive, user.CreatedAt)
	return err
}

func (r *Repository) GetUserByID(id string) (*User, error) {
	row := r.db.QueryRow("SELECT id, name, email, phone, role, is_active, created_at FROM users WHERE id = $1", id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.IsActive, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) UpdateUser(user User) error {
	_, err := r.db.Exec("UPDATE users SET name = $1, email = $2, phone = $3, role = $4, is_active = $5 WHERE id = $6",
		user.Name, user.Email, user.Phone, user.Role, user.IsActive, user.ID)
	return err
}

func (r *Repository) DeleteUser(id string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
