package tourdaysapp

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

	err := h.service.UnassignGuide(request.TourDayIDs, request.GuideID)
	logger.Info("error %s", err)
	if err != nil {
		http.Error(w, "failed to unassign guide", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ConfirmTourDay(w http.ResponseWriter, r *http.Request) {
	if err := h.service.ConfirmTourDay(mux.Vars(r)["id"]); err != nil {
		http.Error(w, "failed to confirm tour day", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetTourDaysAvailableForGuide(w http.ResponseWriter, r *http.Request) {
	guideID := mux.Vars(r)["guide_id"]
	tourDays, err := h.service.TourDaysAvailableForGuide(guideID)
	if err != nil {
		http.Error(w, "failed to fetch available tour days", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tourDays)
}
