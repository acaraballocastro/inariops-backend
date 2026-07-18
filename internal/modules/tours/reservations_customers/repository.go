package reservationscustomers

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddCustomerToReservation(reservationID string, customerIDs []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	for _, customerID := range customerIDs {
		_, err = tx.Exec(`
			INSERT INTO reservations_customers (reservation_id, customer_id)
			VALUES ($1, $2)
			ON CONFLICT (reservation_id, customer_id) DO NOTHING
		`, reservationID, customerID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) RemoveCustomerFromReservation(reservationID string, customerID string) error {
	_, err := r.db.Exec(`
		DELETE FROM reservations_customers
		WHERE reservation_id = $1 AND customer_id = $2
	`, reservationID, customerID)
	return err
}

func (r *Repository) GetCustomersByReservationID(reservationID string) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT customer_id
		FROM reservations_customers
		WHERE reservation_id = $1
	`, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customerIDs []string
	for rows.Next() {
		var customerID string
		if err := rows.Scan(&customerID); err != nil {
			return nil, err
		}
		customerIDs = append(customerIDs, customerID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customerIDs, nil
}
