package users

import (
	"inariops/internal/domain"
	"inariops/internal/modules/auth"
	"inariops/internal/modules/guides"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo   *Repository
	auth   *auth.Service
	guides *guides.Service
}

type UpdateUserInput struct {
	ID    string
	Name  *string
	Email *string
	Phone *string
	Role  *domain.UserRole
}

func NewService(repo *Repository, auth *auth.Service, guides *guides.Service) *Service {
	return &Service{repo: repo, auth: auth, guides: guides}
}

func (s *Service) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAllUsers()
}

func (s *Service) CreateUser(name, email, phone string, role string) (domain.UserCredentials, error) {
	user := domain.User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Role:      domain.UserRole(role),
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	validateCreateUserErr := s.validateCreateUser(name, email, role)
	if validateCreateUserErr != nil {
		return domain.UserCredentials{}, validateCreateUserErr
	}

	err := s.repo.CreateUser(user)
	if err != nil {
		return domain.UserCredentials{}, err
	}

	authCredentials, err := s.auth.CreateCredentials(user.ID)
	if err != nil {
		_ = s.repo.DeleteUser(user.ID)
		return domain.UserCredentials{}, err
	}

	guide := domain.Guide{
		ID:             uuid.New().String(),
		UserID:         user.ID,
		MaxToursPerDay: 1, // Set a default value
		CreatedAt:      time.Now(),
	}

	err = s.guides.CreateGuide(guide)
	if err != nil {
		_ = s.repo.DeleteUser(user.ID)
		return domain.UserCredentials{}, err
	}

	return domain.UserCredentials{
		User:            user,
		AuthCredentials: authCredentials,
	}, nil
}

func (s *Service) GetUserByID(id string) (*domain.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *Service) UpdateUser(input UpdateUserInput) (domain.User, error) {
	validateUpdateUserErr := s.validateUpdateUser(input)
	if validateUpdateUserErr != nil {
		return domain.User{}, validateUpdateUserErr
	}

	user, err := s.repo.GetUserByID(input.ID)
	if err != nil {
		return domain.User{}, err
	}

	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Email != nil {
		user.Email = *input.Email
	}

	if input.Phone != nil {
		user.Phone = *input.Phone
	}

	if input.Role != nil {
		user.Role = *input.Role
	}

	err = s.repo.UpdateUser(*user)
	if err != nil {
		return domain.User{}, err
	}

	return *user, nil
}

func (s *Service) DeactivateUser(id string) error {
	return s.repo.DeactivateUser(id)
}

func (s *Service) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}
