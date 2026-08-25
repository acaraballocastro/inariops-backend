package bootstrap

import (
	guidesapp "inariops/internal/application/guides"
	reservationcustomersapp "inariops/internal/application/reservation_customers"
	reservationsapp "inariops/internal/application/reservations"
	tourdaysapp "inariops/internal/application/tour_days"

	"inariops/internal/modules/auth"
	"inariops/internal/modules/customers"
	"inariops/internal/modules/guides"
	"inariops/internal/modules/guides/languages"
	"inariops/internal/modules/guides/zones"
	"inariops/internal/modules/tours/agencies"
	"inariops/internal/modules/tours/reservations"
	tourdays "inariops/internal/modules/tours/tour_days"
	"inariops/internal/modules/users"
)

type Handlers struct {
	Auth *auth.Handler

	User *users.Handler

	Guide *guides.Handler

	Reservation *reservations.Handler

	TourDay *tourdays.Handler

	Customer *customers.Handler

	Agency *agencies.Handler

	Language *languages.Handler

	Zone *zones.Handler

	ReservationApplication *reservationsapp.Handler

	ReservationCustomerApplication *reservationcustomersapp.Handler

	TourDayApplication *tourdaysapp.Handler

	GuideApplication *guidesapp.Handler
}

func newHandlers(
	services *Services,
) *Handlers {

	return &Handlers{

		Auth: auth.NewHandler(
			services.Auth,
		),

		User: users.NewHandler(
			services.User,
		),

		Guide: guides.NewHandler(
			services.Guide,
		),

		Reservation: reservations.NewHandler(
			services.Reservation,
		),

		TourDay: tourdays.NewHandler(
			services.TourDay,
		),

		Customer: customers.NewHandler(
			services.Customer,
		),

		Agency: agencies.NewHandler(
			services.Agency,
		),

		Language: languages.NewHandler(
			services.Language,
		),

		Zone: zones.NewHandler(
			services.Zone,
		),

		// Application Handlers
		ReservationApplication: reservationsapp.NewHandler(
			services.ReservationApplication,
		),

		ReservationCustomerApplication: reservationcustomersapp.NewHandler(
			services.ReservationCustomerApplication,
		),

		TourDayApplication: tourdaysapp.NewHandler(
			services.TourDayApplication,
		),

		GuideApplication: guidesapp.NewHandler(
			services.GuideApplication,
		),
	}
}
