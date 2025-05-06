package dto

import "github.com/go-playground/validator/v10"

type AuthorizePaymentInput struct {
	Amount            float64 `json:"amount" validate:"required,gt=0"`
	ExternalReference string  `json:"external_reference" validate:"required"`
	PaymentMethod     string  `json:"payment_method" validate:"required"`
}

type UpdatePaymentStatusInput struct {
	ExternalReference string `json:"external_reference" validate:"required"`
	Status            string `json:"status" validate:"required,oneof=completed failed"`
}

func ValidatePaymentCreate(input AuthorizePaymentInput) error {
	return validator.New().Struct(input)
}

func ValidatePaymentStatusUpdate(input UpdatePaymentStatusInput) error {
	return validator.New().Struct(input)
}
