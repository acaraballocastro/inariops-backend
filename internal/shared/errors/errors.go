package errors

var (
	// generic
	ErrInvalidInput = AppError{
		Code:    "INVALID_INPUT",
		Message: "invalid input",
	}

	ErrUnauthorized = AppError{
		Code:    "UNAUTHORIZED",
		Message: "unauthorized",
	}

	ErrForbidden = AppError{
		Code:    "FORBIDDEN",
		Message: "forbidden",
	}

	ErrNotFound = AppError{
		Code:    "NOT_FOUND",
		Message: "resource not found",
	}

	// users
	ErrUserNotFound = AppError{
		Code:    "USER_NOT_FOUND",
		Message: "user not found",
	}

	ErrEmailAlreadyExists = AppError{
		Code:    "EMAIL_ALREADY_EXISTS",
		Message: "email already exists",
	}

	ErrNullableFieldEmpty = AppError{
		Code:    "NULLABLE_FIELD_EMPTY",
		Message: "nullable fields cannot be empty if provided",
	}

	ErrInvalidRole = AppError{
		Code:    "INVALID_ROLE",
		Message: "invalid role",
	}
)

type AppError struct {
	Code    string
	Message string
}

func (e AppError) Error() string {
	return e.Message
}
