package reservationsapp

import (
	"fmt"
	"inariops/internal/domain"
	"inariops/internal/modules/customers"
	"inariops/internal/modules/tours/reservations"
	reservationscustomers "inariops/internal/modules/tours/reservations_customers"
	tours "inariops/internal/modules/tours/shared"
	tourdays "inariops/internal/modules/tours/tour_days"
	"inariops/internal/shared/logger"
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
		return err
	}

	err = s.reservationsRepo.DeleteReservation(reservationCode)
	if err != nil {
		return err
	}

	tourDays, err := s.tourDaysRepo.GetTourDaysByReservationID(reservation.ID)
	if err != nil {
		return err
	}

	for _, tourDay := range tourDays {
		err := s.tourDaysRepo.DeleteTourDay(tourDay.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
