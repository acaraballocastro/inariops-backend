package reservations

import (
	"inariops/internal/domain"
	tours "inariops/internal/modules/tours/shared"
	tourdays "inariops/internal/modules/tours/tour_days"
	"inariops/internal/shared/errors"
	"inariops/internal/shared/logger"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo        *Repository
	tourDayRepo *tourdays.Repository
}

func NewService(repo *Repository, tourDayRepo *tourdays.Repository) *Service {
	return &Service{repo: repo, tourDayRepo: tourDayRepo}
}

func (s *Service) GetAllReservations() ([]Reservation, error) {
	reservations, err := s.repo.GetAllReservations()
	if err != nil {
		logger.Error("GetAllReservations: failed to retrieve all reservations: %v", err)
	}
	return reservations, err
}

func (s *Service) GetReservationByCode(code string) (Reservation, error) {
	if code == "" {
		logger.Error("GetReservationByCode: invalid reservation code - code is empty")
		return Reservation{}, errors.ErrInvalidInput
	}

	reservation, err := s.repo.GetReservationByCode(code)
	if err != nil {
		logger.Error("GetReservationByCode: failed to get reservation for code %s: %v", code, err)
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

	err := s.repo.CreateReservation(reservation)
	if err != nil {
		logger.Error("CreateReservation: failed to create reservation: %v", err)
		return err
	}
	days := int(reservation.EndDate.Sub(reservation.StartDate).Hours()/24) + 1

	for i := 0; i < days; i++ {
		tourDay := tours.TourDay{
			ID:            uuid.New().String(),
			ReservationID: reservation.ID,
			Title:         reservation.Title,
			StartDateTime: reservation.StartDate.Add(time.Duration(i) * 24 * time.Hour),
			Duration:      nil,
			Status:        domain.RESERVATION_PENDING_ASSIGNMENT,
			VoucherStatus: domain.VOUCHER_NOT_GENERATED,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		err := s.tourDayRepo.CreateTourDay(tourDay)
		if err != nil {
			logger.Error("CreateReservation: failed to create tour day for reservation %s: %v", reservation.ID, err)
			return err
		}
	}

	return nil
}

func (s *Service) UpdateReservation(reservation Reservation) error {
	if reservation.Code == nil || *reservation.Code == "" {
		logger.Error("UpdateReservation: invalid input - code is required")
		return errors.ErrInvalidReservationCode
	}

	err := s.repo.UpdateReservation(reservation)
	if err != nil {
		logger.Error(
			"UpdateReservation: failed to update reservation %s: %v",
			*reservation.Code,
			err,
		)
		return err
	}

	return nil
}

func (s *Service) DeleteReservation(code string) error {
	if code == "" {

		return errors.ErrInvalidInput
	}

	return s.repo.DeleteReservation(code)
}

func (s *Service) GetReservationsByAgencyID(agencyID string) ([]Reservation, error) {
	return []Reservation{}, errors.ErrNotFound
}
