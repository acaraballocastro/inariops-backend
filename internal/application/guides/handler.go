package guidesapp

import (
	"database/sql"
	"encoding/json"
	"errors"
	"inariops/internal/modules/guides/availabilities"
	appErrors "inariops/internal/shared/errors"
	"inariops/internal/shared/logger"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetGuideDetailByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	guide, err := h.service.GetGuideDetailByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guide)
}

func (h *Handler) GetAllGuidesDetail(w http.ResponseWriter, r *http.Request) {
	guides, err := h.service.GetAllGuidesDetail()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guides)
}

func (h *Handler) AddLanguageToGuide(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	defer r.Body.Close()
	var languageGuide LanguageGuide

	if err := json.NewDecoder(r.Body).Decode(&languageGuide); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	languageGuide.GuideID = id
	if err := h.service.AddLanguageToGuide(languageGuide.GuideID, languageGuide.LanguageID); err != nil {
		logger.Error("Error adding language to guide: %v", err)
		http.Error(w, "failed to add language to guide", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) RemoveLanguageFromGuide(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var languageGuide LanguageGuide

	if err := json.NewDecoder(r.Body).Decode(&languageGuide); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.RemoveLanguageFromGuide(languageGuide.GuideID, languageGuide.LanguageID); err != nil {
		http.Error(w, "failed to remove language from guide", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetLanguagesByGuideID(w http.ResponseWriter, r *http.Request) {
	guideID := mux.Vars(r)["id"]

	languages, err := h.service.GetLanguagesByGuideID(guideID)
	if err != nil {
		http.Error(w, "failed to fetch languages for guide", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(languages)
}

func (h *Handler) CreateGuide(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var guide CreateGuideRequest
	if err := json.NewDecoder(r.Body).Decode(&guide); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	createdGuide, err := h.service.CreateGuide(guide)
	if err != nil {
		http.Error(w, "failed to create guide", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdGuide)
}

func (h *Handler) UpdateGuide(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	var guide UpdateGuideRequest
	if err := json.NewDecoder(r.Body).Decode(&guide); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	updatedGuide, err := h.service.UpdateGuide(id, guide)
	if err != nil {
		http.Error(w, "failed to update guide", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedGuide)
}

// Availability Handlers
func (h *Handler) GetAvailabilitiesByGuideID(w http.ResponseWriter, r *http.Request) {
	guideID := mux.Vars(r)["id"]

	availabilities, err := h.service.GetAvailabilitiesByGuideID(guideID)
	if err != nil {
		http.Error(w, "failed to fetch availabilities for guide", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(availabilities)
}

func (h *Handler) CreateAvailability(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var availability availabilities.CreateAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&availability); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	availability.GuideID = mux.Vars(r)["id"]

	createdAvailability, err := h.service.CreateAvailability(availability)
	if err != nil {
		writeAvailabilityError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdAvailability)
}

func (h *Handler) DeleteAvailability(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["availability_id"]

	err := h.service.DeleteAvailability(id)
	if err != nil {
		http.Error(w, "failed to delete availability", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	guideID := vars["id"]
	availabilityID := vars["availability_id"]

	var availabilityRequest availabilities.CreateAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&availabilityRequest); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	updatedAvailability, err := h.service.UpdateAvailability(guideID, availabilityID, availabilityRequest)
	if err != nil {
		writeAvailabilityError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedAvailability)
}

func writeAvailabilityError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, appErrors.ErrInvalidAvailabilityDate):
		http.Error(w, "invalid availability dates", http.StatusBadRequest)
	case errors.Is(err, appErrors.ErrAvailabilityConflict):
		http.Error(w, "availability range overlaps an existing range", http.StatusConflict)
	case errors.Is(err, sql.ErrNoRows), errors.Is(err, appErrors.ErrAvailabilityNotFound):
		http.Error(w, "guide or availability not found", http.StatusNotFound)
	default:
		http.Error(w, "failed to process availability", http.StatusInternalServerError)
	}
}
