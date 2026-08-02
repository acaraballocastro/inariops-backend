package reservationcustomersapp

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
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
