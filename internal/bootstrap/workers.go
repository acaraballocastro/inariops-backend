package bootstrap

import (
	"context"

	"inariops/internal/workers"
)

func (a *App) StartWorkers(
	ctx context.Context,
) {

	scheduler := workers.NewScheduler()

	scheduler.Register(

		workers.NewGuideAssignmentTimeoutWorker(
			a.Services.TourDayApplication,
		),

		workers.NewReservationStatusWorker(
			a.Services.ReservationApplication,
			a.Repositories.Reservation,
		),
	)

	scheduler.Start(ctx)
}
