package customers

import (
	"inariops/internal/shared/errors"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCustomer(input CreateCustomerInput) error {
	customer := Customer{
		ID:             uuid.New().String(),
		FullName:       input.FullName,
		DocumentNumber: input.DocumentNumber,
		Phone:          input.Phone,
		Email:          input.Email,
		Age:            *input.Age,
		CreatedAt:      time.Now(),
	}

	return s.repo.CreateCustomer(customer)
}

func (s *Service) GetCustomerByID(id string) (Customer, error) {
	customer, err := s.repo.GetCustomerByID(id)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func (s *Service) UpdateCustomer(input UpdateCustomerInput) error {
	if input.ID == nil {
		return errors.ErrInvalidInput
	}

	customer, err := s.repo.GetCustomerByID(*input.ID)
	if err != nil {
		return errors.ErrCustomerNotFound
	}

	if input.FullName != nil {
		customer.FullName = *input.FullName
	}
	if input.DocumentNumber != nil {
		customer.DocumentNumber = *input.DocumentNumber
	}
	if input.Phone != nil {
		customer.Phone = *input.Phone
	}
	if input.Email != nil {
		customer.Email = *input.Email
	}
	if input.Age != nil {
		customer.Age = *input.Age
	}

	return s.repo.UpdateCustomer(customer)
}

func (s *Service) GetAllCustomers() ([]Customer, error) {
	customers, err := s.repo.GetAllCustomers()
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (s *Service) SearchCustomers(query string) ([]Customer, error) {
	customers, err := s.repo.SearchCustomers(query)
	if err != nil {
		return nil, err
	}

	return customers, nil
}
