package reservationsapp

import (
	"encoding/json"
	"inariops/internal/shared/logger"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetReservationDetailByCode(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	reservationDetail, err := h.service.GetReservationDetailByCode(code)
	if err != nil {
		http.Error(w, "failed to fetch reservation detail", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservationDetail)
}

func (h *Handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reservationRequest CreateReservationRequest

	if err := json.NewDecoder(r.Body).Decode(&reservationRequest); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	reservationDetail, err := h.service.CreateReservation(reservationRequest)
	if err != nil {
		http.Error(w, "failed to create reservation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservationDetail)
}

func (h *Handler) GetCustomersByReservationCode(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	customers, err := h.service.GetCustomersByReservationCode(code)
	if err != nil {
		http.Error(w, "failed to fetch customers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func (h *Handler) DeleteReservation(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	err := h.service.DeleteReservation(code)
	if err != nil {
		http.Error(w, "failed to delete reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateReservation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reservation UpdateReservationRequest

	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		logger.Error("Decode error: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateReservation(reservation)
	if err != nil {
		logger.Error("UpdateReservation error: %v", err)
		http.Error(w, "failed to update reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
