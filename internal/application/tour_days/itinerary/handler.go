package itineraryapp

import (
	"encoding/json"
	"inariops/internal/modules/itinerary"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetItineraryByTourDayID(w http.ResponseWriter, r *http.Request) {
	tourDayID := mux.Vars(r)["tour_day_id"]
	items, err := h.service.GetItineraryByTourDayID(tourDayID)
	if err != nil {
		http.Error(w, "failed to fetch itinerary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *Handler) CreateItineraryItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var itineraryItem itinerary.CreateItineraryItemRequest

	if err := json.NewDecoder(r.Body).Decode(&itineraryItem); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.CreateItineraryItem(mux.Vars(r)["tour_day_id"], itineraryItem)
	if err != nil {
		http.Error(w, "failed to create itinerary item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(itineraryItem)
}

func (h *Handler) GetItineraryItemByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	itineraryItem, err := h.service.GetItineraryItemByID(id)
	if err != nil {
		http.Error(w, "failed to fetch itinerary item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(itineraryItem)
}

func (h *Handler) UpdateItineraryItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var itineraryItem itinerary.CreateItineraryItemRequest

	if err := json.NewDecoder(r.Body).Decode(&itineraryItem); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateItineraryItem(mux.Vars(r)["item_id"], itineraryItem)
	if err != nil {
		http.Error(w, "failed to update itinerary item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(itineraryItem)
}

func (h *Handler) DeleteItineraryItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["item_id"]

	err := h.service.DeleteItineraryItem(id)
	if err != nil {
		http.Error(w, "failed to delete itinerary item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
