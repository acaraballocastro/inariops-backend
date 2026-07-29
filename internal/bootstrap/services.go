package bootstrap

import (
	reservationsapp "inariops/internal/application/reservations"
	tourdaysapp "inariops/internal/application/tour_days"
)

type Services struct {
	TourDay     *tourdaysapp.Service
	Reservation *reservationsapp.Service
}

func newServices(repositories *Repositories) *Services {

	return &Services{

		TourDay: tourdaysapp.NewService(
			repositories.TourDay,
			repositories.Guide,
			repositories.TourDayApplication,
		),

		Reservation: reservationsapp.NewService(
			repositories.Reservation,
			repositories.Customer,
			repositories.ReservationCustomer,
			repositories.TourDay,
			repositories.Agency,
		),
	}
}
