package languages

import (
	"inariops/internal/shared/errors"
	"inariops/internal/shared/logger"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateLanguage(language *Language) (*Language, error) {
	language.ID = uuid.New().String()
	language.IsActive = true
	language.CreatedAt = time.Now()
	language.UpdatedAt = time.Now()

	if err := s.LanguageValidations(language); err != nil {
		return nil, err
	}

	logger.Info("Creating new Language with Code: %s", language.Code)

	return s.repo.CreateLanguage(language)
}

func (s *Service) GetLanguageByCode(code string) (*Language, error) {
	return s.repo.GetLanguageByCode(code)
}

func (s *Service) UpdateLanguage(language *Language) error {
	existingLanguage, err := s.repo.GetLanguageByCode(language.Code)
	if err != nil {
		return err
	}
	if existingLanguage == nil {
		return errors.ErrLanguageNotFound
	}

	existingLanguage.UpdatedAt = time.Now()
	existingLanguage.Name = language.Name
	existingLanguage.IsActive = language.IsActive

	return s.repo.UpdateLanguage(existingLanguage)
}

func (s *Service) DeleteLanguage(code string) error {
	return s.repo.DeleteLanguage(code)
}

func (s *Service) ListLanguages() ([]*Language, error) {
	return s.repo.ListLanguages()
}

func (s *Service) LanguageValidations(language *Language) error {
	existingLanguage, err := s.repo.GetLanguageByCode(language.Code)
	if err != nil {
		return err
	}

	if existingLanguage != nil {
		return errors.ErrInvalidLanguage
	}

	return nil
}
