package guides

import (
	"database/sql"
	"inariops/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateGuide(guide domain.Guide) error {
	_, err := r.db.Exec(`
			INSERT INTO guides (
				user_id, max_tours_per_day, created_at
			)
			VALUES (
				$1, $2, $3
			)
		`,
		guide.UserID, guide.MaxToursPerDay, guide.CreatedAt,
	)

	return err
}

func (r *Repository) GetGuideByID(id string) (domain.Guide, error) {
	var guide domain.Guide

	err := r.db.QueryRow(`
			SELECT id,
				user_id,
				max_tours_per_day,
				created_at
			FROM guides WHERE id = $1
		`, id).Scan(&guide.ID, &guide.UserID, &guide.MaxToursPerDay, &guide.CreatedAt)

	return guide, err
}

func (r *Repository) GetGuideByUserID(userID string) (domain.Guide, error) {
	var guide domain.Guide

	err := r.db.QueryRow(`
			SELECT id,
				user_id,
				max_tours_per_day,
				created_at
			FROM guides WHERE user_id = $1
		`, userID).Scan(&guide.ID, &guide.UserID, &guide.MaxToursPerDay, &guide.CreatedAt)

	return guide, err
}

func (r *Repository) GetAllGuides() ([]domain.Guide, error) {
	rows, err := r.db.Query(`
			SELECT id,
				user_id,
				max_tours_per_day,
				created_at
			FROM guides
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guides []domain.Guide
	for rows.Next() {
		var guide domain.Guide
		err := rows.Scan(&guide.ID, &guide.UserID, &guide.MaxToursPerDay, &guide.CreatedAt)
		if err != nil {
			return nil, err
		}
		guides = append(guides, guide)
	}

	return guides, nil
}

func (r *Repository) UpdateGuide(guide domain.Guide) error {
	result, err := r.db.Exec(`
		UPDATE guides
		SET max_tours_per_day = $1
		WHERE user_id = $2
	`, guide.MaxToursPerDay, guide.UserID)

	if err != nil {
		return err
	}

	resultRowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if resultRowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) DeleteGuideByUserID(userID string) error {
	result, err := r.db.Exec(`
		DELETE FROM guides WHERE user_id = $1
	`, userID)

	if err != nil {
		return err
	}

	resultRowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if resultRowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
