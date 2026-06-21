package reservations

import (
	"database/sql"
	"inariops/internal/shared/errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateReservation(reservation Reservation) error {
	_, err := r.db.Exec(`
			INSERT INTO reservations (
				id, title, description, status, quote_number, file_number, agency_id, total_people_count, start_date,
				end_date, signature_status, voucher_status, created_at, updated_at
			)
			VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
				$11, $12, $13, $14
			)
		`,
		reservation.ID, reservation.Title, reservation.Description, reservation.Status,
		reservation.QuoteNumber, reservation.FileNumber, reservation.AgencyID, reservation.TotalPeopleCount, reservation.StartDate,
		reservation.EndDate, reservation.SignatureStatus, reservation.VoucherStatus, reservation.CreatedAt, reservation.UpdatedAt,
	)

	return err
}

func (r *Repository) GetReservationByCode(code string) (Reservation, error) {
	var res Reservation

	err := r.db.QueryRow(`
			SELECT id,
				code,
				title,
				description,
				status,
				quote_number,
				file_number,
				agency_id,
				total_people_count,
				start_date,
				end_date,
				signature_status,
				voucher_status,
				created_at,
				updated_at
			FROM reservations WHERE code = $1
		`, code).Scan(&res.ID, &res.Code, &res.Title, &res.Description, &res.Status, &res.QuoteNumber, &res.FileNumber, &res.AgencyID, &res.TotalPeopleCount, &res.StartDate, &res.EndDate, &res.SignatureStatus, &res.VoucherStatus, &res.CreatedAt, &res.UpdatedAt)

	return res, err
}

func (r *Repository) GetAllReservations() ([]Reservation, error) {
	rows, err := r.db.Query(`
			SELECT
				id,
				code,
				title,
				description,
				status,
				quote_number,
				file_number,
				agency_id,
				total_people_count,
				start_date,
				end_date,
				signature_status,
				voucher_status,
				created_at,
				updated_at
			FROM reservations
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []Reservation
	for rows.Next() {
		var res Reservation
		err := rows.Scan(&res.ID, &res.Code, &res.Title, &res.Description, &res.Status, &res.QuoteNumber, &res.FileNumber, &res.AgencyID, &res.TotalPeopleCount, &res.StartDate, &res.EndDate, &res.SignatureStatus, &res.VoucherStatus, &res.CreatedAt, &res.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}

	return reservations, nil
}

func (r *Repository) UpdateReservation(reservation Reservation) error {
	result, err := r.db.Exec(`
		UPDATE reservations
		SET title = $1,
			description = $2,
			status = $3,
			quote_number = $4,
			file_number = $5,
			agency_id = $6,
			total_people_count = $7,
			start_date = $8,
			end_date = $9,
			signature_status = $10,
			voucher_status = $11,
			updated_at = NOW()
		WHERE code = $12
	`,
		reservation.Title,
		reservation.Description,
		reservation.Status,
		reservation.QuoteNumber,
		reservation.FileNumber,
		reservation.AgencyID,
		reservation.TotalPeopleCount,
		reservation.StartDate,
		reservation.EndDate,
		reservation.SignatureStatus,
		reservation.VoucherStatus,
		reservation.Code,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.ErrReservationNotFound
	}

	return nil
}

func (r *Repository) DeleteReservation(code string) error {
	result, err := r.db.Exec(`
		DELETE FROM reservations
		WHERE code = $1
	`, code)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.ErrReservationNotFound
	}

	return nil
}
