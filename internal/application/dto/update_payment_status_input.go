package dto

type UpdatePaymentStatusInput struct {
	ExternalRef string `json:"external_reference"`
	Status      string `json:"status"`
}
