package tourdaysapp

import (
	"inariops/internal/modules/guides"
	tours "inariops/internal/modules/tours/shared"
	tourdays "inariops/internal/modules/tours/tour_days"
)

type Service struct {
	tourDayRepo *tourdays.Repository
	guideRepo   *guides.Repository
}

func NewService(tourDayRepo *tourdays.Repository, guideRepo *guides.Repository) *Service {
	return &Service{
		tourDayRepo: tourDayRepo,
		guideRepo:   guideRepo,
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

		return nil
	}

	err = s.tourDayRepo.AssignGuideToMultipleTourDays(tourDayIDs, guide.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UnassignGuide(tourDayIDs []string) error {
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

		return nil
	}

	for _, tourday := range tourdays {
		err := s.tourDayRepo.UnassignGuide(tourday.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
