package agencies

import (
	"database/sql"
	"fmt"
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

func (r *Repository) CreateAgency(agency *Agency) (*Agency, error) {
	err := r.db.QueryRow(`
		INSERT INTO agencies (
			id,
			name,
			representative,
			email,
			is_active,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING
			id,
			name,
			representative,
			email,
			is_active,
			created_at,
			updated_at
	`, agency.ID, agency.Name, agency.Representative, agency.Email, agency.IsActive).Scan(
		&agency.ID,
		&agency.Name,
		&agency.Representative,
		&agency.Email,
		&agency.IsActive,
		&agency.CreatedAt,
		&agency.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return agency, nil
}

func (r *Repository) GetAgencyByID(id string) (*Agency, error) {
	agency := &Agency{}

	err := r.db.QueryRow(`
		SELECT
			id,
			name,
			representative,
			email,
			is_active,
			created_at,
			updated_at
		FROM agencies
		WHERE id = $1
	`, id).Scan(
		&agency.ID,
		&agency.Name,
		&agency.Representative,
		&agency.Email,
		&agency.IsActive,
		&agency.CreatedAt,
		&agency.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return agency, nil
}

func (r *Repository) UpdateAgency(agency *Agency) error {
	_, err := r.db.Exec(`
		UPDATE agencies
		SET
			name = $1,
			representative = $2,
			email = $3,
			is_active = $4,
			updated_at = NOW()
		WHERE id = $5
	`, agency.Name, agency.Representative, agency.Email, agency.IsActive, agency.ID)

	return err
}

func (r *Repository) DeleteAgency(id string) error {
	_, err := r.db.Exec(`
		UPDATE agencies
		SET
			is_active = FALSE,
			updated_at = NOW()
		WHERE id = $1
	`, id)

	return err
}

func (r *Repository) ListAgencies(req ListAgenciesRequest) ([]Agency, error) {
	query := `
		SELECT
			id,
			name,
			representative,
			email,
			is_active,
			created_at,
			updated_at
		FROM agencies
		WHERE 1=1
	`

	args := []any{}
	index := 1

	if req.Name != nil {
		query += fmt.Sprintf(" AND name ILIKE $%d", index)
		args = append(args, "%"+*req.Name+"%")
		index++
	}

	if req.Email != nil {
		query += fmt.Sprintf(" AND email = $%d", index)
		args = append(args, *req.Email)
		index++
	}

	if req.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", index)
		args = append(args, *req.IsActive)
		index++
	}

	query += " ORDER BY name ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	agencies := make([]Agency, 0)

	for rows.Next() {
		var agency Agency

		err := rows.Scan(
			&agency.ID,
			&agency.Name,
			&agency.Representative,
			&agency.Email,
			&agency.IsActive,
			&agency.CreatedAt,
			&agency.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		agencies = append(agencies, agency)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return agencies, nil
}

func (r *Repository) GetAgencyByEmail(email string) (*Agency, error) {
	agency := &Agency{}

	err := r.db.QueryRow(`
		SELECT
			id,
			name,
			representative,
			email,
			is_active,
			created_at,
			updated_at
		FROM agencies
		WHERE email = $1
	`, email).Scan(
		&agency.ID,
		&agency.Name,
		&agency.Representative,
		&agency.Email,
		&agency.IsActive,
		&agency.CreatedAt,
		&agency.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return agency, nil
}
