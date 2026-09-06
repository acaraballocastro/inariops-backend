package activity

import (
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetActivityByID(activityID string) (Activity, error) {
	activity, err := s.repo.GetActivityByID(activityID)
	if err != nil {
		return Activity{}, err
	}
	return activity, nil
}

func (s *Service) GetAllActivities() ([]Activity, error) {
	activities, err := s.repo.GetAllActivities()
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (s *Service) CreateActivity(activity *CreateActivityRequest) error {
	newActivity := &Activity{
		ID:          uuid.New().String(),
		Name:        activity.Name,
		Description: activity.Description,
		IsActive:    true,
	}

	return s.repo.CreateActivity(*newActivity)
}

func (s *Service) UpdateActivity(activity *UpdateActivityRequest) error {
	existingActivity, err := s.repo.GetActivityByID(activity.ID)
	if err != nil {
		return err
	}
	existingActivity.Name = *activity.Name
	existingActivity.Description = activity.Description
	existingActivity.IsActive = *activity.IsActive

	return s.repo.UpdateActivity(existingActivity)
}

func (s *Service) DeleteActivity(activityID string) error {
	return s.repo.DeleteActivity(activityID)
}
