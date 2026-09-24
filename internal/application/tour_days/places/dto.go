package placesapp

import (
	"inariops/internal/modules/guides/zones"
	"inariops/internal/modules/itinerary/places"
)

type PlaceWithZone struct {
	Place *places.Place
	Zone  *zones.Zone
}
