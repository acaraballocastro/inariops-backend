package reservationscustomers

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

func (h *Handler) AddCustomerToReservation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	reservationCode := vars["reservation_code"]

	var req UpsertCustomerToReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.AddCustomerToReservation(reservationCode, req.CustomerID); err != nil {
		logger.Error("Error adding customer to reservation: %v", err)
		http.Error(w, "failed to add customer to reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RemoveCustomerFromReservation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	reservationCode := vars["reservation_code"]

	var req UpsertCustomerToReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.RemoveCustomerFromReservation(reservationCode, req.CustomerID); err != nil {
		logger.Error("Error removing customer from reservation: %v", err)
		http.Error(w, "failed to remove customer from reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetCustomersByReservationCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reservationCode := vars["reservation_code"]

	customerList, err := h.service.GetCustomersByReservationCode(reservationCode)
	if err != nil {
		logger.Error("Error fetching customers by reservation code: %v", err)
		http.Error(w, "failed to fetch customers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customerList)
}
