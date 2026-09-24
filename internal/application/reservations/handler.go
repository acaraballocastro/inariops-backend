package reservationsapp

import (
	"encoding/json"
	"inariops/internal/domain"
	"inariops/internal/shared/logger"
	"net/http"

	"github.com/gorilla/mux"
)

type StatusUpdateRequest struct {
	Status string `json:"status"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetReservationDetailByCode(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	reservationDetail, err := h.service.GetReservationDetailByCode(code)
	if err != nil {
		http.Error(w, "failed to fetch reservation detail", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservationDetail)
}

func (h *Handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reservationRequest CreateReservationRequest

	if err := json.NewDecoder(r.Body).Decode(&reservationRequest); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	reservationDetail, err := h.service.CreateReservation(reservationRequest)
	if err != nil {
		http.Error(w, "failed to create reservation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservationDetail)
}

func (h *Handler) DeleteReservation(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	err := h.service.DeleteReservation(code)
	if err != nil {
		http.Error(w, "failed to delete reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateReservation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reservation UpdateReservationRequest

	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		logger.Error("Decode error: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateReservation(reservation)
	if err != nil {
		logger.Error("UpdateReservation error: %v", err)
		http.Error(w, "failed to update reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) EraseReservation(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	err := h.service.EraseReservation(code)
	if err != nil {
		logger.Error("EraseReservation error: %v", err)
		http.Error(w, "failed to erase reservation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateSignatureStatus(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var request StatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	status := domain.SignatureStatus(request.Status)
	if status != domain.SIGNATURE_NOT_SENT && status != domain.SIGNATURE_SENT && status != domain.SIGNATURE_SIGNED && status != domain.SIGNATURE_REJECTED {
		http.Error(w, "invalid signature status", http.StatusBadRequest)
		return
	}
	if err := h.service.UpdateSignatureStatus(mux.Vars(r)["code"], status); err != nil {
		http.Error(w, "failed to update signature status", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateVoucherStatus(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var request StatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	status := domain.VoucherStatus(request.Status)
	if status != domain.VOUCHER_NOT_GENERATED && status != domain.VOUCHER_GENERATED && status != domain.VOUCHER_PARTIALLY_SENT && status != domain.VOUCHER_SENT {
		http.Error(w, "invalid voucher status", http.StatusBadRequest)
		return
	}
	if err := h.service.UpdateVoucherStatus(mux.Vars(r)["code"], status); err != nil {
		http.Error(w, "failed to update voucher status", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
