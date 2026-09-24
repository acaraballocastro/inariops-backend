package availabilities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAvailability(availability *Availability) (*Availability, error) {

	_, err := r.db.Exec(`
		INSERT INTO availabilities (
			id,
			guide_id,
			start_date,
			end_date,
			reason
		)
		VALUES ($1, $2, $3, $4, $5)
	`,
		availability.ID,
		availability.GuideID,
		availability.StartDate,
		availability.EndDate,
		availability.Reason,
	)

	if err != nil {
		return nil, err
	}

	return availability, nil
}

func (r *Repository) GetAvailabilityByID(id string) (*Availability, error) {

	availability := &Availability{}

	err := r.db.QueryRow(`
		SELECT
			id,
			guide_id,
			start_date,
			end_date,
			reason
		FROM availabilities
		WHERE id = $1
	`, id).Scan(
		&availability.ID,
		&availability.GuideID,
		&availability.StartDate,
		&availability.EndDate,
		&availability.Reason,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return availability, nil
}

func (r *Repository) GetAvailabilitiesByGuideID(guideID string) ([]*Availability, error) {

	rows, err := r.db.Query(`
		SELECT
			id,
			guide_id,
			start_date,
			end_date,
			reason
		FROM availabilities
		WHERE guide_id = $1
		ORDER BY start_date ASC
	`, guideID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var availabilities []*Availability

	for rows.Next() {

		availability := &Availability{}

		if err := rows.Scan(
			&availability.ID,
			&availability.GuideID,
			&availability.StartDate,
			&availability.EndDate,
			&availability.Reason,
		); err != nil {
			return nil, err
		}

		availabilities = append(availabilities, availability)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return availabilities, nil
}

func (r *Repository) DeleteAvailability(id string) error {
	result, err := r.db.Exec(`
		DELETE FROM availabilities
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) UpdateAvailability(availability *Availability) error {
	result, err := r.db.Exec(`
		UPDATE availabilities
		SET
			guide_id = $1,
			start_date = $2,
			end_date = $3,
			reason = $4
		WHERE id = $5
	`,
		availability.GuideID,
		availability.StartDate,
		availability.EndDate,
		availability.Reason,
		availability.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) GetAllAvailabilities() ([]*Availability, error) {

	rows, err := r.db.Query(`
		SELECT
			id,
			guide_id,
			start_date,
			end_date,
			reason
		FROM availabilities
		ORDER BY start_date ASC
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var availabilities []*Availability

	for rows.Next() {

		availability := &Availability{}

		if err := rows.Scan(
			&availability.ID,
			&availability.GuideID,
			&availability.StartDate,
			&availability.EndDate,
			&availability.Reason,
		); err != nil {
			return nil, err
		}

		availabilities = append(availabilities, availability)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return availabilities, nil
}

func (r *Repository) HasConflictingAvailability(guideID string, startDate, endDate *time.Time, excludeID *uuid.UUID) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM availabilities
		WHERE guide_id = $1
		AND start_date <= $3
		AND end_date >= $2`
	args := []any{guideID, startDate, endDate}
	if excludeID != nil {
		query += " AND id <> $4"
		args = append(args, excludeID)
	}

	err := r.db.QueryRow(query, args...).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
