package dto

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
