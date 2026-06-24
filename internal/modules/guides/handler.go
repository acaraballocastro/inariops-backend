package guides

import (
	"encoding/json"
	"inariops/internal/domain"
	"inariops/internal/shared/response"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetGuideByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	guide, err := h.service.GetGuideByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, guide)
}

func (h *Handler) GetGuideByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	guide, err := h.service.GetGuideByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, guide)
}

func (h *Handler) GetAllGuides(w http.ResponseWriter, r *http.Request) {
	guides, err := h.service.GetAllGuides()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, guides)
}

func (h *Handler) UpdateGuide(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var guide domain.Guide
	if err := json.NewDecoder(r.Body).Decode(&guide); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	guide.ID = id

	if err := h.service.UpdateGuide(guide); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Guide updated successfully"})
}
