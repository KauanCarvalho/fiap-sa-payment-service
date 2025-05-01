package mappers

import (
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain/entities"
)

func ToPaymentDTO(payment entities.Payment) dto.PaymentOutput {
	return dto.PaymentOutput{
		Amount:            payment.Amount,
		Status:            payment.Status,
		ExternalReference: payment.ExternalReference,
		PaymentMethod:     payment.PaymentMethod,
		Provider:          payment.Provider,
		QRCode:            payment.QRCode,
	}
}
