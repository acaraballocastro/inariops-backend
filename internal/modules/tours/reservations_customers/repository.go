package reservationscustomers

import (
	"database/sql"
	"inariops/internal/db"
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

func (r *Repository) AddCustomerToReservation(reservationID string, customerIDs []string) error {
	for _, customerID := range customerIDs {
		_, err := r.db.Exec(`
			INSERT INTO reservation_customers (reservation_id, customer_id)
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
		DELETE FROM reservation_customers
		WHERE reservation_id = $1 AND customer_id = $2
	`, reservationID, customerID)
	return err
}

func (r *Repository) GetCustomersByReservationID(reservationID string) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT customer_id
		FROM reservation_customers
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

func (r *Repository) IsCustomerInReservation(reservationID string, customerID string) bool {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM reservation_customers
			WHERE reservation_id = $1 AND customer_id = $2
		)
	`, reservationID, customerID).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (r *Repository) RemoveAllReservationsFromCustomer(customerID string) error {
	_, err := r.db.Exec(`
		DELETE FROM reservation_customers
		WHERE customer_id = $1
	`, customerID)
	return err
}
