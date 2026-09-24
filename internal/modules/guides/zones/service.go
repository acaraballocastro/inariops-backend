package zones

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

func (s *Service) CreateZone(zone *CreateZoneRequest) (*Zone, error) {
	newzone := &Zone{
		ID:        uuid.New().String(),
		Name:      zone.Name,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.ZoneValidations(newzone); err != nil {
		return nil, err
	}

	return s.repo.CreateZone(newzone)
}

func (s *Service) GetZoneByName(zoneName string) (*Zone, error) {
	return s.repo.GetZoneByName(zoneName)
}

func (s *Service) UpdateZone(zone *UpdateZoneRequest) error {
	existingZone, err := s.repo.GetZoneByID(zone.ID)
	if err != nil {
		return err
	}
	if existingZone == nil {
		return appErrors.ErrZoneNotFound
	}

	existingZone.UpdatedAt = time.Now()
	existingZone.Name = zone.Name
	existingZone.IsActive = zone.IsActive

	return s.repo.UpdateZone(existingZone)
}

func (s *Service) DeleteZone(zoneName string) error {
	return s.repo.DeleteZone(zoneName)
}

func (s *Service) ListZones() ([]*Zone, error) {
	return s.repo.ListZones()
}

func (s *Service) ZoneValidations(zone *Zone) error {
	existingZone, err := s.repo.GetZoneByName(zone.Name)
	if err != nil {
		return err
	}

	if existingZone != nil {
		return appErrors.ErrInvalidZone
	}

	return nil
}
