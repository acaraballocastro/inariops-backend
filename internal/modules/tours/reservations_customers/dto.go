package reservationscustomers

type ReservationCustomer struct {
	ReservationID string `json:"reservation_id"`
	CustomerID    string `json:"customer_id"`
}

type UpsertCustomerToReservationRequest struct {
	CustomerID []string `json:"customer_id"`
}
