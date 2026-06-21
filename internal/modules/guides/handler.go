package guides

import "inariops/internal/domain"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateGuide(guide domain.Guide) error {
	return h.service.CreateGuide(guide)
}

func (h *Handler) GetGuideByID(id string) (domain.Guide, error) {
	return h.service.GetGuideByID(id)
}

func (h *Handler) GetGuideByUserID(userID string) (domain.Guide, error) {
	return h.service.GetGuideByUserID(userID)
}

func (h *Handler) GetAllGuides() ([]domain.Guide, error) {
	return h.service.GetAllGuides()
}

func (h *Handler) UpdateGuide(guide domain.Guide) error {
	return h.service.UpdateGuide(guide)
}

func (h *Handler) DeleteGuide(id string) error {
	return h.service.DeleteGuide(id)
}
