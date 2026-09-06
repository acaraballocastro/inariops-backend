package placesapp

import (
	"encoding/json"
	"inariops/internal/modules/itinerary/places"
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

func (h *Handler) CreatePlace(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request places.CreatePlaceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.CreatePlace(&request)
	if err != nil {
		http.Error(w, "failed to create place", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(request)
}

func (h *Handler) GetPlacesByTourDayID(w http.ResponseWriter, r *http.Request) {
	places, err := h.service.GetAllPlaces()
	if err != nil {
		http.Error(w, "failed to fetch places", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(places)
}

func (h *Handler) GetPlaceByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	place, err := h.service.GetPlaceByID(id)
	if err != nil {
		http.Error(w, "failed to fetch place", http.StatusInternalServerError)
		return
	}

	if place == nil {
		http.Error(w, "place not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(place)
}

func (h *Handler) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request places.CreatePlaceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	id := mux.Vars(r)["id"]
	err := h.service.UpdatePlace(id, &request)
	if err != nil {
		http.Error(w, "failed to update place", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeletePlace(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := h.service.DeletePlace(id)
	if err != nil {
		http.Error(w, "failed to delete place", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
