package reservations

import (
	"encoding/json"
	"io"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleReservations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Get("id") != "" {
			h.GetReservationByID(w, r)
		} else {
			h.GetAllReservations(w, r)
		}
	case http.MethodPost:
		h.CreateReservation(w, r)
	case http.MethodPatch:
		h.UpdateReservation(w, r)
	case http.MethodDelete:
		h.DeleteReservation(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) GetAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := h.service.GetAllReservations()
	if err != nil {
		http.Error(w, "failed to fetch reservations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}

func (h *Handler) GetReservationByID(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("id")

	reservation, err := h.service.GetReservationByID(code)
	if err != nil {
		http.Error(w, "failed to fetch reservation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}

func (h *Handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
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
	var reservation Reservation

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(bodyBytes, &reservation); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateReservation(reservation)
	if err != nil {
		http.Error(w, "failed to update reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteReservation(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("id")

	err := h.service.DeleteReservation(code)
	if err != nil {
		http.Error(w, "failed to delete reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
