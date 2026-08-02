package guidesapp

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
