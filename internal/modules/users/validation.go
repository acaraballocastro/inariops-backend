package users

import (
	"inariops/internal/shared/errors"
	"regexp"
)

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
