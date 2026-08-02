package languages

import (
	"fmt"
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

	if _, isValid := s.LanguageValidations(language); !isValid {
		return nil, fmt.Errorf("invalid language data")
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
		return fmt.Errorf("language with code %s not found", language.Code)
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

func (s *Service) LanguageValidations(language *Language) (*Language, bool) {
	existingLanguage, err := s.repo.GetLanguageByCode(language.Code)
	if err != nil {
		return nil, false
	}
	if existingLanguage != nil {
		return nil, false
	}

	return nil, true
}
