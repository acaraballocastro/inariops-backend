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

	// reservations
	ErrReservationNotFound = AppError{
		Code:    "RESERVATION_NOT_FOUND",
		Message: "reservation not found",
	}

	ErrInvalidReservationCode = AppError{
		Code:    "INVALID_RESERVATION_CODE",
		Message: "invalid reservation code",
	}

	ErrWeakPassword = AppError{
		Code:    "WEAK_PASSWORD",
		Message: "password does not meet security requirements, must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character",
	}

	ErrFailedToUpdateReservation = AppError{
		Code:    "FAILED_TO_UPDATE_RESERVATION",
		Message: "failed to update reservation",
	}

	// tour days
	ErrTourDayNotFound = AppError{
		Code:    "TOUR_DAY_NOT_FOUND",
		Message: "tour day not found",
	}

	ErrTourDayCancelled = AppError{
		Code:    "TOUR_DAY_CANCELLED",
		Message: "tour day cancelled",
	}

	// guides
	ErrGuideNotFound = AppError{
		Code:    "GUIDE_NOT_FOUND",
		Message: "guide not found",
	}

	ErrGuideAlreadyAssigned = AppError{
		Code:    "GUIDE_ALREADY_ASSIGNED",
		Message: "guide already assigned to this tour day",
	}

	ErrGuideNotAssigned = AppError{
		Code:    "GUIDE_NOT_ASSIGNED",
		Message: "guide is not assigned to this tour day",
	}

	ErrCustomerNotFound = AppError{
		Code:    "CUSTOMER_NOT_FOUND",
		Message: "customer not found",
	}

	ErrExistingAgency = AppError{
		Code:    "EXISTING_AGENCY",
		Message: "agency already exists",
	}
)

type AppError struct {
	Code    string
	Message string
}

func (e AppError) Error() string {
	return e.Message
}
