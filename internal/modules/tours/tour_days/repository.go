package tourdays

import (
	"context"
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

func (r *Repository) AssignGuide(tourDayID string, guideID string) error {
	_, err := r.db.Exec(`
		UPDATE tour_days SET
			guide_id = $1,
			status = $2,
			updated_at = $3
		WHERE id = $4
	`,
		guideID,
		domain.RESERVATION_GUIDE_PREASSIGNED,
		time.Now(),
		tourDayID,
	)
	return err
}

func (r *Repository) AssignGuideToMultipleTourDays(tourDayIDs []string, guideID string) error {
	for _, tourDayID := range tourDayIDs {
		_, err := r.db.Exec(`
			UPDATE tour_days SET
				guide_id = $1,
				status = $2,
				updated_at = $3
			WHERE id = $4
		`,
			guideID,
			domain.RESERVATION_GUIDE_PREASSIGNED,
			time.Now(),
			tourDayID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) UnassignGuide(tourDayID string) error {
	_, err := r.db.Exec(`
		UPDATE tour_days SET
			guide_id = NULL,
			status = $1,
			updated_at = $2
		WHERE id = $3
	`,
		domain.RESERVATION_PENDING_ASSIGNMENT,
		time.Now(),
		tourDayID,
	)
	return err
}

// Worker jobs
func (r *Repository) GetExpiredGuideAssignments(ctx context.Context) ([]tours.TourDay, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, reservation_id, status, updated_at
		FROM tour_days
		WHERE status = $1 AND updated_at < $2
	`, domain.RESERVATION_GUIDE_PREASSIGNED, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expiredTourDays []tours.TourDay
	for rows.Next() {
		var tourDay tours.TourDay
		err := rows.Scan(&tourDay.ID, &tourDay.ReservationID, &tourDay.Status, &tourDay.UpdatedAt)
		if err != nil {
			return nil, err
		}
		expiredTourDays = append(expiredTourDays, tourDay)
	}

	return expiredTourDays, nil
}

func (r *Repository) UpdateTourDayStatus(tourDayID string, status domain.ReservationStatus) error {
	_, err := r.db.Exec(`
		UPDATE tour_days SET
			status = $1,
			updated_at = $2
		WHERE id = $3
	`,
		status,
		time.Now(),
		tourDayID,
	)
	return err
}

func (r *Repository) DeleteTourDaysByReservationID(reservationID string) error {
	_, err := r.db.Exec(`
		UPDATE tour_days SET
			status = $1,
			updated_at = $2
		WHERE reservation_id = $3
	`, domain.RESERVATION_CANCELLED, time.Now(), reservationID)
	return err
}

func (r *Repository) DeleteTourDaysByID(id string) error {
	_, err := r.db.Exec(`
		DELETE FROM tour_days WHERE id = $1
	`, id)
	return err
}
