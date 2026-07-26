package tourdaysapp

type AssignGuideRequest struct {
	TourDayIDs []string `json:"tour_day_ids"`
	GuideID    string   `json:"guide_id"`
}
