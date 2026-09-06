package tourdaysapp

import (
	"context"
	"inariops/internal/domain"
	"inariops/internal/modules/guides"
	"inariops/internal/modules/guides/availabilities"
	tours "inariops/internal/modules/tours/shared"
	tourdays "inariops/internal/modules/tours/tour_days"
	"inariops/internal/shared/errors"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo             *Repository
	tourDayRepo      *tourdays.Repository
	guideRepo        *guides.Repository
	availabilityRepo *availabilities.Repository
}

func NewService(tourDayRepo *tourdays.Repository, guideRepo *guides.Repository, availabilityRepo *availabilities.Repository, assigmentRepo *Repository) *Service {
	return &Service{
		tourDayRepo:      tourDayRepo,
		guideRepo:        guideRepo,
		availabilityRepo: availabilityRepo,
		repo:             assigmentRepo,
	}
}

func (s *Service) AssignGuide(tourDayIDs []string, guideID string) error {
	tourdays := make([]tours.TourDay, 0, len(tourDayIDs))
	for _, tourDayID := range tourDayIDs {
		tourDay, err := s.tourDayRepo.GetTourDayByID(tourDayID)
		if err != nil {
			return err
		}
		tourdays = append(tourdays, tourDay)
	}

	guide, err := s.guideRepo.GetGuideByID(guideID)
	if err != nil {
		return err
	}

	if len(tourdays) == 1 {
		err = s.tourDayRepo.AssignGuide(tourdays[0].ID, guide.ID)
		if err != nil {
			return err
		}
		err = s.SetAssignmentStatusHistoryRegistry(tourdays[0].ID, guide.ID, "ASSIGNED", "GUIDE")
		return nil
	}

	err = s.tourDayRepo.AssignGuideToMultipleTourDays(tourDayIDs, guide.ID)
	for _, tourday := range tourdays {
		err = s.SetAssignmentStatusHistoryRegistry(tourday.ID, guide.ID, "ASSIGNED", "GUIDE")
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UnassignGuide(tourDayIDs []string, guideID string) error {
	tourdays := make([]tours.TourDay, 0, len(tourDayIDs))
	for _, tourDayID := range tourDayIDs {
		tourDay, err := s.tourDayRepo.GetTourDayByID(tourDayID)
		if err != nil {
			return err
		}
		tourdays = append(tourdays, tourDay)
	}

	if len(tourdays) == 1 {
		err := s.tourDayRepo.UnassignGuide(tourdays[0].ID)
		if err != nil {
			return err
		}
		err = s.SetAssignmentStatusHistoryRegistry(tourdays[0].ID, guideID, "UNASSIGNED", "GUIDE")
		if err != nil {
			return err
		}

		return nil
	}

	for _, tourday := range tourdays {
		err := s.tourDayRepo.UnassignGuide(tourday.ID)
		if err != nil {
			return err
		}
		err = s.SetAssignmentStatusHistoryRegistry(tourday.ID, guideID, "UNASSIGNED", "GUIDE")
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) SetAssignmentStatusHistoryRegistry(tourDayID string, guideID string, action string, changedBy string) error {
	tourDay, err := s.tourDayRepo.GetTourDayByID(tourDayID)
	if err != nil {
		return err
	}

	guide, err := s.guideRepo.GetGuideByID(guideID)
	if err != nil {
		return err
	}

	assigment := TourDayAssigmentHistory{
		ID:        uuid.New().String(),
		TourDayID: tourDay.ID,
		GuideID:   guide.ID,
		Action:    action,
		ChangedBy: changedBy,
		CreatedAt: tourDay.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	err = s.repo.AddAssignmentStatus(assigment)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SetTourDayHistory(tourDayID string, newStatus string, changedBy string, reason string) error {
	tourDay, err := s.tourDayRepo.GetTourDayByID(tourDayID)
	if err != nil {
		return err
	}

	history := TourDayStatusHistory{
		ID:             uuid.New().String(),
		TourDayID:      tourDay.ID,
		PreviousStatus: string(tourDay.Status),
		NewStatus:      newStatus,
		ChangedBy:      changedBy,
		Reason:         reason,
		CreatedAt:      tourDay.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	err = s.repo.AddTourStatus(history)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Add a Worker to handle automatic changing of status
func (s *Service) ProcessExpiredGuideAssignments(ctx context.Context) (domain.Result, error) {
	var result domain.Result

	tours, err := s.tourDayRepo.GetExpiredGuideAssignments(ctx)
	if err != nil {
		return result, err
	}

	for _, tour := range tours {
		result.Found++

		err := s.tourDayRepo.UpdateTourDayStatus(tour.ID, domain.RESERVATION_PAYMENT_PENDING)
		if err != nil {
			result.Failed++
			continue
		}

		err = s.SetTourDayHistory(tour.ID, string(domain.RESERVATION_PAYMENT_PENDING), "SYSTEM", "Guide assignment expired")

		result.Processed++

	}

	return result, nil
}

func (s *Service) TourDaysAvailableForGuide(guideID string) ([]tours.TourDay, error) {
	guide, err := s.guideRepo.GetGuideByID(guideID)
	if err != nil {
		return nil, err
	}

	tourDays, err := s.tourDayRepo.GetTourDaysAvailableForGuide(guideID)
	if err != nil {
		return nil, err
	}

	blockedRanges, err := s.availabilityRepo.GetAvailabilitiesByGuideID(guideID)
	if err != nil {
		return nil, err
	}

	available := make([]tours.TourDay, 0, len(tourDays))
	dailyAssignments := make(map[string]int)
	for _, tourDay := range tourDays {
		if availabilityCoversDate(blockedRanges, tourDay.StartDateTime) {
			continue
		}

		dateKey := tourDay.StartDateTime.Format("2006-01-02")
		if _, counted := dailyAssignments[dateKey]; !counted {
			count, err := s.tourDayRepo.CountGuideTourDaysByDate(guideID, tourDay.StartDateTime)
			if err != nil {
				return nil, err
			}
			dailyAssignments[dateKey] = count
		}
		if dailyAssignments[dateKey] >= guide.MaxToursPerDay {
			continue
		}
		available = append(available, tourDay)
	}

	return available, nil
}

func (s *Service) ConfirmTourDay(tourDayID string) error {
	tourDay, err := s.tourDayRepo.GetTourDayByID(tourDayID)
	if err != nil {
		return err
	}
	if tourDay.Status != domain.RESERVATION_GUIDE_PREASSIGNED {
		return errors.ErrInvalidInput
	}
	if err := s.tourDayRepo.UpdateTourDayStatus(tourDayID, domain.RESERVATION_GUIDE_CONFIRMED); err != nil {
		return err
	}
	return s.repo.AddTourStatus(TourDayStatusHistory{
		ID:             uuid.New().String(),
		TourDayID:      tourDay.ID,
		PreviousStatus: string(tourDay.Status),
		NewStatus:      string(domain.RESERVATION_GUIDE_CONFIRMED),
		ChangedBy:      "GUIDE",
		Reason:         "Guide confirmed tour",
		CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
	})
}

func availabilityCoversDate(ranges []*availabilities.Availability, date time.Time) bool {
	date = date.Truncate(24 * time.Hour)
	for _, availability := range ranges {
		if availability.StartDate == nil || availability.EndDate == nil {
			continue
		}
		start := availability.StartDate.Truncate(24 * time.Hour)
		end := availability.EndDate.Truncate(24 * time.Hour)
		if !date.Before(start) && !date.After(end) {
			return true
		}
	}
	return false
}
