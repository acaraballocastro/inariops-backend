package tourdaysapp

import (
	"context"
	"database/sql"
	"inariops/internal/db"
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
	db               *sql.DB
	repo             *Repository
	tourDayRepo      *tourdays.Repository
	guideRepo        *guides.Repository
	availabilityRepo *availabilities.Repository
}

func NewService(
	database *sql.DB,
	tourDayRepo *tourdays.Repository,
	guideRepo *guides.Repository,
	availabilityRepo *availabilities.Repository,
	assigmentRepo *Repository,
) *Service {
	return &Service{
		db:               database,
		tourDayRepo:      tourDayRepo,
		guideRepo:        guideRepo,
		availabilityRepo: availabilityRepo,
		repo:             assigmentRepo,
	}
}

func (s *Service) AssignGuide(tourDayIDs []string, guideID string) error {
	return db.WithTransaction(s.db, func(tx *sql.Tx) error {

		tourDayRepo := s.tourDayRepo.WithTx(tx)
		guideRepo := s.guideRepo.WithTx(tx)
		assignmentRepo := s.repo.WithTx(tx)

		tourdays := make([]tours.TourDay, 0, len(tourDayIDs))

		for _, tourDayID := range tourDayIDs {
			tourDay, err := tourDayRepo.GetTourDayByID(tourDayID)
			if err != nil {
				return err
			}

			tourdays = append(tourdays, tourDay)
		}

		guide, err := guideRepo.GetGuideByID(guideID)
		if err != nil {
			return err
		}

		if len(tourdays) == 1 {
			if err := tourDayRepo.AssignGuide(
				tourdays[0].ID,
				guide.ID,
			); err != nil {
				return err
			}

			return s.setAssignmentStatusHistoryRegistry(
				tourdays[0].ID,
				guide.ID,
				"ASSIGNED",
				"GUIDE",
				tourDayRepo,
				guideRepo,
				assignmentRepo,
			)
		}

		if err := tourDayRepo.AssignGuideToMultipleTourDays(
			tourDayIDs,
			guide.ID,
		); err != nil {
			return err
		}

		for _, tourday := range tourdays {
			if err := s.setAssignmentStatusHistoryRegistry(
				tourday.ID,
				guide.ID,
				"ASSIGNED",
				"GUIDE",
				tourDayRepo,
				guideRepo,
				assignmentRepo,
			); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Service) UnassignGuide(tourDayIDs []string, guideID string) error {
	return db.WithTransaction(s.db, func(tx *sql.Tx) error {

		tourDayRepo := s.tourDayRepo.WithTx(tx)
		guideRepo := s.guideRepo.WithTx(tx)
		assignmentRepo := s.repo.WithTx(tx)

		tourdays := make([]tours.TourDay, 0, len(tourDayIDs))

		for _, tourDayID := range tourDayIDs {
			tourDay, err := tourDayRepo.GetTourDayByID(tourDayID)
			if err != nil {
				return err
			}

			tourdays = append(tourdays, tourDay)
		}

		// Validate guide before modifying anything.
		guide, err := guideRepo.GetGuideByID(guideID)
		if err != nil {
			return err
		}

		for _, tourday := range tourdays {

			if err := tourDayRepo.UnassignGuide(tourday.ID); err != nil {
				return err
			}

			if err := s.setAssignmentStatusHistoryRegistry(
				tourday.ID,
				guide.ID,
				"UNASSIGNED",
				"GUIDE",
				tourDayRepo,
				guideRepo,
				assignmentRepo,
			); err != nil {
				return err
			}
		}

		return nil
	})
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

	expiredTours, err := s.tourDayRepo.GetExpiredGuideAssignments(ctx)
	if err != nil {
		return result, err
	}

	for _, tour := range expiredTours {
		result.Found++

		err := db.WithTransaction(s.db, func(tx *sql.Tx) error {

			tourDayRepo := s.tourDayRepo.WithTx(tx)
			historyRepo := s.repo.WithTx(tx)

			if err := tourDayRepo.UpdateTourDayStatus(
				tour.ID,
				domain.RESERVATION_PAYMENT_PENDING,
			); err != nil {
				return err
			}

			history := TourDayStatusHistory{
				ID:             uuid.New().String(),
				TourDayID:      tour.ID,
				PreviousStatus: string(tour.Status),
				NewStatus:      string(domain.RESERVATION_PAYMENT_PENDING),
				ChangedBy:      "SYSTEM",
				Reason:         "Guide assignment expired",
				CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
			}

			if err := historyRepo.AddTourStatus(history); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			result.Failed++
			continue
		}

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
	return db.WithTransaction(s.db, func(tx *sql.Tx) error {

		tourDayRepo := s.tourDayRepo.WithTx(tx)
		historyRepo := s.repo.WithTx(tx)

		tourDay, err := tourDayRepo.GetTourDayByID(tourDayID)
		if err != nil {
			return err
		}

		if tourDay.Status != domain.RESERVATION_GUIDE_PREASSIGNED {
			return errors.ErrInvalidInput
		}

		if err := tourDayRepo.UpdateTourDayStatus(
			tourDayID,
			domain.RESERVATION_GUIDE_CONFIRMED,
		); err != nil {
			return err
		}

		history := TourDayStatusHistory{
			ID:             uuid.New().String(),
			TourDayID:      tourDay.ID,
			PreviousStatus: string(tourDay.Status),
			NewStatus:      string(domain.RESERVATION_GUIDE_CONFIRMED),
			ChangedBy:      "GUIDE",
			Reason:         "Guide confirmed tour",
			CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
		}

		if err := historyRepo.AddTourStatus(history); err != nil {
			return err
		}

		return nil
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

func (s *Service) setAssignmentStatusHistoryRegistry(
	tourDayID string,
	guideID string,
	action string,
	changedBy string,
	tourDayRepo *tourdays.Repository,
	guideRepo *guides.Repository,
	assignmentRepo *Repository) error {

	tourDay, err := tourDayRepo.GetTourDayByID(tourDayID)
	if err != nil {
		return err
	}

	guide, err := guideRepo.GetGuideByID(guideID)
	if err != nil {
		return err
	}

	assignment := TourDayAssigmentHistory{
		ID:        uuid.New().String(),
		TourDayID: tourDay.ID,
		GuideID:   guide.ID,
		Action:    action,
		ChangedBy: changedBy,
		CreatedAt: tourDay.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return assignmentRepo.AddAssignmentStatus(assignment)
}

func (s *Service) SetAssignmentStatusHistoryRegistry(
	tourDayID string,
	guideID string,
	action string,
	changedBy string,
) error {
	return s.setAssignmentStatusHistoryRegistry(
		tourDayID,
		guideID,
		action,
		changedBy,
		s.tourDayRepo,
		s.guideRepo,
		s.repo,
	)
}
