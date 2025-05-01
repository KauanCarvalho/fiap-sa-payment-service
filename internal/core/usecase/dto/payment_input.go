package dto

import "github.com/go-playground/validator/v10"

type AuthorizePaymentInput struct {
	Amount            float64 `json:"amount" validate:"required,gt=0"`
	ExternalReference string  `json:"external_reference" validate:"required"`
	PaymentMethod     string  `json:"payment_method" validate:"required"`
}

func ValidatePaymentCreate(input AuthorizePaymentInput) error {
	return validator.New().Struct(input)
}
