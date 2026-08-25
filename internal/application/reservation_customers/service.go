package reservationcustomersapp

import (
	"fmt"
	"inariops/internal/modules/customers"
	"inariops/internal/modules/tours/reservations"
	reservationscustomers "inariops/internal/modules/tours/reservations_customers"
	"inariops/internal/shared/logger"
)

type Service struct {
	reservationsRepo         *reservations.Repository
	customersRepo            *customers.Repository
	reservationCustomersRepo *reservationscustomers.Repository
}

func NewService(
	reservationsRepo *reservations.Repository,
	customersRepo *customers.Repository,
	reservationCustomersRepo *reservationscustomers.Repository,
) *Service {
	return &Service{
		reservationsRepo:         reservationsRepo,
		customersRepo:            customersRepo,
		reservationCustomersRepo: reservationCustomersRepo,
	}
}

func (s *Service) AddCustomerToReservation(reservationCode string, customerIDs []string) error {
	reservation, err := s.reservationsRepo.GetReservationByCode(reservationCode)
	if err != nil {
		return err
	}

	var customersList []string
	for _, customerID := range customerIDs {
		customer, err := s.customersRepo.GetCustomerByID(customerID)
		if err != nil {
			return err
		}

		if s.reservationCustomersRepo.IsCustomerInReservation(reservation.ID, customer.ID) {
			return fmt.Errorf("customer with ID %s is already in the reservation", customer.ID)
		}

		customersList = append(customersList, customer.ID)
	}

	err = s.reservationCustomersRepo.AddCustomerToReservation(reservation.ID, customersList)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveCustomersFromReservation(reservationCode string, customerIDs []string) error {
	reservation, err := s.reservationsRepo.GetReservationByCode(reservationCode)
	if err != nil {
		return err
	}

	for _, customerID := range customerIDs {
		err := s.reservationCustomersRepo.RemoveCustomerFromReservation(reservation.ID, customerID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetCustomersByReservationCode(reservationCode string) ([]customers.Customer, error) {
	reservation, err := s.reservationsRepo.GetReservationByCode(reservationCode)
	if err != nil {
		return nil, err
	}

	customerIDs, err := s.reservationCustomersRepo.GetCustomersByReservationID(reservation.ID)

	if len(customerIDs) == 0 {
		return []customers.Customer{}, nil
	}

	var customersList []customers.Customer
	for _, id := range customerIDs {
		customer, err := s.customersRepo.GetCustomerByID(id)
		if err != nil {
			return nil, err
		}
		customersList = append(customersList, customer)
	}

	return customersList, nil
}

func (s *Service) DeleteCustomer(customerID string) error {
	// Check if the customer exists
	customer, err := s.customersRepo.GetCustomerByID(customerID)
	if err != nil {
		logger.Error("Error retrieving customer with ID %s: %v", customerID, err)
		return err
	}
	if customer.ID == "" {
		logger.Error("Customer with ID %s not found", customerID)
		return fmt.Errorf("customer with ID %s not found", customerID)
	}

	// Remove the customer from all reservations
	err = s.reservationCustomersRepo.RemoveAllReservationsFromCustomer(customerID)
	if err != nil {
		logger.Error("Error removing customer with ID %s from reservations: %v", customerID, err)
		return err
	}

	// Delete the customer from the customers table
	err = s.customersRepo.DeleteCustomer(customerID)
	if err != nil {
		logger.Error("Error deleting customer with ID %s: %v", customerID, err)
		return err
	}

	return nil
}
