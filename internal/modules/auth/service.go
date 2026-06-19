package auth

import (
	"inariops/internal/domain"
	"inariops/internal/shared/errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Login(email, password string) (LoginResponse, error) {

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return LoginResponse{}, errors.ErrUserNotFound
	}
	if !user.IsActive {
		return LoginResponse{}, errors.ErrUnauthorized
	}

	cred, err := s.repo.GetCredentials(user.ID)
	if err != nil {
		return LoginResponse{}, errors.ErrUnauthorized
	}
	if !cred.IsActive {
		return LoginResponse{}, errors.ErrUnauthorized
	}

	if !CheckPassword(cred.PasswordHash, password) {
		return LoginResponse{}, errors.ErrUnauthorized
	}

	token, err := GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		Token:              token,
		UserID:             user.ID,
		Role:               string(user.Role),
		MustChangePassword: cred.MustChangePassword,
	}, nil
}

func (s *Service) ChangePassword(userID, oldPass, newPass string) error {

	log.Printf("ChangePassword: starting for user %s", userID)

	cred, err := s.repo.GetCredentials(userID)
	if err != nil {
		log.Printf("ChangePassword: failed to get credentials for user %s: %v", userID, err)
		return errors.ErrUserNotFound
	}

	if !CheckPassword(cred.PasswordHash, oldPass) {
		log.Printf("ChangePassword: old password does not match for user %s", userID)
		return errors.ErrUnauthorized
	}

	hash, err := HashPassword(newPass)
	if err != nil {
		log.Printf("ChangePassword: failed to hash new password for user %s: %v", userID, err)
		return err
	}

	if err := s.repo.UpdatePassword(userID, hash); err != nil {
		log.Printf("ChangePassword: failed to update password for user %s: %v", userID, err)
		return err
	}

	log.Printf("ChangePassword: success for user %s", userID)
	return nil
}

func (s *Service) CreateCredentials(userID string) (domain.AuthCredentials, error) {
	password, err := createDefaultPassword()
	if err != nil {
		return domain.AuthCredentials{}, err
	}

	log.Printf("Default password for user %s: %s", userID, password)

	passwordHash, err := HashPassword(password)
	if err != nil {
		return domain.AuthCredentials{}, err
	}

	credentials := Credentials{
		ID:                 uuid.New().String(),
		UserID:             userID,
		PasswordHash:       passwordHash,
		MustChangePassword: true,
		IsActive:           true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	return domain.AuthCredentials{
		UserID:             credentials.UserID,
		MustChangePassword: credentials.MustChangePassword,
		IsActive:           credentials.IsActive,
	}, s.repo.CreateCredentials(credentials)
}
