package workers

import (
	"context"
	"time"

	reservationsapp "inariops/internal/application/reservations"
	"inariops/internal/domain"
	reservations "inariops/internal/modules/tours/reservations"
	"inariops/internal/shared/logger"
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
	return 1 * time.Minute
}

func (w *ReservationStatusWorker) Run(
	ctx context.Context,
) (domain.Result, error) {

	reservations, err := w.reservationsRepo.GetAllReservations()
	logger.Info("reservation-status-sync", "found reservations: %d", len(reservations))

	if err != nil {
		return domain.Result{}, err
	}

	result := domain.Result{
		Found: len(reservations),
	}

	for _, reservation := range reservations {

		err := w.service.SyncReservationStatus(
			*reservation.Code,
		)

		logger.Info("reservation-status-sync", "syncing reservation: %s, error: %v", *reservation.Code, err)
		if err != nil {
			result.Failed++
			continue
		}

		result.Processed++
	}

	return result, nil
}
