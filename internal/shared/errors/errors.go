package errors

var (
	// =========================================================
	// GENERIC
	// =========================================================

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

	ErrConflict = AppError{
		Code:    "CONFLICT",
		Message: "resource conflict",
	}

	// =========================================================
	// AUTH
	// =========================================================

	ErrInvalidCredentials = AppError{
		Code:    "INVALID_CREDENTIALS",
		Message: "invalid credentials",
	}

	ErrInvalidCurrentPassword = AppError{
		Code:    "INVALID_CURRENT_PASSWORD",
		Message: "current password is incorrect",
	}

	ErrWeakPassword = AppError{
		Code:    "WEAK_PASSWORD",
		Message: "password does not meet security requirements, must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character",
	}

	// =========================================================
	// USERS
	// =========================================================

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

	ErrInvalidUserID = AppError{
		Code:    "INVALID_USER_ID",
		Message: "invalid user id",
	}

	ErrUserInactive = AppError{
		Code:    "USER_INACTIVE",
		Message: "user is inactive",
	}

	// =========================================================
	// CUSTOMERS
	// =========================================================

	ErrCustomerNotFound = AppError{
		Code:    "CUSTOMER_NOT_FOUND",
		Message: "customer not found",
	}

	ErrCustomerAlreadyExists = AppError{
		Code:    "CUSTOMER_ALREADY_EXISTS",
		Message: "customer already exists",
	}

	ErrCustomerAlreadyInReservation = AppError{
		Code:    "CUSTOMER_ALREADY_IN_RESERVATION",
		Message: "customer is already in the reservation",
	}

	ErrCustomerNotInReservation = AppError{
		Code:    "CUSTOMER_NOT_IN_RESERVATION",
		Message: "customer is not in the reservation",
	}

	// =========================================================
	// GUIDES
	// =========================================================

	ErrGuideNotFound = AppError{
		Code:    "GUIDE_NOT_FOUND",
		Message: "guide not found",
	}

	ErrGuideAlreadyExists = AppError{
		Code:    "GUIDE_ALREADY_EXISTS",
		Message: "guide already exists",
	}

	ErrGuideAlreadyAssigned = AppError{
		Code:    "GUIDE_ALREADY_ASSIGNED",
		Message: "guide already assigned to this tour day",
	}

	ErrGuideNotAssigned = AppError{
		Code:    "GUIDE_NOT_ASSIGNED",
		Message: "guide is not assigned to this tour day",
	}

	// =========================================================
	// LANGUAGES
	// =========================================================

	ErrLanguageNotFound = AppError{
		Code:    "LANGUAGE_NOT_FOUND",
		Message: "language not found",
	}

	ErrInvalidLanguage = AppError{
		Code:    "INVALID_LANGUAGE",
		Message: "invalid language data",
	}

	ErrLanguageCodeRequired = AppError{
		Code:    "LANGUAGE_CODE_REQUIRED",
		Message: "language code is required",
	}

	ErrLanguageAlreadyAssigned = AppError{
		Code:    "LANGUAGE_ALREADY_ASSIGNED",
		Message: "language is already assigned to the guide",
	}

	ErrLanguageNotAssigned = AppError{
		Code:    "LANGUAGE_NOT_ASSIGNED",
		Message: "language is not assigned to the guide",
	}

	// =========================================================
	// ZONES
	// =========================================================

	ErrZoneNotFound = AppError{
		Code:    "ZONE_NOT_FOUND",
		Message: "zone not found",
	}

	ErrInvalidZone = AppError{
		Code:    "INVALID_ZONE",
		Message: "invalid zone data",
	}

	ErrZoneAlreadyAssigned = AppError{
		Code:    "ZONE_ALREADY_ASSIGNED",
		Message: "zone is already assigned to the guide",
	}

	ErrZoneNotAssigned = AppError{
		Code:    "ZONE_NOT_ASSIGNED",
		Message: "zone is not assigned to the guide",
	}

	// =========================================================
	// AVAILABILITY
	// =========================================================

	ErrInvalidAvailabilityDate = AppError{
		Code:    "INVALID_AVAILABILITY_DATE",
		Message: "invalid availability date",
	}

	ErrAvailabilityConflict = AppError{
		Code:    "AVAILABILITY_CONFLICT",
		Message: "availability overlaps an existing range",
	}

	ErrAvailabilityNotFound = AppError{
		Code:    "AVAILABILITY_NOT_FOUND",
		Message: "availability not found",
	}

	// =========================================================
	// AGENCIES
	// =========================================================

	ErrAgencyNotFound = AppError{
		Code:    "AGENCY_NOT_FOUND",
		Message: "agency not found",
	}

	ErrExistingAgency = AppError{
		Code:    "EXISTING_AGENCY",
		Message: "agency already exists",
	}

	ErrInvalidAgency = AppError{
		Code:    "INVALID_AGENCY",
		Message: "invalid agency data",
	}

	// =========================================================
	// RESERVATIONS
	// =========================================================

	ErrReservationNotFound = AppError{
		Code:    "RESERVATION_NOT_FOUND",
		Message: "reservation not found",
	}

	ErrFailedToCreateReservation = AppError{
		Code:    "FAILED_TO_CREATE_RESERVATION",
		Message: "failed to create reservation",
	}

	ErrFailedToUpdateReservation = AppError{
		Code:    "FAILED_TO_UPDATE_RESERVATION",
		Message: "failed to update reservation",
	}

	ErrInvalidReservationCode = AppError{
		Code:    "INVALID_RESERVATION_CODE",
		Message: "invalid reservation code",
	}

	ErrInvalidReservationDates = AppError{
		Code:    "INVALID_RESERVATION_DATES",
		Message: "invalid reservation dates",
	}

	ErrReservationCancelled = AppError{
		Code:    "RESERVATION_CANCELLED",
		Message: "reservation is cancelled",
	}

	ErrInvalidReservationStatus = AppError{
		Code:    "INVALID_RESERVATION_STATUS",
		Message: "invalid reservation status",
	}

	ErrInvalidSignatureStatus = AppError{
		Code:    "INVALID_SIGNATURE_STATUS",
		Message: "invalid signature status",
	}

	ErrInvalidVoucherStatus = AppError{
		Code:    "INVALID_VOUCHER_STATUS",
		Message: "invalid voucher status",
	}

	// =========================================================
	// TOUR DAYS
	// =========================================================

	ErrTourDayNotFound = AppError{
		Code:    "TOUR_DAY_NOT_FOUND",
		Message: "tour day not found",
	}

	ErrTourDayCancelled = AppError{
		Code:    "TOUR_DAY_CANCELLED",
		Message: "tour day cancelled",
	}

	ErrTourDayAlreadyConfirmed = AppError{
		Code:    "TOUR_DAY_ALREADY_CONFIRMED",
		Message: "tour day is already confirmed",
	}

	ErrInvalidTourDayStatus = AppError{
		Code:    "INVALID_TOUR_DAY_STATUS",
		Message: "invalid tour day status",
	}

	ErrGuideAssignmentNotAllowed = AppError{
		Code:    "GUIDE_ASSIGNMENT_NOT_ALLOWED",
		Message: "guide cannot be assigned to this tour day",
	}

	// =========================================================
	// ITINERARIES
	// =========================================================

	ErrItineraryNotFound = AppError{
		Code:    "ITINERARY_NOT_FOUND",
		Message: "itinerary not found",
	}

	ErrActivityNotFound = AppError{
		Code:    "ACTIVITY_NOT_FOUND",
		Message: "activity not found",
	}

	ErrActivityAlreadyExists = AppError{
		Code:    "ACTIVITY_ALREADY_EXISTS",
		Message: "activity already exists",
	}

	ErrInvalidActivity = AppError{
		Code:    "INVALID_ACTIVITY",
		Message: "invalid activity data",
	}

	ErrActivityAlreadyInItinerary = AppError{
		Code:    "ACTIVITY_ALREADY_IN_ITINERARY",
		Message: "activity is already in the itinerary",
	}

	ErrActivityNotInItinerary = AppError{
		Code:    "ACTIVITY_NOT_IN_ITINERARY",
		Message: "activity is not in the itinerary",
	}

	ErrPlaceNotFound = AppError{
		Code:    "PLACE_NOT_FOUND",
		Message: "place not found",
	}

	ErrPlaceAlreadyExists = AppError{
		Code:    "PLACE_ALREADY_EXISTS",
		Message: "place already exists",
	}

	ErrInvalidPlace = AppError{
		Code:    "INVALID_PLACE",
		Message: "invalid place data",
	}

	ErrPlaceAlreadyInItinerary = AppError{
		Code:    "PLACE_ALREADY_IN_ITINERARY",
		Message: "place is already in the itinerary",
	}

	ErrPlaceNotInItinerary = AppError{
		Code:    "PLACE_NOT_IN_ITINERARY",
		Message: "place is not in the itinerary",
	}

	// =========================================================
	// TOUR DAY ITINERARY
	// =========================================================

	ErrTourDayItineraryNotFound = AppError{
		Code:    "TOUR_DAY_ITINERARY_NOT_FOUND",
		Message: "tour day itinerary not found",
	}

	ErrItineraryAlreadyAssigned = AppError{
		Code:    "ITINERARY_ALREADY_ASSIGNED",
		Message: "itinerary is already assigned to this tour day",
	}

	ErrItineraryNotAssigned = AppError{
		Code:    "ITINERARY_NOT_ASSIGNED",
		Message: "itinerary is not assigned to this tour day",
	}
)

type AppError struct {
	Code    string
	Message string
}

func (e AppError) Error() string {
	return e.Message
}
