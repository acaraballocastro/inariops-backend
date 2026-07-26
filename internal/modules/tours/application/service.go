package reservationsapp

import (
	"fmt"
	"inariops/internal/domain"
	"inariops/internal/modules/customers"
	"inariops/internal/modules/tours/reservations"
	reservationscustomers "inariops/internal/modules/tours/reservations_customers"
	tours "inariops/internal/modules/tours/shared"
	tourdays "inariops/internal/modules/tours/tour_days"
	"inariops/internal/shared/errors"
	"inariops/internal/shared/logger"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	reservationsRepo         *reservations.Repository
	customersRepo            *customers.Repository
	reservationCustomersRepo *reservationscustomers.Repository
	tourDaysRepo             *tourdays.Repository
}

func NewService(
	reservationsRepo *reservations.Repository,
	customersRepo *customers.Repository,
	reservationCustomersRepo *reservationscustomers.Repository,
	tourDaysRepo *tourdays.Repository,
) *Service {
	return &Service{
		reservationsRepo:         reservationsRepo,
		customersRepo:            customersRepo,
		reservationCustomersRepo: reservationCustomersRepo,
		tourDaysRepo:             tourDaysRepo,
	}
}

func (s *Service) GetReservationDetailByCode(reservationCode string) (ReservationDetail, error) {
	var detail ReservationDetail

	reservation, err := s.reservationsRepo.GetReservationByCode(reservationCode)
	if err != nil {
		return detail, err
	}
	detail.Reservation = reservation

	customerIDs, err := s.reservationCustomersRepo.GetCustomersByReservationID(reservation.ID)
	if err != nil {
		return detail, err
	}

	var customersList []customers.Customer
	for _, id := range customerIDs {
		customer, err := s.customersRepo.GetCustomerByID(id)
		if err != nil {
			return detail, err
		}
		customersList = append(customersList, customer)
	}
	detail.Customers = customersList

	tourDays, err := s.tourDaysRepo.GetTourDaysByReservationID(reservation.ID)
	if err != nil {
		return detail, err
	}
	detail.TourDays = tourDays

	return detail, nil
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

func (s *Service) CreateReservation(reservationRequest CreateReservationRequest) (ReservationDetail, error) {
	reservation := reservations.Reservation{
		ID:               uuid.New().String(),
		Title:            reservationRequest.Title,
		Description:      reservationRequest.Description,
		AgencyID:         reservationRequest.AgencyID,
		TotalPeopleCount: reservationRequest.TotalPeopleCount,
		StartDate:        reservationRequest.StartDate,
		EndDate:          reservationRequest.EndDate,
		Status:           domain.RESERVATION_PENDING_ASSIGNMENT,
		VoucherStatus:    domain.VOUCHER_NOT_GENERATED,
		SignatureStatus:  domain.SIGNATURE_NOT_SENT,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	createdReservation, err := s.reservationsRepo.CreateReservation(reservation)
	logger.Info("CreateReservation: created reservation with ID %v", createdReservation.ID)
	if err != nil {
		logger.Error("CreateReservation: failed to create reservation: %v", err)
		return ReservationDetail{}, err
	}

	daysRange := int(reservation.EndDate.Sub(reservation.StartDate).Hours()/24) + 1

	var tourDays []tours.TourDay
	for i := 0; i < daysRange; i++ {
		day := reservation.StartDate.AddDate(0, 0, i)

		tourDay := tours.TourDay{
			ID:            uuid.New().String(),
			ReservationID: createdReservation.ID,
			Title:         createdReservation.Title,
			StartDateTime: day,
			PeopleCount:   createdReservation.TotalPeopleCount,
			Duration:      nil,
			Status:        domain.RESERVATION_PENDING_ASSIGNMENT,
			VoucherStatus: domain.VOUCHER_NOT_GENERATED,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		err := s.tourDaysRepo.CreateTourDay(tourDay)
		logger.Info("CreateReservation: created tour day with ID %v for reservation ID %v", tourDay.ID, createdReservation.ID)
		if err != nil {
			logger.Error("CreateReservation: failed to create tour day for reservation ID %v: %v", createdReservation.ID, err)
			return ReservationDetail{}, err
		}
		tourDays = append(tourDays, tourDay)
	}

	var customersList []customers.Customer
	var customerIDs []string
	for _, customerRequest := range reservationRequest.Customers {
		if customerRequest.ID == nil {
			customer := customers.Customer{
				ID:             uuid.New().String(),
				FullName:       customerRequest.FullName,
				DocumentNumber: customerRequest.DocumentNumber,
				Phone:          customerRequest.Phone,
				Email:          customerRequest.Email,
				Age:            customerRequest.Age,
				CreatedAt:      time.Now(),
			}
			err := s.customersRepo.CreateCustomer(customer)
			logger.Info("CreateReservation: created customer with ID %v", customer.ID)
			if err != nil {
				logger.Error("CreateReservation: failed to create customer: %v", err)
				return ReservationDetail{}, err
			}
			customerRequest.ID = &customer.ID
			customersList = append(customersList, customer)
			customerIDs = append(customerIDs, customer.ID)
		} else {
			customer, err := s.customersRepo.GetCustomerByID(*customerRequest.ID)
			if err != nil {
				return ReservationDetail{}, err
			}
			customersList = append(customersList, customer)
			customerIDs = append(customerIDs, customer.ID)
		}
	}
	err = s.reservationCustomersRepo.AddCustomerToReservation(createdReservation.ID, customerIDs)
	logger.Info("CreateReservation: added customers to reservation ID %v", createdReservation.ID)

	if err != nil {
		logger.Error("CreateReservation: failed to add customers to reservation ID %v: %v", createdReservation.ID, err)
		return ReservationDetail{}, err
	}

	return ReservationDetail{
		Reservation: createdReservation,
		Customers:   customersList,
		TourDays:    tourDays,
	}, nil
}

func (s *Service) DeleteReservation(reservationCode string) error {
	reservation, err := s.reservationsRepo.GetReservationByCode(reservationCode)
	if err != nil {
		logger.Error("DeleteReservation: failed to get reservation for code %s: %v", reservationCode, err)
		return err
	}

	err = s.reservationsRepo.DeleteReservation(reservationCode)
	if err != nil {
		logger.Error("DeleteReservation: failed to delete reservation for code %s: %v", reservationCode, err)
		return err
	}

	tourDays, err := s.tourDaysRepo.GetTourDaysByReservationID(reservation.ID)
	if err != nil {
		logger.Error("DeleteReservation: failed to get tour days for reservation ID %s: %v", reservation.ID, err)
		return err
	}

	for _, tourDay := range tourDays {
		err := s.tourDaysRepo.CancelTourDay(tourDay.ID)
		if err != nil {
			logger.Error("DeleteReservation: failed to cancel tour day with ID %s: %v", tourDay.ID, err)
			return err
		}
	}

	return nil
}

func (s *Service) UpdateReservation(reservation UpdateReservationRequest) error {
	if reservation.Code == "" {
		logger.Error("updateReservation: invalid input - code is required")
		return errors.ErrInvalidReservationCode
	}

	logger.Info("updateReservation: updating reservation with code %s", reservation.Code)
	oldReservation, err := s.reservationsRepo.GetReservationByCode(reservation.Code)
	if err != nil {
		logger.Error("updateReservation: failed to retrieve reservation with code %s: %v", reservation.Code, err)
		return errors.ErrReservationNotFound
	}

	if reservation.Title != nil {
		oldReservation.Title = reservation.Title
	}
	if reservation.Description != nil {
		oldReservation.Description = reservation.Description
	}
	if reservation.QuoteNumber != nil {
		oldReservation.QuoteNumber = reservation.QuoteNumber
	}
	if reservation.FileNumber != nil {
		oldReservation.FileNumber = reservation.FileNumber
	}
	if reservation.AgencyID != nil {
		oldReservation.AgencyID = reservation.AgencyID
	}
	if reservation.TotalPeopleCount != nil {
		oldReservation.TotalPeopleCount = reservation.TotalPeopleCount
	}
	oldReservation.StartDate = reservation.StartDate
	oldReservation.EndDate = reservation.EndDate
	oldReservation.UpdatedAt = time.Now()

	err = s.reservationsRepo.UpdateReservation(oldReservation)
	if err != nil {
		logger.Error(
			"updateReservation: failed to update reservation %s: %v",
			reservation.Code,
			err,
		)
		return errors.ErrFailedToUpdateReservation
	}

	existingCustomers, err := s.GetCustomersByReservationCode(reservation.Code)

	if len(reservation.Customers) == 0 {
		for _, customer := range existingCustomers {
			err = s.reservationCustomersRepo.RemoveCustomerFromReservation(oldReservation.ID, customer.ID)
			if err != nil {
				logger.Error(
					"updateReservation: failed to remove customer %s from reservation %s: %v",
					customer.ID,
					reservation.Code,
					err,
				)
				return errors.ErrFailedToUpdateReservation
			}
		}

		return nil
	}

	existingCustomerIDs := make(map[string]struct{}, len(existingCustomers))
	for _, customer := range existingCustomers {
		existingCustomerIDs[customer.ID] = struct{}{}
	}

	desiredCustomerIDs := make([]string, 0, len(reservation.Customers))
	desiredCustomerSet := make(map[string]struct{}, len(reservation.Customers))
	for _, customerRequest := range reservation.Customers {
		var customerID string

		if customerRequest.ID == nil || strings.TrimSpace(*customerRequest.ID) == "" {
			customer := customers.Customer{
				ID:             uuid.New().String(),
				FullName:       customerRequest.FullName,
				DocumentNumber: customerRequest.DocumentNumber,
				Phone:          customerRequest.Phone,
				Email:          customerRequest.Email,
				Age:            customerRequest.Age,
				CreatedAt:      time.Now(),
			}

			err = s.customersRepo.CreateCustomer(customer)
			if err != nil {
				logger.Error(
					"updateReservation: failed to create customer for reservation %s: %v",
					reservation.Code,
					err,
				)
				return errors.ErrFailedToUpdateReservation
			}

			customerID = customer.ID
		} else {
			customer, err := s.customersRepo.GetCustomerByID(*customerRequest.ID)
			if err != nil {
				logger.Error(
					"updateReservation: failed to get customer %s for reservation %s: %v",
					strings.TrimSpace(*customerRequest.ID),
					reservation.Code,
					err,
				)
				return errors.ErrFailedToUpdateReservation
			}

			customerID = customer.ID
		}

		desiredCustomerIDs = append(desiredCustomerIDs, customerID)
		desiredCustomerSet[customerID] = struct{}{}
	}

	for _, customerID := range desiredCustomerIDs {
		if _, exists := existingCustomerIDs[customerID]; exists {
			continue
		}

		err = s.reservationCustomersRepo.AddCustomerToReservation(oldReservation.ID, []string{customerID})
		if err != nil {
			logger.Error(
				"updateReservation: failed to add customer %s to reservation %s: %v",
				customerID,
				reservation.Code,
				err,
			)
			return errors.ErrFailedToUpdateReservation
		}
	}

	for _, customer := range existingCustomers {
		if _, exists := desiredCustomerSet[customer.ID]; exists {
			continue
		}

		err = s.reservationCustomersRepo.RemoveCustomerFromReservation(oldReservation.ID, customer.ID)
		if err != nil {
			logger.Error(
				"updateReservation: failed to remove customer %s from reservation %s: %v",
				customer.ID,
				reservation.Code,
				err,
			)
			return errors.ErrFailedToUpdateReservation
		}
	}
	return nil
}
