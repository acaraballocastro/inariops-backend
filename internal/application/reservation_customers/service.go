package reservationcustomersapp

import (
	"database/sql"
	"fmt"
	"inariops/internal/db"
	"inariops/internal/modules/customers"
	"inariops/internal/modules/tours/reservations"
	reservationscustomers "inariops/internal/modules/tours/reservations_customers"
	"inariops/internal/shared/logger"
)

type Service struct {
	db                       *sql.DB
	reservationsRepo         *reservations.Repository
	customersRepo            *customers.Repository
	reservationCustomersRepo *reservationscustomers.Repository
}

func NewService(
	database *sql.DB,
	reservationsRepo *reservations.Repository,
	customersRepo *customers.Repository,
	reservationCustomersRepo *reservationscustomers.Repository,
) *Service {
	return &Service{
		db:                       database,
		reservationsRepo:         reservationsRepo,
		customersRepo:            customersRepo,
		reservationCustomersRepo: reservationCustomersRepo,
	}
}

func (s *Service) AddCustomerToReservation(
	reservationCode string,
	customerIDs []string,
) error {
	return db.WithTransaction(nil, func(tx *sql.Tx) error {

		reservationsRepo := s.reservationsRepo.WithTx(tx)
		customersRepo := s.customersRepo.WithTx(tx)
		reservationCustomersRepo := s.reservationCustomersRepo.WithTx(tx)

		reservation, err := reservationsRepo.GetReservationByCode(reservationCode)
		if err != nil {
			return err
		}

		var customersList []string

		for _, customerID := range customerIDs {
			customer, err := customersRepo.GetCustomerByID(customerID)
			if err != nil {
				return err
			}

			if reservationCustomersRepo.IsCustomerInReservation(
				reservation.ID,
				customer.ID,
			) {
				return fmt.Errorf(
					"customer with ID %s is already in the reservation",
					customer.ID,
				)
			}

			customersList = append(customersList, customer.ID)
		}

		if err := reservationCustomersRepo.AddCustomerToReservation(
			reservation.ID,
			customersList,
		); err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) RemoveCustomersFromReservation(
	reservationCode string,
	customerIDs []string,
) error {
	return db.WithTransaction(s.db, func(tx *sql.Tx) error {

		reservationsRepo := s.reservationsRepo.WithTx(tx)
		reservationCustomersRepo := s.reservationCustomersRepo.WithTx(tx)

		reservation, err := reservationsRepo.GetReservationByCode(reservationCode)
		if err != nil {
			return err
		}

		for _, customerID := range customerIDs {
			if err := reservationCustomersRepo.RemoveCustomerFromReservation(
				reservation.ID,
				customerID,
			); err != nil {
				return err
			}
		}

		return nil
	})
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
	return db.WithTransaction(s.db, func(tx *sql.Tx) error {

		customersRepo := s.customersRepo.WithTx(tx)
		reservationCustomersRepo := s.reservationCustomersRepo.WithTx(tx)

		customer, err := customersRepo.GetCustomerByID(customerID)
		if err != nil {
			logger.Error(
				"Error retrieving customer with ID %s: %v",
				customerID,
				err,
			)
			return err
		}

		if customer.ID == "" {
			logger.Error(
				"Customer with ID %s not found",
				customerID,
			)
			return fmt.Errorf(
				"customer with ID %s not found",
				customerID,
			)
		}

		if err := reservationCustomersRepo.RemoveAllReservationsFromCustomer(
			customerID,
		); err != nil {
			logger.Error(
				"Error removing customer with ID %s from reservations: %v",
				customerID,
				err,
			)
			return err
		}

		if err := customersRepo.DeleteCustomer(customerID); err != nil {
			logger.Error(
				"Error deleting customer with ID %s: %v",
				customerID,
				err,
			)
			return err
		}

		return nil
	})
}
