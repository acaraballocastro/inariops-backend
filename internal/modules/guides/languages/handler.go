package languages

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

func (h *Handler) CreateLanguage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var languageReq Language

	if err := json.NewDecoder(r.Body).Decode(&languageReq); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	language, err := h.service.CreateLanguage(&languageReq)
	if err != nil {
		http.Error(w, "failed to create Language", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(language)
}

func (h *Handler) GetLanguageByCode(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	Language, err := h.service.GetLanguageByCode(code)
	if err != nil {
		http.Error(w, "failed to fetch Language", http.StatusInternalServerError)
		return
	}

	if Language == nil {
		http.Error(w, "Language not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Language)
}

func (h *Handler) UpdateLanguage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var Language Language

	if err := json.NewDecoder(r.Body).Decode(&Language); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	Language.Code = mux.Vars(r)["code"]

	if err := h.service.UpdateLanguage(&Language); err != nil {
		http.Error(w, "failed to update Language", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Language)
}

func (h *Handler) DeleteLanguage(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	if err := h.service.DeleteLanguage(code); err != nil {
		http.Error(w, "failed to delete Language", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListLanguages(w http.ResponseWriter, r *http.Request) {
	languages, err := h.service.ListLanguages()
	if err != nil {
		http.Error(w, "failed to fetch languages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(languages)
}
