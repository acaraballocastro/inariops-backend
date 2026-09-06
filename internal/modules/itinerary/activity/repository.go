package activity

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateActivity(activity Activity) error {
	_, err := r.db.Exec(`
		INSERT INTO activities (
			id,
			name,
			description,
			is_active,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		activity.ID,
		activity.Name,
		activity.Description,
		activity.IsActive,
		activity.CreatedAt,
		activity.UpdatedAt,
	)

	return err
}

func (r *Repository) GetActivityByID(id string) (Activity, error) {
	activity := Activity{}

	err := r.db.QueryRow(`
		SELECT
			id,
			name,
			description,
			is_active,
			created_at,
			updated_at
		FROM activities
		WHERE id = $1
	`, id).Scan(
		&activity.ID,
		&activity.Name,
		&activity.Description,
		&activity.IsActive,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)

	if err != nil {
		return Activity{}, err
	}

	return activity, nil
}

func (r *Repository) GetAllActivities() ([]Activity, error) {
	rows, err := r.db.Query(`
		SELECT
			id,
			name,
			description,
			is_active,
			created_at,
			updated_at
		FROM activities
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		var activity Activity
		err := rows.Scan(
			&activity.ID,
			&activity.Name,
			&activity.Description,
			&activity.IsActive,
			&activity.CreatedAt,
			&activity.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

func (r *Repository) UpdateActivity(activity Activity) error {
	_, err := r.db.Exec(`
		UPDATE activities SET
			name = $1,
			description = $2,
			is_active = $3,
			updated_at = $4
		WHERE id = $5
	`,
		activity.Name,
		activity.Description,
		activity.IsActive,
		activity.UpdatedAt,
		activity.ID,
	)

	return err
}

func (r *Repository) DeleteActivity(id string) error {
	_, err := r.db.Exec(`
		DELETE FROM activities WHERE id = $1
	`, id)

	return err
}
