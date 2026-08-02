package languages

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateLanguage(language *Language) (*Language, error) {
	_, err := r.db.Exec(`
		INSERT INTO languages (
			id,
			name,
			code,
			is_active,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, language.ID, language.Name, language.Code, language.IsActive, language.CreatedAt, language.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return language, nil
}

func (r *Repository) GetLanguageByCode(code string) (*Language, error) {
	language := &Language{}
	err := r.db.QueryRow(`
		SELECT
			id,
			name,
			code,
			is_active,
			created_at,
			updated_at
		FROM languages WHERE code = $1
	`, code).Scan(&language.ID, &language.Name, &language.Code, &language.IsActive, &language.CreatedAt, &language.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Language not found
		}
		return nil, err // Other error occurred
	}

	return language, nil
}

func (r *Repository) UpdateLanguage(language *Language) error {
	_, err := r.db.Exec(`
		UPDATE languages SET
			name = $1,
			code = $2,
			is_active = $3,
			updated_at = $4
		WHERE id = $5
	`, language.Name, language.Code, language.IsActive, language.UpdatedAt, language.ID)

	return err
}

func (r *Repository) DeleteLanguage(id string) error {
	_, err := r.db.Exec(`
		DELETE FROM languages WHERE id = $1
	`, id)
	return err
}

func (r *Repository) DeactivateLanguage(id string) error {
	_, err := r.db.Exec(`
		UPDATE languages SET is_active = false, updated_at = NOW() WHERE id = $1
	`, id)
	return err
}

func (r *Repository) ListLanguages() ([]*Language, error) {
	rows, err := r.db.Query(`
		SELECT
			id,
			name,
			code,
			is_active,
			created_at,
			updated_at
		FROM languages
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var languages []*Language
	for rows.Next() {
		language := &Language{}
		err := rows.Scan(&language.ID, &language.Name, &language.Code, &language.IsActive, &language.CreatedAt, &language.UpdatedAt)
		if err != nil {
			return nil, err
		}
		languages = append(languages, language)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return languages, nil
}
