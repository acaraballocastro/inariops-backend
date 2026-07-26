package tourdays

import (
	"context"
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

	err := s.repo.CreateTourDay(tourDay)
	if err != nil {
		logger.Error("CreateTourDay: failed to create tour day for reservation ID %s: %v", tourDay.ReservationID, err)
		return err
	}

	return nil
}

func (s *Service) UpdateTourDay(input UpdateTourDayInput) error {
	logger.Info("UpdateTourDay: updating tour day")
	logger.Info("UpdateTourDay: input: %+v", input)

	tourDay, err := s.repo.GetTourDayByID(input.ID)
	if err != nil {
		logger.Error("UpdateTourDay: failed to get tour day: %v", err)
		return err
	}

	if input.Title != nil {
		tourDay.Title = input.Title
	}

	if !input.StartDateTime.IsZero() {
		tourDay.StartDateTime = input.StartDateTime
	}

	if input.Duration != nil {
		tourDay.Duration = input.Duration
	}

	if input.ZoneID != nil {
		tourDay.ZoneID = input.ZoneID
	}

	if input.GuideID != nil {
		tourDay.GuideID = input.GuideID
	}

	if input.PeopleCount != nil {
		tourDay.PeopleCount = input.PeopleCount
	}

	if input.MeetingPoint != nil {
		tourDay.MeetingPoint = input.MeetingPoint
	}

	if input.GuideHotelID != nil {
		tourDay.GuideHotelID = input.GuideHotelID
	}

	if input.CustomerHotelID != nil {
		tourDay.CustomerHotelID = input.CustomerHotelID
	}

	if input.GuideLiabilityNotes != nil {
		tourDay.GuideLiabilityNotes = input.GuideLiabilityNotes
	}

	if input.CustomerLiabilityNotes != nil {
		tourDay.CustomerLiabilityNotes = input.CustomerLiabilityNotes
	}

	if input.Remuneration != nil {
		tourDay.Remuneration = input.Remuneration
	}

	if input.Status != "" {
		tourDay.Status = input.Status
	}

	if input.VoucherStatus != "" {
		tourDay.VoucherStatus = input.VoucherStatus
	}

	tourDay.UpdatedAt = time.Now()

	return s.repo.UpdateTourDay(tourDay)
}
func (s *Service) CancelTourDay(id string) error {
	logger.Info("CancelTourDay: canceling tour day")
	return s.repo.CancelTourDay(id)
}

// TODO: Add a Worker to handle automatic changing of status
func (s *Service) ProcessExpiredGuideAssignments(ctx context.Context) (domain.Result, error) {
	var result domain.Result

	tours, err := s.repo.GetExpiredGuideAssignments(ctx)
	if err != nil {
		return result, err
	}

	for _, tour := range tours {
		result.Found++

		err := s.repo.UpdateTourDayStatus(tour.ID, domain.RESERVATION_GUIDE_CONFIRMED)
		if err != nil {
			result.Failed++
			continue
		}

		result.Processed++

	}

	return result, nil
}
