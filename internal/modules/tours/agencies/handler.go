package agencies

import (
	"encoding/json"
	"net/http"
	"strconv"

	appErrors "inariops/internal/shared/errors"
	"inariops/internal/shared/logger"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateAgency(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var agency Agency

	if err := json.NewDecoder(r.Body).Decode(&agency); err != nil {
		logger.Error("CreateAgency", "failed to decode request body: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	createdAgency, err := h.service.CreateAgency(&agency)
	if err != nil {
		logger.Error("CreateAgency", "failed to create agency: %v", err)

		switch err {
		case appErrors.ErrExistingAgency:
			http.Error(w, err.Error(), http.StatusConflict)

		case appErrors.ErrInvalidAgency:
			http.Error(w, err.Error(), http.StatusBadRequest)

		default:
			http.Error(w, "failed to create agency", http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdAgency)
}

func (h *Handler) GetAgencyByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	agency, err := h.service.GetAgencyByID(id)
	if err != nil {
		if err == appErrors.ErrAgencyNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "failed to fetch agency", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agency)
}

func (h *Handler) UpdateAgency(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var agency Agency

	if err := json.NewDecoder(r.Body).Decode(&agency); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	agency.ID = mux.Vars(r)["id"]

	if err := h.service.UpdateAgency(&agency); err != nil {
		if err == appErrors.ErrAgencyNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "failed to update agency", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agency)
}

func (h *Handler) DeleteAgency(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.service.DeleteAgency(id); err != nil {
		if err == appErrors.ErrAgencyNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "failed to delete agency", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListAgencies(w http.ResponseWriter, r *http.Request) {
	var req ListAgenciesRequest

	if name := r.URL.Query().Get("name"); name != "" {
		req.Name = &name
	}

	if email := r.URL.Query().Get("email"); email != "" {
		req.Email = &email
	}

	if isActive := r.URL.Query().Get("is_active"); isActive != "" {
		value, err := strconv.ParseBool(isActive)
		if err != nil {
			http.Error(w, "invalid is_active parameter", http.StatusBadRequest)
			return
		}

		req.IsActive = &value
	}

	agencies, err := h.service.ListAgencies(req)
	if err != nil {
		http.Error(w, "failed to fetch agencies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agencies)
}
