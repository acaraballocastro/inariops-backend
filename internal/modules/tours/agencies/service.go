package agencies

import (
	"time"

	appErrors "inariops/internal/shared/errors"

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

	if err := s.AgencyValidations(agency); err != nil {
		return nil, err
	}

	createdAgency, err := s.repo.CreateAgency(agency)

	if err != nil {
		return nil, err
	}

	return createdAgency, nil
}

func (s *Service) GetAgencyByID(id string) (*Agency, error) {
	agency, err := s.repo.GetAgencyByID(id)
	if err != nil {
		return nil, err
	}

	if agency == nil {
		return nil, appErrors.ErrAgencyNotFound
	}

	return agency, nil
}

func (s *Service) UpdateAgency(agency *Agency) error {
	existingAgency, err := s.repo.GetAgencyByID(agency.ID)
	if err != nil {
		return err
	}

	if existingAgency == nil {
		return appErrors.ErrAgencyNotFound
	}

	return s.repo.UpdateAgency(agency)
}

func (s *Service) DeleteAgency(id string) error {
	agency, err := s.repo.GetAgencyByID(id)
	if err != nil {
		return err
	}

	if agency == nil {
		return appErrors.ErrAgencyNotFound
	}

	return s.repo.DeleteAgency(id)
}

func (s *Service) ListAgencies(req ListAgenciesRequest) ([]Agency, error) {
	return s.repo.ListAgencies(req)
}

func (s *Service) AgencyValidations(agency *Agency) error {
	existingAgency, err := s.repo.GetAgencyByEmail(agency.Email)
	if err != nil {
		return err
	}

	if existingAgency != nil {
		return appErrors.ErrExistingAgency
	}

	return nil
}
