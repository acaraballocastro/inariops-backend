package reservationscustomers

import (
	"inariops/internal/modules/customers"
	"inariops/internal/modules/tours/reservations"
	"inariops/internal/shared/errors"
)

type Service struct {
	repo             *Repository
	reservationsRepo *reservations.Repository
	customersRepo    *customers.Repository
}

func NewService(repo *Repository, reservationsRepo *reservations.Repository, customersRepo *customers.Repository) *Service {
	return &Service{
		repo:             repo,
		reservationsRepo: reservationsRepo,
		customersRepo:    customersRepo,
	}
}

func (s *Service) AddCustomerToReservation(reservationCode string, customerIDs []string) error {
	var reservation reservations.Reservation

	for _, customerID := range customerIDs {
		if _, err := s.customersRepo.GetCustomerByID(customerID); err != nil {
			return errors.ErrCustomerNotFound
		}
	}

	reservation, err := s.reservationsRepo.GetReservationByCode(reservationCode)
	if err != nil {
		return errors.ErrReservationNotFound
	}

	return s.repo.AddCustomerToReservation(reservation.ID, customerIDs)
}

func (s *Service) RemoveCustomerFromReservation(reservationCode string, customerIDs []string) error {
	reservation, err := s.reservationsRepo.GetReservationByCode(reservationCode)
	if err != nil {
		return errors.ErrReservationNotFound
	}

	for _, customerID := range customerIDs {
		if _, err := s.customersRepo.GetCustomerByID(customerID); err != nil {
			return errors.ErrCustomerNotFound
		}

		s.repo.RemoveCustomerFromReservation(reservation.ID, customerID)
		if err != nil {
			return err
		}

	}
	return nil
}

func (s *Service) GetCustomersByReservationCode(reservationCode string) ([]customers.Customer, error) {
	var customerList []customers.Customer

	reservation, err := s.reservationsRepo.GetReservationByCode(reservationCode)
	if err != nil {
		return nil, errors.ErrReservationNotFound
	}

	reservationCustomer, err := s.repo.GetCustomersByReservationID(reservation.ID)

	for _, customerID := range reservationCustomer {
		customer, err := s.customersRepo.GetCustomerByID(customerID)
		if err != nil {
			return nil, errors.ErrCustomerNotFound
		}
		customerList = append(customerList, customer)
	}

	return customerList, nil
}
