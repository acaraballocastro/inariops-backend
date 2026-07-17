package customers

import "time"

type Customer struct {
	ID             string
	FullName       string
	DocumentNumber string
	Phone          string
	Email          string
	Age            int
	CreatedAt      time.Time
}

type CreateCustomerInput struct {
	FullName       string `json:"full_name" validate:"required"`
	DocumentNumber string `json:"document_number" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Age            int    `json:"age" validate:"required,gte=0"`
}

type UpdateCustomerInput struct {
	ID             *string `json:"id" validate:"required"`
	FullName       *string `json:"full_name" validate:"required"`
	DocumentNumber *string `json:"document_number" validate:"required"`
	Phone          *string `json:"phone" validate:"required"`
	Email          *string `json:"email" validate:"required,email"`
	Age            *int    `json:"age" validate:"required,gte=0"`
}

type CustomerResponse struct {
	ID             string `json:"id"`
	FullName       string `json:"full_name"`
	DocumentNumber string `json:"document_number"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Age            int    `json:"age"`
}
