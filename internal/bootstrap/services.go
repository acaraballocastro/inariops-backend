package bootstrap

import (
	reservationcustomersapp "inariops/internal/application/reservation_customers"
	reservationsapp "inariops/internal/application/reservations"
	tourdaysapp "inariops/internal/application/tour_days"

	"inariops/internal/modules/auth"
	"inariops/internal/modules/customers"
	"inariops/internal/modules/guides"
	"inariops/internal/modules/guides/languages"
	"inariops/internal/modules/tours/agencies"
	"inariops/internal/modules/tours/reservations"
	tourdays "inariops/internal/modules/tours/tour_days"
	"inariops/internal/modules/users"
)

type Services struct {
	Auth *auth.Service

	User *users.Service

	Guide *guides.Service

	Reservation *reservations.Service

	TourDay *tourdays.Service

	Customer *customers.Service

	Agency *agencies.Service

	Language *languages.Service

	ReservationApplication *reservationsapp.Service

	ReservationCustomerApplication *reservationcustomersapp.Service

	TourDayApplication *tourdaysapp.Service
}

func newServices(
	repositories *Repositories,
) *Services {

	authService := auth.NewService(
		repositories.Auth,
	)

	guideService := guides.NewService(
		repositories.Guide,
	)

	return &Services{

		Auth: authService,

		User: users.NewService(
			repositories.User,
			authService,
			guideService,
		),

		Guide: guideService,

		Reservation: reservations.NewService(
			repositories.Reservation,
		),

		TourDay: tourdays.NewService(
			repositories.TourDay,
		),

		Customer: customers.NewService(
			repositories.Customer,
		),

		Agency: agencies.NewService(
			repositories.Agency,
		),

		Language: languages.NewService(
			repositories.Language,
		),

		// Application Services
		ReservationApplication: reservationsapp.NewService(
			repositories.Reservation,
			repositories.Customer,
			repositories.ReservationCustomer,
			repositories.TourDay,
			repositories.Agency,
		),

		ReservationCustomerApplication: reservationcustomersapp.NewService(
			repositories.Reservation,
			repositories.Customer,
			repositories.ReservationCustomer,
		),

		TourDayApplication: tourdaysapp.NewService(
			repositories.TourDay,
			repositories.Guide,
			repositories.TourDayApplication,
		),
	}
}
