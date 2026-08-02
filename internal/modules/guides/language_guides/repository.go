package languageguides

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateLanguageGuide(languageGuide *LanguageGuide) (*LanguageGuide, error) {
	_, err := r.db.Exec(`
		INSERT INTO guide_languages (
			guide_id,
			language_id
		)
		VALUES ($1, $2)
	`, languageGuide.GuideID, languageGuide.LanguageID)
	if err != nil {
		return nil, err
	}
	return languageGuide, nil
}

func (r *Repository) GetLanguageGuideByGuideID(guideID string) ([]*LanguageGuide, error) {
	rows, err := r.db.Query(`
		SELECT
			guide_id,
			language_id
		FROM guide_languages WHERE guide_id = $1
	`, guideID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var languageGuides []*LanguageGuide
	for rows.Next() {
		lg := &LanguageGuide{}
		if err := rows.Scan(&lg.GuideID, &lg.LanguageID); err != nil {
			return nil, err
		}
		languageGuides = append(languageGuides, lg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return languageGuides, nil
}

func (r *Repository) IsLanguageGuideExists(guideID, languageID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM guide_languages
			WHERE guide_id = $1 AND language_id = $2
		)
	`, guideID, languageID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) DeleteLanguageGuide(guideID, languageID string) error {
	_, err := r.db.Exec(`
		DELETE FROM guide_languages WHERE guide_id = $1 AND language_id = $2
	`, guideID, languageID)
	return err
}
