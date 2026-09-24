package users

import (
	"database/sql"
	"inariops/internal/db"
	"inariops/internal/domain"
	"inariops/internal/modules/auth"
	"inariops/internal/modules/guides"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	db     *sql.DB
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

func NewService(
	database *sql.DB,
	repo *Repository,
	authService *auth.Service,
	guidesService *guides.Service,
) *Service {
	return &Service{
		db:     database,
		repo:   repo,
		auth:   authService,
		guides: guidesService,
	}
}

func (s *Service) WithTx(tx *sql.Tx) *Service {
	return &Service{
		db:     s.db,
		repo:   s.repo.WithTx(tx),
		auth:   s.auth.WithTx(tx),
		guides: s.guides.WithTx(tx),
	}
}

func (s *Service) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAllUsers()
}

func (s *Service) GetAllGuides() ([]domain.GuideUser, error) {
	guides, err := s.guides.GetAllGuides()
	if err != nil {
		return nil, err
	}

	var guideUsers []domain.GuideUser
	for _, guide := range guides {
		user, err := s.repo.GetUserByID(guide.UserID)
		if err != nil {
			return nil, err
		}

		guideUser := domain.GuideUser{
			ID:             guide.ID,
			UserID:         guide.UserID,
			Name:           user.Name,
			Email:          user.Email,
			Phone:          user.Phone,
			MaxToursPerDay: guide.MaxToursPerDay,
			CreatedAt:      guide.CreatedAt,
		}
		guideUsers = append(guideUsers, guideUser)
	}

	return guideUsers, nil

}

func (s *Service) CreateUser(
	name,
	email,
	phone string,
	role string,
) (domain.UserCredentials, error) {

	user := domain.User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Role:      domain.UserRole(role),
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := s.validateCreateUser(name, email, role); err != nil {
		return domain.UserCredentials{}, err
	}

	var result domain.UserCredentials

	err := db.WithTransaction(s.db, func(tx *sql.Tx) error {

		txService := s.WithTx(tx)

		if err := txService.repo.CreateUser(user); err != nil {
			return err
		}

		authCredentials, err := txService.auth.CreateCredentials(user.ID)
		if err != nil {
			return err
		}

		if user.Role == domain.RoleGuide {
			guide := domain.Guide{
				ID:             uuid.New().String(),
				UserID:         user.ID,
				MaxToursPerDay: 1,
				CreatedAt:      time.Now(),
			}

			if err := txService.guides.CreateGuide(guide); err != nil {
				return err
			}
		}

		result = domain.UserCredentials{
			User:            user,
			AuthCredentials: authCredentials,
		}

		return nil
	})

	if err != nil {
		return domain.UserCredentials{}, err
	}

	return result, nil
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
