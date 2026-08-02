package bootstrap

import (
	"database/sql"

	tourdaysapp "inariops/internal/application/tour_days"

	"inariops/internal/modules/auth"
	"inariops/internal/modules/customers"
	"inariops/internal/modules/guides"
	languageguides "inariops/internal/modules/guides/language_guides"
	"inariops/internal/modules/guides/languages"
	"inariops/internal/modules/tours/agencies"
	"inariops/internal/modules/tours/reservations"
	reservationscustomers "inariops/internal/modules/tours/reservations_customers"
	tourdays "inariops/internal/modules/tours/tour_days"
	"inariops/internal/modules/users"
)

type Repositories struct {
	Auth *auth.Repository

	User *users.Repository

	Guide *guides.Repository

	Reservation *reservations.Repository

	ReservationCustomer *reservationscustomers.Repository

	TourDay *tourdays.Repository

	TourDayApplication *tourdaysapp.Repository

	Customer *customers.Repository

	Agency *agencies.Repository

	Language *languages.Repository

	LanguageGuide *languageguides.Repository
}

func newRepositories(db *sql.DB) *Repositories {

	return &Repositories{

		Auth: auth.NewRepository(db),

		User: users.NewRepository(db),

		Guide: guides.NewRepository(db),

		Reservation: reservations.NewRepository(db),

		ReservationCustomer: reservationscustomers.NewRepository(db),

		TourDay: tourdays.NewRepository(db),

		TourDayApplication: tourdaysapp.NewRepository(db),

		Customer: customers.NewRepository(db),

		Agency: agencies.NewRepository(db),

		Language: languages.NewRepository(db),

		LanguageGuide: languageguides.NewRepository(db),
	}
}
