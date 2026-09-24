package reservationsapp

import (
	"database/sql"
	"inariops/internal/db"
	"inariops/internal/domain"
	"inariops/internal/modules/customers"
	"inariops/internal/modules/tours/agencies"
	"inariops/internal/modules/tours/reservations"
	tours "inariops/internal/modules/tours/shared"

	reservationscustomers "inariops/internal/modules/tours/reservations_customers"
	tourdays "inariops/internal/modules/tours/tour_days"
	"inariops/internal/shared/errors"
	"inariops/internal/shared/logger"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	db                       *sql.DB
	reservationsRepo         *reservations.Repository
	customersRepo            *customers.Repository
	reservationCustomersRepo *reservationscustomers.Repository
	tourDaysRepo             *tourdays.Repository
	agenciesRepo             *agencies.Repository
}

func NewService(
	database *sql.DB,
	reservationsRepo *reservations.Repository,
	customersRepo *customers.Repository,
	reservationCustomersRepo *reservationscustomers.Repository,
	tourDaysRepo *tourdays.Repository,
	agenciesRepo *agencies.Repository,
) *Service {
	return &Service{
		db:                       database,
		reservationsRepo:         reservationsRepo,
		customersRepo:            customersRepo,
		reservationCustomersRepo: reservationCustomersRepo,
		tourDaysRepo:             tourDaysRepo,
		agenciesRepo:             agenciesRepo,
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

func (s *Service) CreateReservation(reservationRequest CreateReservationRequest) (ReservationDetail, error) {

	var detail ReservationDetail

	err := db.WithTransaction(s.db, func(tx *sql.Tx) error {

		reservationsRepo := s.reservationsRepo.WithTx(tx)
		customersRepo := s.customersRepo.WithTx(tx)
		reservationCustomersRepo := s.reservationCustomersRepo.WithTx(tx)
		tourDaysRepo := s.tourDaysRepo.WithTx(tx)
		agenciesRepo := s.agenciesRepo.WithTx(tx)

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

		// Validate agency
		if reservation.AgencyID != nil &&
			strings.TrimSpace(*reservation.AgencyID) != "" {

			agency, err := agenciesRepo.GetAgencyByID(*reservation.AgencyID)
			if err != nil {
				logger.Error(
					"CreateReservation: failed to get agency for reservation: %v",
					err,
				)
				return err
			}

			if agency == nil {
				logger.Error(
					"CreateReservation: agency not found: %s",
					*reservation.AgencyID,
				)
				return errors.ErrFailedToCreateReservation
			}
		}

		// Create reservation
		createdReservation, err := reservationsRepo.CreateReservation(reservation)
		if err != nil {
			logger.Error(
				"CreateReservation: failed to create reservation: %v",
				err,
			)
			return err
		}

		// Create tour days
		tourDays := []tours.TourDay{}

		daysRange := int(
			createdReservation.EndDate.
				Sub(createdReservation.StartDate).
				Hours()/24,
		) + 1

		for i := 0; i < daysRange; i++ {

			day := createdReservation.StartDate.AddDate(0, 0, i)

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

			err := tourDaysRepo.CreateTourDay(tourDay)
			if err != nil {
				logger.Error(
					"CreateReservation: failed to create tour day for reservation ID %v: %v",
					createdReservation.ID,
					err,
				)
				return err
			}

			tourDays = append(tourDays, tourDay)
		}

		// Resolve customers
		var customersList []customers.Customer
		var customerIDs []string

		for _, customerRequest := range reservationRequest.Customers {

			if customerRequest.ID == nil ||
				strings.TrimSpace(*customerRequest.ID) == "" {

				customer, err := customersRepo.GetCustomerByEmail(
					customerRequest.Email,
				)

				if err == nil {
					customersList = append(customersList, customer)
					customerIDs = append(customerIDs, customer.ID)
					continue
				}

				customer = customers.Customer{
					ID:             uuid.New().String(),
					FullName:       customerRequest.FullName,
					DocumentNumber: customerRequest.DocumentNumber,
					Phone:          customerRequest.Phone,
					Email:          customerRequest.Email,
					Age:            customerRequest.Age,
					CreatedAt:      time.Now(),
				}

				err = customersRepo.CreateCustomer(customer)
				if err != nil {
					logger.Error(
						"CreateReservation: failed to create customer: %v",
						err,
					)
					return err
				}

				customersList = append(customersList, customer)
				customerIDs = append(customerIDs, customer.ID)

				continue
			}

			customer, err := customersRepo.GetCustomerByEmail(
				customerRequest.Email,
			)

			switch {
			case err == nil:
				customersList = append(customersList, customer)
				customerIDs = append(customerIDs, customer.ID)
				continue

			case err == sql.ErrNoRows:
				// Customer does not exist, create it

			default:
				logger.Error(
					"CreateReservation: failed to find customer by email %s: %v",
					customerRequest.Email,
					err,
				)
				return err
			}

			customer = customers.Customer{
				ID:             uuid.New().String(),
				FullName:       customerRequest.FullName,
				DocumentNumber: customerRequest.DocumentNumber,
				Phone:          customerRequest.Phone,
				Email:          customerRequest.Email,
				Age:            customerRequest.Age,
				CreatedAt:      time.Now(),
			}

			err = customersRepo.CreateCustomer(customer)
			if err != nil {
				logger.Error(
					"CreateReservation: failed to create customer: %v",
					err,
				)
				return err
			}

			customersList = append(customersList, customer)
			customerIDs = append(customerIDs, customer.ID)
		}

		// Link customers to reservation
		err = reservationCustomersRepo.AddCustomerToReservation(
			createdReservation.ID,
			customerIDs,
		)
		if err != nil {
			logger.Error(
				"CreateReservation: failed to add customers to reservation ID %v: %v",
				createdReservation.ID,
				err,
			)
			return err
		}

		detail = ReservationDetail{
			Reservation: createdReservation,
			Customers:   customersList,
			TourDays:    tourDays,
		}

		return nil
	})

	if err != nil {
		return ReservationDetail{}, err
	}

	return detail, nil
}

func (s *Service) DeleteReservation(reservationCode string) error {

	err := db.WithTransaction(s.db, func(tx *sql.Tx) error {

		reservationsRepo := s.reservationsRepo.WithTx(tx)
		tourDaysRepo := s.tourDaysRepo.WithTx(tx)

		// Get reservation
		reservation, err :=
			reservationsRepo.GetReservationByCode(reservationCode)

		if err != nil {
			logger.Error(
				"DeleteReservation: failed to get reservation for code %s: %v",
				reservationCode,
				err,
			)
			return err
		}

		// Get tour days before changing reservation state
		tourDays, err :=
			tourDaysRepo.GetTourDaysByReservationID(reservation.ID)

		if err != nil {
			logger.Error(
				"DeleteReservation: failed to get tour days for reservation ID %s: %v",
				reservation.ID,
				err,
			)
			return err
		}

		// Cancel reservation
		err = reservationsRepo.DeleteReservation(reservationCode)

		if err != nil {
			logger.Error(
				"DeleteReservation: failed to delete reservation for code %s: %v",
				reservationCode,
				err,
			)
			return err
		}

		// Cancel all tour days
		for _, tourDay := range tourDays {

			err = tourDaysRepo.CancelTourDay(tourDay.ID)

			if err != nil {
				logger.Error(
					"DeleteReservation: failed to cancel tour day with ID %s: %v",
					tourDay.ID,
					err,
				)
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateReservation(reservation UpdateReservationRequest) error {

	if reservation.Code == "" {
		logger.Error("updateReservation: invalid input - code is required")
		return errors.ErrInvalidReservationCode
	}

	logger.Info(
		"updateReservation: updating reservation with code %s",
		reservation.Code,
	)

	err := db.WithTransaction(s.db, func(tx *sql.Tx) error {

		reservationsRepo := s.reservationsRepo.WithTx(tx)
		customersRepo := s.customersRepo.WithTx(tx)
		reservationCustomersRepo := s.reservationCustomersRepo.WithTx(tx)
		agenciesRepo := s.agenciesRepo.WithTx(tx)

		// Get current reservation
		oldReservation, err := reservationsRepo.GetReservationByCode(
			reservation.Code,
		)
		if err != nil {
			logger.Error(
				"updateReservation: failed to retrieve reservation with code %s: %v",
				reservation.Code,
				err,
			)
			return errors.ErrReservationNotFound
		}

		// Validate agency
		if reservation.AgencyID != nil &&
			strings.TrimSpace(*reservation.AgencyID) != "" {

			agency, err := agenciesRepo.GetAgencyByID(
				*reservation.AgencyID,
			)

			if err != nil {
				logger.Error(
					"updateReservation: failed to get agency for reservation: %v",
					err,
				)
				return errors.ErrFailedToUpdateReservation
			}

			if agency == nil {
				logger.Error(
					"updateReservation: agency not found: %s",
					*reservation.AgencyID,
				)
				return errors.ErrFailedToUpdateReservation
			}
		}

		// Update reservation fields
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
			oldReservation.TotalPeopleCount =
				reservation.TotalPeopleCount
		}

		oldReservation.StartDate = reservation.StartDate
		oldReservation.EndDate = reservation.EndDate
		oldReservation.UpdatedAt = time.Now()

		err = reservationsRepo.UpdateReservation(oldReservation)
		if err != nil {
			logger.Error(
				"updateReservation: failed to update reservation %s: %v",
				reservation.Code,
				err,
			)
			return errors.ErrFailedToUpdateReservation
		}

		// Get current customer relations
		existingCustomerIDs, err :=
			reservationCustomersRepo.GetCustomersByReservationID(
				oldReservation.ID,
			)

		if err != nil {
			logger.Error(
				"updateReservation: failed to get existing customers for reservation %s: %v",
				reservation.Code,
				err,
			)
			return errors.ErrFailedToUpdateReservation
		}

		existingCustomers := make([]customers.Customer, 0, len(existingCustomerIDs))

		for _, customerID := range existingCustomerIDs {

			customer, err :=
				customersRepo.GetCustomerByID(customerID)

			if err != nil {
				logger.Error(
					"updateReservation: failed to get customer %s for reservation %s: %v",
					customerID,
					reservation.Code,
					err,
				)
				return errors.ErrFailedToUpdateReservation
			}

			existingCustomers = append(
				existingCustomers,
				customer,
			)
		}

		// Empty customer list means remove all relations
		if len(reservation.Customers) == 0 {

			for _, customer := range existingCustomers {

				err = reservationCustomersRepo.
					RemoveCustomerFromReservation(
						oldReservation.ID,
						customer.ID,
					)

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

		existingCustomerSet :=
			make(map[string]struct{}, len(existingCustomers))

		for _, customer := range existingCustomers {
			existingCustomerSet[customer.ID] = struct{}{}
		}

		desiredCustomerIDs :=
			make([]string, 0, len(reservation.Customers))

		desiredCustomerSet :=
			make(map[string]struct{}, len(reservation.Customers))

		// Resolve requested customers
		for _, customerRequest := range reservation.Customers {

			var customerID string

			// Customer identified by ID
			if customerRequest.ID != nil &&
				strings.TrimSpace(*customerRequest.ID) != "" {

				customer, err :=
					customersRepo.GetCustomerByID(
						strings.TrimSpace(*customerRequest.ID),
					)

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

			} else {

				// Try to find customer by email
				customer, err :=
					customersRepo.GetCustomerByEmail(
						customerRequest.Email,
					)

				switch {
				case err == nil:
					customerID = customer.ID

				case err == sql.ErrNoRows:
					// Customer does not exist, create it
					newCustomer := customers.Customer{
						ID:             uuid.New().String(),
						FullName:       customerRequest.FullName,
						DocumentNumber: customerRequest.DocumentNumber,
						Phone:          customerRequest.Phone,
						Email:          customerRequest.Email,
						Age:            customerRequest.Age,
						CreatedAt:      time.Now(),
					}

					err = customersRepo.CreateCustomer(
						newCustomer,
					)

					if err != nil {
						logger.Error(
							"updateReservation: failed to create customer for reservation %s: %v",
							reservation.Code,
							err,
						)
						return errors.ErrFailedToUpdateReservation
					}

					customerID = newCustomer.ID

				default:
					logger.Error(
						"updateReservation: failed to find customer by email for reservation %s: %v",
						reservation.Code,
						err,
					)
					return errors.ErrFailedToUpdateReservation
				}
			}

			desiredCustomerIDs =
				append(desiredCustomerIDs, customerID)

			desiredCustomerSet[customerID] = struct{}{}
		}

		// Add missing relations
		for _, customerID := range desiredCustomerIDs {

			if _, exists := existingCustomerSet[customerID]; exists {
				continue
			}

			err = reservationCustomersRepo.
				AddCustomerToReservation(
					oldReservation.ID,
					[]string{customerID},
				)

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

		// Remove obsolete relations
		for _, customer := range existingCustomers {

			if _, exists := desiredCustomerSet[customer.ID]; exists {
				continue
			}

			err = reservationCustomersRepo.
				RemoveCustomerFromReservation(
					oldReservation.ID,
					customer.ID,
				)

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
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SyncReservationStatus(reservationCode string) error {
	reservation, err := s.reservationsRepo.GetReservationByCode(reservationCode)
	if err != nil {
		logger.Error("SyncReservationStatus: failed to get reservation for code %s: %v", reservationCode, err)
		return err
	}

	tourDays, err := s.tourDaysRepo.GetTourDaysByReservationID(reservation.ID)
	if err != nil {
		logger.Error("SyncReservationStatus: failed to get tour days for reservation ID %s: %v", reservation.ID, err)
		return err
	}

	newStatus := domain.RESERVATION_PENDING_ASSIGNMENT
	allGuideConfirmed := len(tourDays) > 0
	for _, tourDay := range tourDays {
		if tourDay.Status != domain.RESERVATION_GUIDE_CONFIRMED {
			allGuideConfirmed = false
		}
		if tourDay.Status == domain.RESERVATION_PAYMENT_PENDING {
			newStatus = domain.RESERVATION_PAYMENT_PENDING
			break
		} else if tourDay.Status == domain.RESERVATION_GUIDE_PREASSIGNED && newStatus != domain.RESERVATION_PAYMENT_PENDING {
			newStatus = domain.RESERVATION_GUIDE_PREASSIGNED
		} else if tourDay.Status == domain.RESERVATION_PENDING_ASSIGNMENT && newStatus != domain.RESERVATION_PAYMENT_PENDING && newStatus != domain.RESERVATION_GUIDE_PREASSIGNED {
			newStatus = domain.RESERVATION_PENDING_ASSIGNMENT
		}
	}
	if allGuideConfirmed {
		newStatus = domain.RESERVATION_CONFIRMED
	}

	if reservation.Status != newStatus {
		logger.Info("SyncReservationStatus: updating reservation %s status from %s to %s", reservationCode, reservation.Status, newStatus)
		reservation.Status = newStatus
		reservation.UpdatedAt = time.Now()
		err = s.reservationsRepo.UpdateReservationStatus(*reservation.Code, newStatus)
		if err != nil {
			logger.Error("SyncReservationStatus: failed to update reservation status for code %s: %v", reservationCode, err)
			return err
		}
	}
	logger.Info("SyncReservationStatus: reservation %s status updated to %s", reservationCode, newStatus)

	return nil
}

func (s *Service) UpdateSignatureStatus(code string, status domain.SignatureStatus) error {
	return s.reservationsRepo.UpdateSignatureStatus(code, status)
}

func (s *Service) UpdateVoucherStatus(code string, status domain.VoucherStatus) error {
	return s.reservationsRepo.UpdateVoucherStatus(code, status)
}

func (s *Service) EraseReservation(reservationCode string) error {

	err := db.WithTransaction(s.db, func(tx *sql.Tx) error {

		reservationsRepo := s.reservationsRepo.WithTx(tx)
		tourDaysRepo := s.tourDaysRepo.WithTx(tx)
		reservationCustomersRepo := s.reservationCustomersRepo.WithTx(tx)

		// Get reservation
		reservation, err :=
			reservationsRepo.GetReservationByCode(reservationCode)

		if err != nil {
			logger.Error(
				"EraseReservation: failed to get reservation for code %s: %v",
				reservationCode,
				err,
			)
			return err
		}

		// Get tour days
		tourDays, err :=
			tourDaysRepo.GetTourDaysByReservationID(reservation.ID)

		if err != nil {
			logger.Error(
				"EraseReservation: failed to get tour days for reservation ID %s: %v",
				reservation.ID,
				err,
			)
			return err
		}

		// Delete tour days
		for _, tourDay := range tourDays {

			err = tourDaysRepo.DeleteTourDaysByID(tourDay.ID)

			if err != nil {
				logger.Error(
					"EraseReservation: failed to delete tour day with ID %s: %v",
					tourDay.ID,
					err,
				)
				return err
			}
		}

		// Get reservation customers
		customerIDs, err :=
			reservationCustomersRepo.GetCustomersByReservationID(
				reservation.ID,
			)

		if err != nil {
			logger.Error(
				"EraseReservation: failed to get customers for reservation ID %s: %v",
				reservation.ID,
				err,
			)
			return err
		}

		// Delete customer relations
		for _, customerID := range customerIDs {

			err =
				reservationCustomersRepo.RemoveCustomerFromReservation(
					reservation.ID,
					customerID,
				)

			if err != nil {
				logger.Error(
					"EraseReservation: failed to remove customer %s from reservation ID %s: %v",
					customerID,
					reservation.ID,
					err,
				)
				return err
			}
		}

		// Finally delete reservation
		err = reservationsRepo.EraseReservation(reservationCode)

		if err != nil {
			logger.Error(
				"EraseReservation: failed to delete reservation for code %s: %v",
				reservationCode,
				err,
			)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
