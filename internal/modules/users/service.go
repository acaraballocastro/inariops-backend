package users

import (
	"inariops/internal/shared/errors"
	"regexp"
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

	validateCreateUserErr := s.validateCreateUser(name, email, role)
	if validateCreateUserErr != nil {
		return User{}, validateCreateUserErr
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
	validateUpdateUserErr := s.validateUpdateUser(input)
	if validateUpdateUserErr != nil {
		return User{}, validateUpdateUserErr
	}

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

// VALIDATION
func (s *Service) validateCreateUser(name, email, role string) error {
	if name == "" {
		return errors.ErrNullableFieldEmpty
	}

	if email == "" {
		return errors.ErrNullableFieldEmpty
	}

	if !isValidEmail(email) {
		return errors.ErrInvalidInput
	}

	if role != "ADMIN" && role != "GUIDE" {
		return errors.ErrInvalidRole
	}

	return nil
}

func isValidEmail(email string) bool {
	// Simple regex for email validation
	const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func (s *Service) validateUpdateUser(input UpdateUserInput) error {

	if input.Email != nil && *input.Email == "" {
		return errors.ErrNullableFieldEmpty
	}

	if input.Name != nil && *input.Name == "" {
		return errors.ErrNullableFieldEmpty
	}

	if input.Role != nil {
		r := *input.Role
		if r != "ADMIN" && r != "GUIDE" {
			return errors.ErrInvalidRole
		}
	}

	return nil
}
