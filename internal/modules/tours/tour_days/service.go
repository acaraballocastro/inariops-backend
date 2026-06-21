package tourdays

import (
	"inariops/internal/domain"
	tours "inariops/internal/modules/tours/shared"
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

func (s *Service) GetTourDaysByTourID(tourID string) (tours.TourDay, error) {
	tourDay, err := s.repo.GetTourDayByID(tourID)
	if err != nil {
		logger.Error("GetTourDaysByTourID: failed to get tour day for tour ID %s: %v", tourID, err)
		return tours.TourDay{}, err
	}

	return tourDay, nil
}

func (s *Service) GetTourDaysByReservationID(reservationID string) ([]tours.TourDay, error) {
	logger.Info("reservattionID %s", reservationID)
	tourDays, err := s.repo.GetTourDaysByReservationID(reservationID)
	if err != nil {
		logger.Error("GetTourDaysByReservationID: failed to get tour days for reservation ID %s: %v", reservationID, err)
		return nil, err
	}

	return tourDays, nil
}

func (s *Service) CreateTourDay(tourDayInput CreateTourDayInput) error {
	logger.Info("CreateTourDay: creating tour day")
	tourDay := tours.TourDay{
		ID:                     uuid.New().String(),
		ReservationID:          tourDayInput.ReservationID,
		Title:                  tourDayInput.Title,
		StartDateTime:          tourDayInput.StartDateTime,
		Duration:               tourDayInput.Duration,
		ZoneID:                 tourDayInput.ZoneID,
		GuideID:                nil,
		Status:                 domain.RESERVATION_PENDING_ASSIGNMENT,
		PeopleCount:            tourDayInput.PeopleCount,
		Remuneration:           tourDayInput.Remuneration,
		MeetingPoint:           tourDayInput.MeetingPoint,
		GuideHotelID:           tourDayInput.GuideHotelID,
		CustomerHotelID:        tourDayInput.CustomerHotelID,
		GuideLiabilityNotes:    tourDayInput.GuideLiabilityNotes,
		CustomerLiabilityNotes: tourDayInput.CustomerLiabilityNotes,
		VoucherStatus:          domain.VOUCHER_NOT_GENERATED,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	if tourDay.ReservationID == "" {
		logger.Error("CreateTourDay: invalid input - reservation ID is required")
		return errors.ErrInvalidInput
	}

	return s.repo.CreateTourDay(tourDay)
}

func (s *Service) UpdateTourDay(tourDayInput UpdateTourDayInput) error {
	logger.Info("UpdateTourDay: updating tour day")
	tourDay := tours.TourDay{
		ID:                     tourDayInput.ID,
		ReservationID:          tourDayInput.ReservationID,
		Title:                  tourDayInput.Title,
		StartDateTime:          tourDayInput.StartDateTime,
		Duration:               tourDayInput.Duration,
		ZoneID:                 tourDayInput.ZoneID,
		GuideID:                nil,
		Status:                 tourDayInput.Status,
		PeopleCount:            tourDayInput.PeopleCount,
		Remuneration:           tourDayInput.Remuneration,
		MeetingPoint:           tourDayInput.MeetingPoint,
		GuideHotelID:           tourDayInput.GuideHotelID,
		CustomerHotelID:        tourDayInput.CustomerHotelID,
		GuideLiabilityNotes:    tourDayInput.GuideLiabilityNotes,
		CustomerLiabilityNotes: tourDayInput.CustomerLiabilityNotes,
		VoucherStatus:          tourDayInput.VoucherStatus,
		UpdatedAt:              time.Now(),
	}

	if tourDay.ReservationID == "" {
		logger.Error("UpdateTourDay: invalid input - reservation ID is required")
		return errors.ErrInvalidInput
	}

	err := s.repo.UpdateTourDay(tourDay)
	if err != nil {
		logger.Error("UpdateTourDay: failed to update tour day", "tour day ID %s: %v", tourDay.ID, err)
		return err
	}

	return nil
}

func (s *Service) CancelTourDay(id string) error {
	logger.Info("CancelTourDay: canceling tour day")
	return s.repo.CancelTourDay(id)
}
