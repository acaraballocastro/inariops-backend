package users

import (
	"encoding/json"
	"inariops/internal/shared/logger"
	"inariops/internal/shared/response"
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

// DTO separado (mejor práctica)
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

func (h *Handler) GetAllGuides(w http.ResponseWriter, r *http.Request) {
	guides, err := h.service.GetAllGuides()
	if err != nil {
		http.Error(w, "failed to fetch guides", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guides)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, "failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.service.GetUserByID(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	logger.Info("create user request")

	var input CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("invalid body")
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	createdUser, err := h.service.CreateUser(input.Name, input.Email, input.Phone, input.Role)
	if err != nil {
		logger.Error("create user failed: %v", err)
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// err = authService.CreateCredentials(createdUser.ID, createdUser.Email)

	response.JSON(w, http.StatusCreated, createdUser)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input UpdateUserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	input.ID = mux.Vars(r)["id"]

	updatedUser, err := h.service.UpdateUser(input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, updatedUser)
}

func (h *Handler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		response.Error(w, http.StatusBadRequest, "user id is required")
		return
	}

	err := h.service.DeactivateUser(id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to deactivate user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
