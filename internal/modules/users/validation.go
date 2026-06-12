package users

import (
	"inariops/internal/shared/errors"
	"regexp"
)

// VALIDATION
func (s *Service) validateCreateUser(name, email, role string) error {
	if name == "" {
		return errors.ErrInvalidInput
	}

	if email == "" {
		return errors.ErrInvalidInput
	}

	if !isValidEmail(email) {
		return errors.ErrInvalidInput
	}

	r := UserRole(role)
	if r != RoleAdmin && r != RoleGuide {
		return errors.ErrInvalidRole
	}

	return nil
}

var emailRE = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRE.MatchString(email)
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
