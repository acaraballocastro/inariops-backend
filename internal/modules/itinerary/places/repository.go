package places

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreatePlace(place Place) error {
	_, err := r.db.Exec(`
		INSERT INTO places (
			id,
			name,
			description,
			zone_id,
			is_active,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		place.ID,
		place.Name,
		place.Description,
		place.ZoneID,
		place.IsActive,
		place.CreatedAt,
		place.UpdatedAt,
	)

	return err
}

func (r *Repository) GetPlaceByID(id string) (*Place, error) {
	place := Place{}

	err := r.db.QueryRow(`
		SELECT
			id,
			name,
			description,
			zone_id,
			is_active,
			created_at,
			updated_at
		FROM places
		WHERE id = $1
	`, id).Scan(
		&place.ID,
		&place.Name,
		&place.Description,
		&place.ZoneID,
		&place.IsActive,
		&place.CreatedAt,
		&place.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &place, nil
}

func (r *Repository) GetAllPlaces() ([]Place, error) {
	rows, err := r.db.Query(`
		SELECT
			id,
			name,
			description,
			zone_id,
			is_active,
			created_at,
			updated_at
		FROM places
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var places []Place
	for rows.Next() {
		var place Place
		err := rows.Scan(
			&place.ID,
			&place.Name,
			&place.Description,
			&place.ZoneID,
			&place.IsActive,
			&place.CreatedAt,
			&place.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		places = append(places, place)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return places, nil
}

func (r *Repository) UpdatePlace(place Place) error {
	_, err := r.db.Exec(`
		UPDATE places
		SET name = $1,
			description = $2,
			zone_id = $3,
			is_active = $4,
			updated_at = $5
		WHERE id = $6
	`,
		place.Name,
		place.Description,
		place.ZoneID,
		place.IsActive,
		place.UpdatedAt,
		place.ID,
	)

	return err
}

func (r *Repository) DeletePlace(id string) error {
	_, err := r.db.Exec(`
		DELETE FROM places WHERE id = $1
	`, id)

	return err
}
