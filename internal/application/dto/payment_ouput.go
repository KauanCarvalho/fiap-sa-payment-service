package dto

type PaymentOutput struct {
	Amount            float64 `json:"amount"`
	Status            string  `json:"status"`
	ExternalReference string  `json:"external_reference"`
	Provider          string  `json:"provider"`
	PaymentMethod     string  `json:"payment_method"`
	QRCode            string  `json:"qr_code"`
}
