package reservations

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// TODO: Add pagination
func (h *Handler) GetAllReservations(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	reservations, err := h.service.GetAllReservations()
	if err != nil {
		http.Error(w, "failed to fetch reservations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}

func (h *Handler) GetReservationByCode(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	reservation, err := h.service.GetReservationByCode(code)
	if err != nil {
		http.Error(w, "failed to fetch reservation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}

func (h *Handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reservation Reservation

	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.CreateReservation(reservation)
	if err != nil {
		http.Error(w, "failed to create reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateReservation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var reservation Reservation

	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateReservation(reservation)
	if err != nil {
		http.Error(w, "failed to update reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// TODO: Think to add a soft delete for reservations, so we can keep track of deleted reservations and their history. This will include a new field in the reservation model to indicate if a reservation is deleted, as well as functions to handle soft deletion and retrieval of deleted reservations.
func (h *Handler) DeleteReservation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	code := r.URL.Query().Get("code")

	err := h.service.DeleteReservation(code)
	if err != nil {
		http.Error(w, "failed to delete reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//TODO: Add a few functions and handlers to handle the signature and voucher status of a reservation. This will include functions to update the signature status and voucher status, as well as handlers to handle the corresponding HTTP requests.

//TODO: Add filters to search reservations by status, date range, and other criteria. This will include functions to filter reservations based on the specified criteria, as well as handlers to handle the corresponding HTTP requests.

//TODO: Add a SearchByCode

//TODO: Add userAudit for reservation creation, update, and deletion. This will include functions to log user actions related to reservations, as well as handlers to handle the corresponding HTTP requests.
