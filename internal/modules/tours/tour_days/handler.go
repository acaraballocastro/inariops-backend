package tourdays

import (
	"encoding/json"
	"fmt"
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

func (h *Handler) GetTourDaysByReservationID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	reservationID := vars["reservation_id"]

	tourDays, err := h.service.GetTourDaysByReservationID(reservationID)
	if err != nil {
		http.Error(w, "failed to fetch tour days", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tourDays)
}

func (h *Handler) CancelTourDay(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CancelTourDayInput

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.CancelTourDay(req.ID)
	if err != nil {
		http.Error(w, "failed to cancel tour day", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetTourDayByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tourID := vars["id"]
	tourDay, err := h.service.GetTourDaysByTourID(tourID)
	if err != nil {
		http.Error(w, "failed to fetch tour day", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tourDay)
}

func (h *Handler) CreateTourDay(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var tourDayInput CreateTourDayInput

	if err := json.NewDecoder(r.Body).Decode(&tourDayInput); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.CreateTourDay(tourDayInput)
	if err != nil {
		http.Error(w, "failed to create tour day", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateTourDay(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var tourDayInput UpdateTourDayInput
	muxVars := mux.Vars(r)
	tourDayInput.ID = muxVars["id"]

	if err := json.NewDecoder(r.Body).Decode(&tourDayInput); err != nil {
		logger.Error("Decode error: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateTourDay(tourDayInput)
	if err != nil {
		http.Error(w, "failed to update tour day", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AssignGuide(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var guideID string
	var tourDays []string
	muxVars := mux.Vars(r)
	guideID = muxVars["guide_id"]

	if err := json.NewDecoder(r.Body).Decode(&tourDays); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.AssignGuide(tourDays, guideID)
	if err != nil {
		message := "failed to assign guide to tour day: %v"
		http.Error(w, fmt.Sprintf(message, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UnassignGuide(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var guideID string
	var tourDays []string
	muxVars := mux.Vars(r)
	guideID = muxVars["guide_id"]

	if err := json.NewDecoder(r.Body).Decode(&tourDays); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.UnassignGuide(tourDays, guideID)
	if err != nil {
		message := "failed to unassign guide from tour day: %v"
		http.Error(w, fmt.Sprintf(message, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
