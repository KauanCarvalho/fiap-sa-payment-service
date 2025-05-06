package dto

import "github.com/KauanCarvalho/fiap-sa-payment-service/internal/shared/validation"

type APIErrorOutput struct {
	Details string `json:"details,omitempty"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type APIErrorsOutput struct {
	Errors []APIErrorOutput `json:"errors"`
}

func SimpleAPIErrorsOutput(details, field, message string) APIErrorsOutput {
	return APIErrorsOutput{
		Errors: []APIErrorOutput{
			{
				Details: details,
				Field:   field,
				Message: message,
			},
		},
	}
}

func ErrorsFromValidationErrors(validationErrors []validation.ErrorResponse) APIErrorsOutput {
	errors := make([]APIErrorOutput, len(validationErrors))
	for i, validationError := range validationErrors {
		errors[i] = APIErrorOutput{
			Field:   validationError.Field,
			Message: validationError.Message,
		}
	}
	return APIErrorsOutput{Errors: errors}
}
