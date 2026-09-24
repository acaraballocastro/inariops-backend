package bootstrap

import (
	"database/sql"
	guidesapp "inariops/internal/application/guides"
	reservationcustomersapp "inariops/internal/application/reservation_customers"
	reservationsapp "inariops/internal/application/reservations"
	tourdaysapp "inariops/internal/application/tour_days"
	itineraryapp "inariops/internal/application/tour_days/itinerary"
	placesapp "inariops/internal/application/tour_days/places"

	"inariops/internal/modules/auth"
	"inariops/internal/modules/customers"
	"inariops/internal/modules/guides"
	"inariops/internal/modules/guides/languages"
	"inariops/internal/modules/guides/zones"
	"inariops/internal/modules/itinerary/activity"
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

	Activity *activity.Service

	Zone *zones.Service

	ReservationApplication *reservationsapp.Service

	ReservationCustomerApplication *reservationcustomersapp.Service

	TourDayApplication *tourdaysapp.Service

	GuideApplication *guidesapp.Service

	ItineraryApplication *itineraryapp.Service

	PlaceApplication *placesapp.Service
}

func newServices(
	database *sql.DB,
	repositories *Repositories,
	jwtManager *auth.JWTManager,
) *Services {

	authService := auth.NewService(
		repositories.Auth,
		jwtManager,
	)

	guideService := guides.NewService(
		repositories.Guide,
	)

	userService := users.NewService(
		database,
		repositories.User,
		authService,
		guideService,
	)

	return &Services{

		Auth: authService,

		User: userService,

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

		Zone: zones.NewService(
			repositories.Zone,
		),

		Activity: activity.NewService(
			repositories.Activity,
		),

		// Application Services
		ReservationApplication: reservationsapp.NewService(
			database,
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
			database,
			repositories.TourDay,
			repositories.Guide,
			repositories.Availability,
			repositories.TourDayApplication,
		),

		GuideApplication: guidesapp.NewService(
			database,
			userService,
			repositories.User,
			repositories.Guide,
			repositories.Language,
			repositories.LanguageGuide,
			repositories.Zone,
			repositories.ZoneGuide,
			repositories.Availability,
		),

		ItineraryApplication: itineraryapp.NewService(
			repositories.TourDay,
			repositories.Activity,
			repositories.Place,
			repositories.ItineraryItem,
		),

		PlaceApplication: placesapp.NewService(
			repositories.Place,
			repositories.Zone,
		),
	}
}
