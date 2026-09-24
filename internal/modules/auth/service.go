package auth

import (
	"inariops/internal/domain"
	"inariops/internal/shared/errors"
	"inariops/internal/shared/logger"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo       *Repository
	jwtManager *JWTManager
}

func NewService(repo *Repository, jwtManager *JWTManager) *Service {
	return &Service{repo: repo, jwtManager: jwtManager}
}

func (s *Service) Login(email, password string) (LoginResponse, error) {

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		logger.Error("Login: failed to get user for email %s: %v", email, err)
		return LoginResponse{}, errors.ErrUserNotFound
	}
	if !user.IsActive {
		logger.Error("Login: user %s is not active", user.ID)
		return LoginResponse{}, errors.ErrUnauthorized
	}

	cred, err := s.repo.GetCredentials(user.ID)
	if err != nil {
		logger.Error("Login: failed to get credentials for user %s: %v", user.ID, err)
		return LoginResponse{}, errors.ErrUnauthorized
	}
	if !cred.IsActive {
		logger.Error("Login: credentials for user %s are not active", user.ID)
		return LoginResponse{}, errors.ErrUnauthorized
	}

	if !CheckPassword(cred.PasswordHash, password) {
		logger.Error("Login: failed for email %s", user.Email)
		return LoginResponse{}, errors.ErrUnauthorized
	}

	token, err := s.jwtManager.GenerateToken(user.ID, string(user.Role))
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

// TODO: Add a way to get the userID from the token.
func (s *Service) ChangePassword(userID, oldPass, newPass string) error {

	cred, err := s.repo.GetCredentials(userID)
	if err != nil {
		logger.Error("ChangePassword: failed to get credentials for user %s: %v", userID, err)
		return errors.ErrUserNotFound
	}

	if !CheckPassword(cred.PasswordHash, oldPass) {
		logger.Error("ChangePassword: old password does not match for user %s", userID)
		return errors.ErrUnauthorized
	}

	if !CheckNewPassword(newPass) {
		logger.Error("ChangePassword: new password does not meet security requirements for user %s", userID)
		return errors.ErrWeakPassword
	}

	hash, err := HashPassword(newPass)
	if err != nil {
		logger.Error("ChangePassword: failed to hash new password for user %s: %v", userID, err)
		return err
	}

	if err := s.repo.UpdatePassword(userID, hash); err != nil {
		logger.Error("ChangePassword: failed to update password for user %s: %v", userID, err)
		return err
	}

	return nil
}

// TODO: Add a function to reset password and send email with new password
func (s *Service) CreateCredentials(userID string) (domain.AuthCredentials, error) {
	password, err := createDefaultPassword()
	if err != nil {
		logger.Error("CreateCredentials: failed to create default password for user %s: %v", userID, err)
		return domain.AuthCredentials{}, err
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		logger.Error("CreateCredentials: failed to hash password for user %s: %v", userID, err)
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

//TODO: Add a function to reset password and send email with new password
