package tourdaysapp

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

func (r *Repository) AddAssignmentStatus(assigment TourDayAssigmentHistory) error {
	_, err := r.db.Exec(`
		INSERT INTO tour_day_assignments_history (
			id, tour_day_id, guide_id, action, changed_by, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`,
		assigment.ID,
		assigment.TourDayID,
		assigment.GuideID,
		assigment.Action,
		assigment.ChangedBy,
		assigment.CreatedAt,
	)
	return err
}

func (r *Repository) AddTourStatus(status TourDayStatusHistory) error {
	_, err := r.db.Exec(`
		INSERT INTO tour_day_status_history (
			id, tour_day_id, previous_status, new_status, changed_by, reason, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		status.ID,
		status.TourDayID,
		status.PreviousStatus,
		status.NewStatus,
		status.ChangedBy,
		status.Reason,
		status.CreatedAt,
	)
	return err
}
