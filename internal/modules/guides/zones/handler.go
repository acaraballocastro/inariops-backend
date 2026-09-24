package zones

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

func (h *Handler) CreateZone(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var zone CreateZoneRequest

	if err := json.NewDecoder(r.Body).Decode(&zone); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	createdZone, err := h.service.CreateZone(&zone)
	if err != nil {
		http.Error(w, "failed to create zone", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdZone)
}

func (h *Handler) GetZoneByName(w http.ResponseWriter, r *http.Request) {
	zoneName := mux.Vars(r)["name"]

	zone, err := h.service.GetZoneByName(zoneName)
	if err != nil {
		http.Error(w, "failed to get zone", http.StatusInternalServerError)
		return
	}
	if zone == nil {
		http.Error(w, "zone not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zone)
}

func (h *Handler) UpdateZone(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var zone UpdateZoneRequest

	if err := json.NewDecoder(r.Body).Decode(&zone); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	zone.ID = mux.Vars(r)["id"]

	err := h.service.UpdateZone(&zone)
	if err != nil {
		http.Error(w, "failed to update zone", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteZone(w http.ResponseWriter, r *http.Request) {
	zoneName := mux.Vars(r)["name"]

	err := h.service.DeleteZone(zoneName)
	if err != nil {
		http.Error(w, "failed to delete zone", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListZones(w http.ResponseWriter, r *http.Request) {
	zones, err := h.service.ListZones()
	if err != nil {
		http.Error(w, "failed to list zones", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zones)
}
