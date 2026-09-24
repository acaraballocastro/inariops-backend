package guides

import (
	"database/sql"
	"inariops/internal/domain"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) WithTx(tx *sql.Tx) *Service {
	return &Service{
		repo: s.repo.WithTx(tx),
	}
}

func (s *Service) CreateGuide(guide domain.Guide) error {
	return s.repo.CreateGuide(guide)
}

func (s *Service) GetGuideByID(id string) (domain.Guide, error) {
	return s.repo.GetGuideByID(id)
}

func (s *Service) GetGuideByUserID(userID string) (domain.Guide, error) {
	return s.repo.GetGuideByUserID(userID)
}

func (s *Service) GetAllGuides() ([]domain.Guide, error) {
	return s.repo.GetAllGuides()
}

func (s *Service) UpdateGuide(guide domain.Guide) error {
	return s.repo.UpdateGuide(guide)
}

func (s *Service) DeleteGuide(id string) error {
	return s.repo.DeleteGuideByUserID(id)
}
