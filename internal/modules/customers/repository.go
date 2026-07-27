package customers

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateCustomer(customer Customer) error {
	_, err := r.db.Exec(`
		INSERT INTO customers (id, full_name, document_number, phone, email, age, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		customer.ID,
		customer.FullName,
		customer.DocumentNumber,
		customer.Phone,
		customer.Email,
		customer.Age,
		customer.CreatedAt,
	)

	return err
}

func (r *Repository) GetCustomerByID(id string) (Customer, error) {
	var customer Customer

	err := r.db.QueryRow(`
		SELECT id, full_name, document_number, phone, email, age, created_at
		FROM customers
		WHERE id = $1
	`, id).Scan(
		&customer.ID,
		&customer.FullName,
		&customer.DocumentNumber,
		&customer.Phone,
		&customer.Email,
		&customer.Age,
		&customer.CreatedAt,
	)

	return customer, err
}

func (r *Repository) UpdateCustomer(customer Customer) error {
	_, err := r.db.Exec(`
		UPDATE customers
		SET full_name = $1, document_number = $2, phone = $3, email = $4, age = $5
		WHERE id = $6
	`,
		customer.FullName,
		customer.DocumentNumber,
		customer.Phone,
		customer.Email,
		customer.Age,
		customer.ID,
	)

	return err
}

func (r *Repository) DeleteCustomer(id string) error {
	_, err := r.db.Exec(`
		DELETE FROM customers
		WHERE id = $1
	`, id)

	return err
}

func (r *Repository) GetAllCustomers() ([]Customer, error) {
	rows, err := r.db.Query(`
		SELECT id, full_name, document_number, phone, email, age, created_at
		FROM customers
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var customer Customer
		err := rows.Scan(
			&customer.ID,
			&customer.FullName,
			&customer.DocumentNumber,
			&customer.Phone,
			&customer.Email,
			&customer.Age,
			&customer.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

func (r *Repository) SearchCustomers(term string) ([]Customer, error) {
	rows, err := r.db.Query(`
		SELECT
			id,
			full_name,
			document_number,
			phone,
			email,
			age,
			created_at
		FROM customers
		WHERE
			full_name ILIKE '%' || $1 || '%'
			OR document_number ILIKE '%' || $1 || '%'
			OR phone ILIKE '%' || $1 || '%'
			OR email ILIKE '%' || $1 || '%'
		ORDER BY full_name
		LIMIT 20
	`, term)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customer

	for rows.Next() {
		var customer Customer

		err := rows.Scan(
			&customer.ID,
			&customer.FullName,
			&customer.DocumentNumber,
			&customer.Phone,
			&customer.Email,
			&customer.Age,
			&customer.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *Repository) GetCustomerByEmail(email string) (Customer, error) {
	var customer Customer

	err := r.db.QueryRow(`
		SELECT id, full_name, document_number, phone, email, age, created_at
		FROM customers
		WHERE email = $1
	`, email).Scan(
		&customer.ID,
		&customer.FullName,
		&customer.DocumentNumber,
		&customer.Phone,
		&customer.Email,
		&customer.Age,
		&customer.CreatedAt,
	)

	return customer, err
}
