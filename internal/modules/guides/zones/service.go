package zones

import (
	"fmt"
	"time"

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

	if _, isValid := s.ZoneValidations(newzone); !isValid {
		return nil, fmt.Errorf("invalid zone data")
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
		return fmt.Errorf("zone with ID %s not found", zone.ID)
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

func (s *Service) ZoneValidations(zone *Zone) (*Zone, bool) {
	existingZone, err := s.repo.GetZoneByName(zone.Name)
	if err != nil {
		return nil, false
	}
	if existingZone != nil {
		return existingZone, false
	}
	return nil, true
}
