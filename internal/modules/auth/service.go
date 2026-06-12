package auth

import (
	"inariops/internal/shared/errors"
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

	cred, err := s.repo.GetCredentials(user.ID)
	if err != nil {
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

	cred, err := s.repo.GetCredentials(userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	if !CheckPassword(cred.PasswordHash, oldPass) {
		return errors.ErrUnauthorized
	}

	hash, err := HashPassword(newPass)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(userID, hash)
}
