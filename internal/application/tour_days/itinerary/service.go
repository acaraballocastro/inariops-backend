package itineraryapp

import (
	"inariops/internal/modules/itinerary"
	"inariops/internal/modules/itinerary/activity"
	"inariops/internal/modules/itinerary/places"
	tourdays "inariops/internal/modules/tours/tour_days"

	"github.com/google/uuid"
)

type Service struct {
	tourDaysRepository   *tourdays.Repository
	activitiesRepository *activity.Repository
	placesRepository     *places.Repository
	itineraryRepository  *itinerary.Repository
}

func NewService(tourDaysRepository *tourdays.Repository, activitiesRepository *activity.Repository, placesRepository *places.Repository, itineraryRepository *itinerary.Repository) *Service {
	return &Service{
		tourDaysRepository:   tourDaysRepository,
		activitiesRepository: activitiesRepository,
		placesRepository:     placesRepository,
		itineraryRepository:  itineraryRepository,
	}
}

type ItineraryItem struct {
	ID         string  `json:"id"`
	TourDayID  string  `json:"tour_day_id"`
	OrderIndex int     `json:"order_index"`
	StartTime  *string `json:"start_time"`
	EndTime    *string `json:"end_time"`
	Type       *string `json:"type"`
	PlaceID    *string `json:"place_id"`
	ActivityID *string `json:"activity_id"`
	Notes      *string `json:"notes"`
}

type CreateItineraryItemRequest struct {
	OrderIndex int     `json:"order_index"`
	StartTime  *string `json:"start_time"`
	EndTime    *string `json:"end_time"`
	Type       *string `json:"type"`
	PlaceID    *string `json:"place_id"`
	ActivityID *string `json:"activity_id"`
	Notes      *string `json:"notes"`
}

func (s *Service) CreateItineraryItem(tourDayID string, request itinerary.CreateItineraryItemRequest) error {
	// Check if the tour day exists
	_, err := s.tourDaysRepository.GetTourDayByID(tourDayID)
	if err != nil {
		return err
	}

	// Create a new itinerary item
	itineraryItem := itinerary.ItineraryItem{
		ID:         uuid.New().String(),
		TourDayID:  tourDayID,
		OrderIndex: request.OrderIndex,
		StartTime:  request.StartTime,
		EndTime:    request.EndTime,
		Type:       request.Type,
		PlaceID:    request.PlaceID,
		ActivityID: request.ActivityID,
		Notes:      request.Notes,
	}

	// Check if the activity exists if ActivityID is provided
	if itineraryItem.ActivityID != nil {
		_, err := s.activitiesRepository.GetActivityByID(*itineraryItem.ActivityID)
		if err != nil {
			return err
		}
	}

	// Check if the place exists if PlaceID is provided
	if itineraryItem.PlaceID != nil {
		_, err := s.placesRepository.GetPlaceByID(*itineraryItem.PlaceID)
		if err != nil {
			return err
		}
	}

	return s.itineraryRepository.CreateItineraryItem(itineraryItem)
}

func (s *Service) GetItineraryByTourDayID(tourDayID string) ([]itinerary.ItineraryItem, error) {
	// Check if the tour day exists
	_, err := s.tourDaysRepository.GetTourDayByID(tourDayID)
	if err != nil {
		return nil, err
	}

	return s.itineraryRepository.GetItineraryItemsByTourDayID(tourDayID)
}

func (s *Service) GetItineraryItemByID(itemID string) (itinerary.ItineraryItem, error) {
	// Check if the itinerary item exists
	return s.itineraryRepository.GetItineraryItemByID(itemID)
}

func (s *Service) UpdateItineraryItem(item string, req itinerary.CreateItineraryItemRequest) error {
	// Check if the itinerary item exists
	itineraryItem, err := s.itineraryRepository.GetItineraryItemByID(item)
	if err != nil {
		return err
	}

	// Update the itinerary item fields
	itineraryItem.OrderIndex = req.OrderIndex
	itineraryItem.StartTime = req.StartTime
	itineraryItem.EndTime = req.EndTime
	itineraryItem.Type = req.Type
	itineraryItem.PlaceID = req.PlaceID
	itineraryItem.ActivityID = req.ActivityID
	itineraryItem.Notes = req.Notes

	// Check if the activity exists if ActivityID is provided
	if itineraryItem.ActivityID != nil {
		_, err := s.activitiesRepository.GetActivityByID(*itineraryItem.ActivityID)
		if err != nil {
			return err
		}
	}

	// Check if the place exists if PlaceID is provided
	if itineraryItem.PlaceID != nil {
		_, err := s.placesRepository.GetPlaceByID(*itineraryItem.PlaceID)
		if err != nil {
			return err
		}
	}

	// Update the itinerary item
	return s.itineraryRepository.UpdateItineraryItem(itineraryItem)
}

func (s *Service) DeleteItineraryItem(itemID string) error {
	// Check if the itinerary item exists
	_, err := s.itineraryRepository.GetItineraryItemByID(itemID)
	if err != nil {
		return err
	}

	// Delete the itinerary item
	return s.itineraryRepository.DeleteItineraryItem(itemID)
}
