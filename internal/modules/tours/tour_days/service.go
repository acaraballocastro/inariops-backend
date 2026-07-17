package tourdays

import (
	"inariops/internal/domain"
	"inariops/internal/modules/guides"
	tours "inariops/internal/modules/tours/shared"
	"inariops/internal/shared/errors"
	"inariops/internal/shared/logger"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo      *Repository
	guideRepo *guides.Repository
}

func NewService(repo *Repository, guideRepo *guides.Repository) *Service {
	return &Service{repo: repo, guideRepo: guideRepo}
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

func (s *Service) AssignGuide(tourDays []string, guideID string) error {
	logger.Info("AssignGuide: assigning guide to tour days")

	if len(tourDays) == 0 {
		logger.Error("AssignGuide: invalid input - at least one tour day ID is required")
		return errors.ErrInvalidInput
	}

	if guideID == "" {
		logger.Error("AssignGuide: invalid input - guide ID is required")
		return errors.ErrInvalidInput
	}

	if guideID == "" {
		logger.Error("AssignGuide: invalid input - guide ID is required")
		return errors.ErrInvalidInput
	}

	guide, err := s.guideRepo.GetGuideByID(guideID)
	if err != nil {
		logger.Error("AssignGuide: failed to get guide: %v", err)
		return errors.ErrGuideNotFound
	}

	for _, tourDayID := range tourDays {
		tourDay, err := s.repo.GetTourDayByID(tourDayID)
		if err != nil {
			logger.Error("AssignGuide: failed to get tour day: %v", err)
			return errors.ErrTourDayNotFound
		}

		if tourDay.Status == domain.RESERVATION_CANCELLED {
			logger.Error("AssignGuide: cannot assign guide to a cancelled tour day")
			return errors.ErrTourDayCancelled
		}

		if tourDay.GuideID != nil {
			logger.Error("AssignGuide: tour day already has a guide assigned")
			return errors.ErrGuideAlreadyAssigned
		}

		err = s.repo.AssignGuide(tourDayID, guide.ID)
		if err != nil {
			logger.Error("AssignGuide: failed to assign guide to tour day: %v", err)
			return err
		}
	}

	return nil
}

func (s *Service) UnassignGuide(tourDays []string, guideID string) error {
	logger.Info("UnassignGuide: unassigning guide from tour day")

	if len(tourDays) == 0 {
		logger.Error("UnassignGuide: invalid input - at least one tour day ID is required")
		return errors.ErrInvalidInput
	}

	if guideID == "" {
		logger.Error("UnassignGuide: invalid input - guide ID is required")
		return errors.ErrInvalidInput
	}

	guide, err := s.guideRepo.GetGuideByID(guideID)
	if err != nil {
		logger.Error("UnassignGuide: failed to get guide: %v", err)
		return errors.ErrGuideNotFound
	}

	for _, tourDayID := range tourDays {
		tourDay, err := s.repo.GetTourDayByID(tourDayID)
		if err != nil {
			logger.Error("UnassignGuide: failed to get tour day: %v", err)
			return errors.ErrTourDayNotFound
		}

		if tourDay.GuideID == nil || *tourDay.GuideID != guide.ID {
			logger.Error("UnassignGuide: guide is not assigned to this tour day")
			return errors.ErrGuideNotAssigned
		}

		err = s.repo.UnassignGuide(tourDayID)
		if err != nil {
			logger.Error("UnassignGuide: failed to unassign guide from tour day: %v", err)
			return err
		}
	}
	return nil

}

//TODO: Add a Worker to handle automatic changing of status
