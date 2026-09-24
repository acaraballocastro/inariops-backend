package zoneguides

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

func (r *Repository) CreateZoneGuide(zoneGuide *ZoneGuide) (*ZoneGuide, error) {
	_, err := r.db.Exec(`
		INSERT INTO guide_zones (
			guide_id,
			zone_id
		)
		VALUES ($1, $2)
	`, zoneGuide.GuideID, zoneGuide.ZoneID)
	if err != nil {
		return nil, err
	}
	return zoneGuide, nil
}

func (r *Repository) GetZoneGuideByGuideID(guideID string) ([]*ZoneGuide, error) {
	rows, err := r.db.Query(`
		SELECT
			guide_id,
			zone_id
		FROM guide_zones WHERE guide_id = $1
	`, guideID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var zoneGuides []*ZoneGuide
	for rows.Next() {
		zg := &ZoneGuide{}
		if err := rows.Scan(&zg.GuideID, &zg.ZoneID); err != nil {
			return nil, err
		}
		zoneGuides = append(zoneGuides, zg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return zoneGuides, nil
}

func (r *Repository) IsZoneGuideExists(guideID, zoneID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM guide_zones
			WHERE guide_id = $1 AND zone_id = $2
		)
	`, guideID, zoneID).Scan(&exists)

	return exists, err
}

func (r *Repository) DeleteZoneGuide(guideID, zoneID string) error {
	_, err := r.db.Exec(`
		DELETE FROM guide_zones WHERE guide_id = $1 AND zone_id = $2
	`, guideID, zoneID)

	return err
}
