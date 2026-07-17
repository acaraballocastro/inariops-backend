package customers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input CreateCustomerInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	if err := h.service.CreateCustomer(input); err != nil {
		http.Error(w, "failed to create customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	customer, err := h.service.GetCustomerByID(id)
	if err != nil {
		http.Error(w, "failed to fetch customer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (h *Handler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := mux.Vars(r)["id"]

	var input UpdateCustomerInput
	input.ID = &id
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateCustomer(input); err != nil {
		http.Error(w, "failed to update customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.service.DeleteCustomer(id); err != nil {
		http.Error(w, "failed to delete customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := h.service.GetAllCustomers()
	if err != nil {
		http.Error(w, "failed to fetch customers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func (h *Handler) SearchCustomers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	customers, err := h.service.SearchCustomers(query)
	if err != nil {
		http.Error(w, "failed to search customers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}
