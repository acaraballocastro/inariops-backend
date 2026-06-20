package reservations

import (
	"inariops/internal/domain"
	"inariops/internal/shared/errors"
	"inariops/internal/shared/logger"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllReservations() ([]Reservation, error) {
	reservations, err := s.repo.GetAllReservations()
	if err != nil {
		logger.Error("GetAllReservations: failed to retrieve all reservations: %v", err)
	}
	return reservations, err
}

func (s *Service) GetReservationByID(code string) (Reservation, error) {
	if code == "" {
		logger.Error("GetReservationByID: invalid reservation code - code is empty")
		return Reservation{}, errors.ErrInvalidInput
	}

	reservation, err := s.repo.GetReservationByID(code)
	if err != nil {
		logger.Error("GetReservationByID: failed to get reservation for code %s: %v", code, err)
		return Reservation{}, errors.ErrReservationNotFound
	}

	return reservation, nil
}

func (s *Service) CreateReservation(reservation Reservation) error {
	if reservation.Code == nil || *reservation.Code == "" {
		logger.Error("CreateReservation: invalid input - code is required")
		return errors.ErrInvalidInput
	}

	reservation.Status = domain.RESERVATION_PENDING_ASSIGNMENT
	reservation.VoucherStatus = domain.VOUCHER_NOT_GENERATED
	reservation.SignatureStatus = domain.SIGNATURE_NOT_SENT
	reservation.ID = uuid.New().String()
	reservation.CreatedAt = time.Now()
	reservation.UpdatedAt = time.Now()

	return s.repo.CreateReservation(reservation)
}

func (s *Service) UpdateReservation(reservation Reservation) error {
	if reservation.Code == nil || *reservation.Code == "" {
		logger.Error("UpdateReservation: invalid input - code is required")
		return errors.ErrInvalidReservationCode
	}

	existingReservation, err := s.repo.GetReservationByID(*reservation.Code)
	if err != nil {
		logger.Error("UpdateReservation: failed to get reservation for code %s: %v", *reservation.Code, err)
		return errors.ErrReservationNotFound
	}

	reservation.CreatedAt = existingReservation.CreatedAt
	reservation.ID = existingReservation.ID
	reservation.UpdatedAt = time.Now()

	err = s.repo.UpdateReservation(reservation)
	if err != nil {
		logger.Error("UpdateReservation: failed to persist reservation for code %s: %v", *reservation.Code, err)
		return err
	}

	logger.Info("UpdateReservation: reservation updated successfully for code %s", *reservation.Code)
	return nil
}

func (s *Service) DeleteReservation(code string) error {
	if code == "" {
		logger.Error("DeleteReservation: invalid input - code is required")
		return errors.ErrInvalidInput
	}

	reservation, err := s.repo.GetReservationByID(code)
	if err != nil {
		logger.Error("DeleteReservation: failed to get reservation for code %s: %v", code, err)
		return errors.ErrReservationNotFound
	}

	err = s.repo.DeleteReservation(*reservation.Code)
	if err != nil {
		logger.Error("DeleteReservation: failed to delete reservation for code %s: %v", code, err)
		return err
	}

	logger.Info("DeleteReservation: reservation deleted successfully for code %s", code)
	return nil
}

func (s *Service) GetReservationsByAgencyID(agencyID string) ([]Reservation, error) {
	return []Reservation{}, errors.ErrNotFound
}
