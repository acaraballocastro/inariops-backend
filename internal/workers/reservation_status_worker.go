package workers

import (
	"context"
	"time"

	"inariops/internal/domain"
	reservationsapp "inariops/internal/modules/tours/application"
	reservations "inariops/internal/modules/tours/reservations"
)

type ReservationStatusWorker struct {
	service          *reservationsapp.Service
	reservationsRepo *reservations.Repository
}

func NewReservationStatusWorker(
	service *reservationsapp.Service,
	reservationsRepo *reservations.Repository,
) *ReservationStatusWorker {

	return &ReservationStatusWorker{
		service:          service,
		reservationsRepo: reservationsRepo,
	}
}

func (w *ReservationStatusWorker) Name() string {
	return "reservation-status-sync"
}

func (w *ReservationStatusWorker) Interval() time.Duration {
	return 5 * time.Minute
}

func (w *ReservationStatusWorker) Run(
	ctx context.Context,
) (domain.Result, error) {

	reservations, err := w.reservationsRepo.GetAllReservations()

	if err != nil {
		return domain.Result{}, err
	}

	result := domain.Result{
		Found: len(reservations),
	}

	for _, reservation := range reservations {

		err := w.service.SyncReservationStatus(
			reservation.ID,
		)

		if err != nil {
			result.Failed++
			continue
		}

		result.Processed++
	}

	return result, nil
}
