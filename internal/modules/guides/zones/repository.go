package zones

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

func (r *Repository) CreateZone(zone *Zone) (*Zone, error) {
	_, err := r.db.Exec(`
		INSERT INTO zones (id, name, type, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, zone.ID, zone.Name, zone.Type, zone.IsActive, zone.CreatedAt, zone.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return zone, nil
}

func (r *Repository) GetZoneByName(zoneName string) (*Zone, error) {
	var zone Zone
	err := r.db.QueryRow(`
		SELECT id, name, type, is_active, created_at, updated_at
		FROM zones WHERE name = $1
	`, zoneName).Scan(&zone.ID, &zone.Name, &zone.Type, &zone.IsActive, &zone.CreatedAt, &zone.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &zone, nil
}

func (r *Repository) UpdateZone(zone *Zone) error {
	_, err := r.db.Exec(`
		UPDATE zones SET name = $1, is_active = $2, updated_at = $3 WHERE type = $4
	`, zone.Name, zone.IsActive, zone.UpdatedAt, zone.Type)

	return err
}

func (r *Repository) DeleteZone(zoneName string) error {
	_, err := r.db.Exec(`
		DELETE FROM zones WHERE name = $1
	`, zoneName)

	return err
}

func (r *Repository) ListZones() ([]*Zone, error) {
	rows, err := r.db.Query(`
		SELECT id, name, type, is_active, created_at, updated_at
		FROM zones
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var zones []*Zone
	for rows.Next() {
		var zone Zone
		if err := rows.Scan(&zone.ID, &zone.Name, &zone.Type, &zone.IsActive, &zone.CreatedAt, &zone.UpdatedAt); err != nil {
			return nil, err
		}
		zones = append(zones, &zone)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return zones, nil
}

func (r *Repository) GetZoneByID(zoneID string) (*Zone, error) {
	var zone Zone
	err := r.db.QueryRow(`
		SELECT id, name, type, is_active, created_at, updated_at
		FROM zones WHERE id = $1
	`, zoneID).Scan(&zone.ID, &zone.Name, &zone.Type, &zone.IsActive, &zone.CreatedAt, &zone.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &zone, nil
}
