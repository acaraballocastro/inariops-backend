package agencies

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

func (s *Service) CreateAgency(agency *Agency) (*Agency, error) {
	agency.ID = uuid.New().String()
	agency.IsActive = true
	agency.CreatedAt = time.Now()
	agency.UpdatedAt = time.Now()

	if _, isValid := s.AgencyValidations(agency); !isValid {
		return nil, fmt.Errorf("invalid agency data")
	}

	createdAgency, err := s.repo.CreateAgency(agency)

	if err != nil {
		return nil, err
	}

	return createdAgency, nil
}

func (s *Service) GetAgencyByID(id string) (*Agency, error) {
	return s.repo.GetAgencyByID(id)
}

func (s *Service) UpdateAgency(agency *Agency) error {
	return s.repo.UpdateAgency(agency)
}

func (s *Service) DeleteAgency(id string) error {
	return s.repo.DeleteAgency(id)
}

func (s *Service) ListAgencies(req ListAgenciesRequest) ([]Agency, error) {
	return s.repo.ListAgencies(req)
}
func (s *Service) AgencyValidations(agency *Agency) (Agency, bool) {
	existingAgency, err := s.repo.GetAgencyByEmail(agency.Email)
	if err != nil {
		logger.Error("AgencyValidations", "failed to fetch agency by email: %v", err)
		return Agency{}, false
	}
	if existingAgency != nil {
		logger.Error("AgencyValidations", "agency with email %s already exists", agency.Email)
		return Agency{}, false
	}

	return Agency{}, true
}
