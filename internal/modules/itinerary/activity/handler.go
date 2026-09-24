package activity

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

func (h *Handler) GetActivityByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	activityID := vars["id"]

	activity, err := h.service.GetActivityByID(activityID)
	if err != nil {
		http.Error(w, "failed to fetch activity", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activity)
}

func (h *Handler) GetAllActivities(w http.ResponseWriter, r *http.Request) {
	activities, err := h.service.GetAllActivities()
	if err != nil {
		http.Error(w, "failed to fetch activities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}

func (h *Handler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var activity CreateActivityRequest

	if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateActivity(&activity); err != nil {
		http.Error(w, "failed to create activity", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var activity UpdateActivityRequest

	if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateActivity(&activity); err != nil {
		http.Error(w, "failed to update activity", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	activityID := vars["id"]

	if err := h.service.DeleteActivity(activityID); err != nil {
		http.Error(w, "failed to delete activity", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
