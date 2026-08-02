package tourdaysapp

type AssignGuideRequest struct {
	TourDayIDs []string `json:"tour_day_ids"`
	GuideID    string   `json:"guide_id"`
}

type TourDayAssigmentHistory struct {
	ID        string `json:"id"`
	TourDayID string `json:"tour_day_id"`
	GuideID   string `json:"guide_id"`

	Action    string `json:"action"`
	ChangedBy string `json:"changed_by"`
	Reason    string `json:"reason"`

	CreatedAt string `json:"created_at"`
}

type TourDayStatusHistory struct {
	ID             string `json:"id"`
	TourDayID      string `json:"tour_day_id"`
	PreviousStatus string `json:"previous_status"`
	NewStatus      string `json:"new_status"`
	ChangedBy      string `json:"changed_by"`
	Reason         string `json:"reason"`
	CreatedAt      string `json:"created_at"`
}
