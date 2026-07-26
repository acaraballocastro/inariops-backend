package tourdaysapp

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

func (h *Handler) AssignGuide(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request AssignGuideRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.AssignGuide(request.TourDayIDs, request.GuideID)
	if err != nil {
		http.Error(w, "failed to assign guide", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UnassignGuide(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request AssignGuideRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.UnassignGuide(request.TourDayIDs)
	if err != nil {
		http.Error(w, "failed to unassign guide", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
