package itinerary

type ItineraryItem struct {
	ID         string  `json:"id"`
	TourDayID  string  `json:"tour_day_id"`
	OrderIndex int     `json:"order_index"`
	StartTime  *string `json:"start_time"`
	EndTime    *string `json:"end_time"`
	Type       *string `json:"type"`
	PlaceID    *string `json:"place_id"`
	ActivityID *string `json:"activity_id"`
	Notes      *string `json:"notes"`
}

type CreateItineraryItemRequest struct {
	OrderIndex int     `json:"order_index"`
	StartTime  *string `json:"start_time"`
	EndTime    *string `json:"end_time"`
	Type       *string `json:"type"`
	PlaceID    *string `json:"place_id"`
	ActivityID *string `json:"activity_id"`
	Notes      *string `json:"notes"`
}
