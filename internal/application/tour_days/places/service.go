package placesapp

import (
	"inariops/internal/modules/guides/zones"
	"inariops/internal/modules/itinerary/places"

	"github.com/google/uuid"
)

type Service struct {
	placesRepository *places.Repository
	zonesRepository  *zones.Repository
}

func NewService(placesRepository *places.Repository, zonesRepository *zones.Repository) *Service {
	return &Service{
		placesRepository: placesRepository,
		zonesRepository:  zonesRepository,
	}
}

func (s *Service) GetPlaceByID(placeID string) (*PlaceWithZone, error) {
	place, err := s.placesRepository.GetPlaceByID(placeID)
	if err != nil {
		return nil, err
	}

	zone, err := s.zonesRepository.GetZoneByID(*place.ZoneID)
	if err != nil {
		return nil, err
	}

	return &PlaceWithZone{
		Place: place,
		Zone:  zone,
	}, nil
}

func (s *Service) GetAllPlaces() ([]*PlaceWithZone, error) {
	places, err := s.placesRepository.GetAllPlaces()
	if err != nil {
		return nil, err
	}

	var placesWithZones []*PlaceWithZone
	for _, place := range places {
		zone, err := s.zonesRepository.GetZoneByID(*place.ZoneID)
		if err != nil {
			return nil, err
		}
		placesWithZones = append(placesWithZones, &PlaceWithZone{
			Place: &place,
			Zone:  zone,
		})
	}

	return placesWithZones, nil
}

func (s *Service) CreatePlace(place *places.CreatePlaceRequest) error {
	newPlace := &places.Place{
		ID:          uuid.New().String(),
		Name:        place.Name,
		Description: place.Description,
		ZoneID:      place.ZoneID,
		IsActive:    true,
	}

	return s.placesRepository.CreatePlace(*newPlace)
}

func (s *Service) UpdatePlace(placeID string, place *places.CreatePlaceRequest) error {
	existingPlace, err := s.placesRepository.GetPlaceByID(placeID)
	if err != nil {
		return err
	}

	existingPlace.Name = place.Name
	existingPlace.Description = place.Description
	existingPlace.ZoneID = place.ZoneID
	existingPlace.IsActive = place.IsActive

	err = s.placesRepository.UpdatePlace(*existingPlace)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeletePlace(placeID string) error {
	err := s.placesRepository.DeletePlace(placeID)
	if err != nil {
		return err
	}

	return nil
}
