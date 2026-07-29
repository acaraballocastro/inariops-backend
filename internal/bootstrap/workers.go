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
			a.Services.TourDay,
		),

		workers.NewReservationStatusWorker(
			a.Services.Reservation,
			a.Repositories.Reservation,
		),
	)

	scheduler.Start(ctx)
}
