package reservations

import (
	"inariops/internal/domain"
	"inariops/internal/shared/errors"
	"log"
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
	return s.repo.GetAllReservations()
}

func (s *Service) GetReservationByID(code string) (Reservation, error) {
	reservation, err := s.repo.GetReservationByID(code)
	if err != nil {
		log.Printf("GetReservationByID: failed to get reservation for code %s: %v", code, err)
		return Reservation{}, errors.ErrReservationNotFound
	}

	return reservation, nil
}

func (s *Service) CreateReservation(reservation Reservation) error {
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
		log.Printf("UpdateReservation: invalid reservation code")
		return errors.ErrInvalidReservationCode
	}

	existingReservation, err := s.repo.GetReservationByID(*reservation.Code)
	if err != nil {
		log.Printf("UpdateReservation: failed to get reservation for code %s: %v", *reservation.Code, err)
		return errors.ErrReservationNotFound
	}
	existingReservation = reservation
	existingReservation.UpdatedAt = time.Now()

	err = s.repo.UpdateReservation(existingReservation)
	if err != nil {
		log.Printf("UpdateReservation: failed to persist reservation for code %s: %v", *reservation.Code, err)
		return err
	}

	return nil
}

func (s *Service) DeleteReservation(code string) error {
	reservation, err := s.repo.GetReservationByID(code)
	if err != nil {
		log.Printf("DeleteReservation: failed to get reservation for code %s: %v", code, err)
		return errors.ErrReservationNotFound
	}

	return s.repo.DeleteReservation(*reservation.Code)
}

func (s *Service) GetReservationsByAgencyID(agencyID string) ([]Reservation, error) {
	//TODO: Implement this method to fetch reservations by agency ID
	return nil, nil
}
