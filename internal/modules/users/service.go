package users

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

type UpdateUserInput struct {
	ID    string
	Name  *string
	Email *string
	Phone *string
	Role  *UserRole
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllUsers() ([]User, error) {
	return s.repo.GetAllUsers()
}

func (s *Service) CreateUser(name, email, phone string, role string) (User, error) {
	user := User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Role:      UserRole(role),
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	err := s.repo.CreateUser(user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Service) GetUserByID(id string) (*User, error) {
	return s.repo.GetUserByID(id)
}

func (s *Service) UpdateUser(input UpdateUserInput) (User, error) {

	user, err := s.repo.GetUserByID(input.ID)
	if err != nil {
		return User{}, err
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
		return User{}, err
	}

	return *user, nil
}

func (s *Service) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}
