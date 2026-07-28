package tourdaysapp

import (
	"context"
	"inariops/internal/domain"
	"inariops/internal/modules/guides"
	tours "inariops/internal/modules/tours/shared"
	tourdays "inariops/internal/modules/tours/tour_days"

	"github.com/google/uuid"
)

type Service struct {
	repo        *Repository
	tourDayRepo *tourdays.Repository
	guideRepo   *guides.Repository
}

func NewService(tourDayRepo *tourdays.Repository, guideRepo *guides.Repository, assigmentRepo *Repository) *Service {
	return &Service{
		tourDayRepo: tourDayRepo,
		guideRepo:   guideRepo,
		repo:        assigmentRepo,
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

		err := s.tourDayRepo.UpdateTourDayStatus(tour.ID, domain.RESERVATION_GUIDE_CONFIRMED)
		if err != nil {
			result.Failed++
			continue
		}

		err = s.SetTourDayHistory(tour.ID, string(domain.RESERVATION_GUIDE_CONFIRMED), "SYSTEM", "Guide assignment expired")

		result.Processed++

	}

	return result, nil
}
