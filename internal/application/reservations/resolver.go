package reservationsapp

import (
	"inariops/internal/domain"
	tours "inariops/internal/modules/tours/shared"
)

type Resolver struct{}

func NewResolver() *Resolver {
	return &Resolver{}
}

func (r *Resolver) Resolve(
	tourDays []tours.TourDay,
) domain.ReservationStatus {

	if len(tourDays) == 0 {
		return domain.RESERVATION_PENDING_ASSIGNMENT
	}

	currentStatus := tourDays[0].Status

	for _, tourDay := range tourDays {

		if tourDay.Status != currentStatus {
			return domain.RESERVATION_PENDING
		}
	}

	switch currentStatus {

	case domain.RESERVATION_PENDING_ASSIGNMENT:
		return domain.RESERVATION_PENDING_ASSIGNMENT

	case domain.RESERVATION_GUIDE_PREASSIGNED:
		return domain.RESERVATION_GUIDE_PREASSIGNED

	case domain.RESERVATION_PAYMENT_PENDING:
		return domain.RESERVATION_PAYMENT_PENDING

	case domain.RESERVATION_PARTIALLY_CONFIRMED:
		return domain.RESERVATION_PARTIALLY_CONFIRMED

	case domain.RESERVATION_CONFIRMED:
		return domain.RESERVATION_CONFIRMED

	case domain.RESERVATION_COMPLETED:
		return domain.RESERVATION_COMPLETED

	case domain.RESERVATION_CANCELLED:
		return domain.RESERVATION_CANCELLED

	case domain.RESERVATION_FORCE_MAJEURE_CANCELLED:
		return domain.RESERVATION_FORCE_MAJEURE_CANCELLED

	default:
		return domain.RESERVATION_PENDING
	}
}
