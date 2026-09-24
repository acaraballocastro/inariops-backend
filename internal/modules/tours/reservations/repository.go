package reservations

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

func (r *Repository) CreateReservation(reservation Reservation) (Reservation, error) {

	err := r.db.QueryRow(`
        INSERT INTO reservations (
            id,
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
        )
        VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
            $11, $12, $13, $14
        )
        RETURNING code
    `,
		reservation.ID,
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
		reservation.CreatedAt,
		reservation.UpdatedAt,
	).Scan(
		&reservation.Code,
	)

	if err != nil {
		return Reservation{}, err
	}

	return reservation, nil
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
			quote_number = $3,
			file_number = $4,
			agency_id = $5,
			total_people_count = $6,
			start_date = $7,
			end_date = $8,
			updated_at = NOW()
		WHERE code = $9
	`,
		reservation.Title,
		reservation.Description,
		reservation.QuoteNumber,
		reservation.FileNumber,
		reservation.AgencyID,
		reservation.TotalPeopleCount,
		reservation.StartDate,
		reservation.EndDate,
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
		UPDATE reservations
		SET status = $1,
			updated_at = NOW()
		WHERE code = $2
	`, domain.RESERVATION_CANCELLED, code)

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

func (r *Repository) UpdateReservationStatus(code string, status domain.ReservationStatus) error {
	result, err := r.db.Exec(`
		UPDATE reservations
		SET status = $1,
			updated_at = NOW()
		WHERE code = $2
	`, status, code)

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

func (r *Repository) UpdateSignatureStatus(code string, status domain.SignatureStatus) error {
	result, err := r.db.Exec(`
		UPDATE reservations SET signature_status = $1, updated_at = NOW() WHERE code = $2
	`, status, code)
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

func (r *Repository) UpdateVoucherStatus(code string, status domain.VoucherStatus) error {
	result, err := r.db.Exec(`
		UPDATE reservations SET voucher_status = $1, updated_at = NOW() WHERE code = $2
	`, status, code)
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

func (r *Repository) EraseReservation(code string) error {
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
