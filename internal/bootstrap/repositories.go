package bootstrap

import (
	"database/sql"

	"inariops/internal/modules/customers"
	"inariops/internal/modules/guides"
	"inariops/internal/modules/tours/agencies"
	"inariops/internal/modules/tours/reservations"
	reservationscustomers "inariops/internal/modules/tours/reservations_customers"
	tourdays "inariops/internal/modules/tours/tour_days"

	tourdaysapp "inariops/internal/application/tour_days"
)

type Repositories struct {
	Guide               *guides.Repository
	TourDay             *tourdays.Repository
	TourDayApplication  *tourdaysapp.Repository
	Reservation         *reservations.Repository
	Customer            *customers.Repository
	ReservationCustomer *reservationscustomers.Repository
	Agency              *agencies.Repository
}

func newRepositories(db *sql.DB) *Repositories {

	return &Repositories{

		Guide: guides.NewRepository(db),

		TourDay: tourdays.NewRepository(db),

		TourDayApplication: tourdaysapp.NewRepository(db),

		Reservation: reservations.NewRepository(db),

		Customer: customers.NewRepository(db),

		ReservationCustomer: reservationscustomers.NewRepository(db),

		Agency: agencies.NewRepository(db),
	}
}
