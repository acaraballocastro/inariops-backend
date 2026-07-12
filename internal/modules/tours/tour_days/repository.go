package tourdays

import (
	"database/sql"
	"inariops/internal/domain"
	tours "inariops/internal/modules/tours/shared"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateTourDay(tourDay tours.TourDay) error {
	_, err := r.db.Exec(`
		INSERT INTO tour_days (
			id, reservation_id, code, title, start_datetime, duration, zone_id, guide_id,
			status, people_count, remuneration, meeting_point, guide_hotel_id, customer_hotel_id,
			guide_liability_notes, customer_liability_notes, voucher_status, created_at, updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19
		)
	`,
		tourDay.ID,
		tourDay.ReservationID,
		tourDay.Code,
		tourDay.Title,
		tourDay.StartDateTime,
		tourDay.Duration,
		tourDay.ZoneID,
		tourDay.GuideID,
		tourDay.Status,
		tourDay.PeopleCount,
		tourDay.Remuneration,
		tourDay.MeetingPoint,
		tourDay.GuideHotelID,
		tourDay.CustomerHotelID,
		tourDay.GuideLiabilityNotes,
		tourDay.CustomerLiabilityNotes,
		tourDay.VoucherStatus,
		tourDay.CreatedAt,
		tourDay.UpdatedAt,
	)

	return err
}

func (r *Repository) GetTourDayByID(id string) (tours.TourDay, error) {
	var tourDay tours.TourDay

	err := r.db.QueryRow(`
		SELECT id, reservation_id, code, title, start_datetime, duration, zone_id, guide_id,
			status, people_count, remuneration, meeting_point, guide_hotel_id, customer_hotel_id,
			guide_liability_notes, customer_liability_notes, voucher_status, created_at, updated_at
		FROM tour_days WHERE id = $1
	`, id).Scan(
		&tourDay.ID,
		&tourDay.ReservationID,
		&tourDay.Code,
		&tourDay.Title,
		&tourDay.StartDateTime,
		&tourDay.Duration,
		&tourDay.ZoneID,
		&tourDay.GuideID,
		&tourDay.Status,
		&tourDay.PeopleCount,
		&tourDay.Remuneration,
		&tourDay.MeetingPoint,
		&tourDay.GuideHotelID,
		&tourDay.CustomerHotelID,
		&tourDay.GuideLiabilityNotes,
		&tourDay.CustomerLiabilityNotes,
		&tourDay.VoucherStatus,
		&tourDay.CreatedAt,
		&tourDay.UpdatedAt,
	)

	return tourDay, err
}

func (r *Repository) GetTourDaysByReservationID(reservationID string) ([]tours.TourDay, error) {
	rows, err := r.db.Query(`
		SELECT id, reservation_id, code, title, start_datetime, duration, zone_id, guide_id,
			status, people_count, remuneration, meeting_point, guide_hotel_id, customer_hotel_id,
			guide_liability_notes, customer_liability_notes, voucher_status, created_at, updated_at
		FROM tour_days WHERE reservation_id = $1
	`, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tourDays []tours.TourDay
	for rows.Next() {
		var tourDay tours.TourDay
		err := rows.Scan(
			&tourDay.ID,
			&tourDay.ReservationID,
			&tourDay.Code,
			&tourDay.Title,
			&tourDay.StartDateTime,
			&tourDay.Duration,
			&tourDay.ZoneID,
			&tourDay.GuideID,
			&tourDay.Status,
			&tourDay.PeopleCount,
			&tourDay.Remuneration,
			&tourDay.MeetingPoint,
			&tourDay.GuideHotelID,
			&tourDay.CustomerHotelID,
			&tourDay.GuideLiabilityNotes,
			&tourDay.CustomerLiabilityNotes,
			&tourDay.VoucherStatus,
			&tourDay.CreatedAt,
			&tourDay.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tourDays = append(tourDays, tourDay)
	}

	return tourDays, nil
}

func (r *Repository) UpdateTourDay(tourDay tours.TourDay) error {
	_, err := r.db.Exec(`
		UPDATE tour_days SET
			title = $1,
			start_datetime = $2,
			duration = $3,
			zone_id = $4,
			guide_id = $5,
			status = $6,
			people_count = $7,
			remuneration = $8,
			meeting_point = $9,
			guide_hotel_id = $10,
			customer_hotel_id = $11,
			guide_liability_notes = $12,
			customer_liability_notes = $13,
			voucher_status = $14,
			updated_at = $15
		WHERE id = $16
	`,
		tourDay.Title,
		tourDay.StartDateTime,
		tourDay.Duration,
		tourDay.ZoneID,
		tourDay.GuideID,
		tourDay.Status,
		tourDay.PeopleCount,
		tourDay.Remuneration,
		tourDay.MeetingPoint,
		tourDay.GuideHotelID,
		tourDay.CustomerHotelID,
		tourDay.GuideLiabilityNotes,
		tourDay.CustomerLiabilityNotes,
		tourDay.VoucherStatus,
		tourDay.UpdatedAt,
		tourDay.ID,
	)
	return err
}

func (r *Repository) DeleteTourDay(id string) error {
	_, err := r.db.Exec("DELETE FROM tour_days WHERE id = $1", id)
	return err
}

func (r *Repository) CancelTourDay(id string) error {
	_, err := r.db.Exec(`
		UPDATE tour_days SET
			status = $1,
			updated_at = $2
		WHERE id = $3
	`,
		domain.RESERVATION_CANCELLED,
		time.Now(),
		id,
	)
	return err
}
