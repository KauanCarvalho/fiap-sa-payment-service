package mappers_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase/mappers"
	"github.com/stretchr/testify/assert"
)

func TestToPaymentDTO(t *testing.T) {
	t.Run("should map Payment entity to PaymentOutput DTO", func(t *testing.T) {
		payment := entities.Payment{
			Amount:            1500,
			Status:            "authorized",
			ExternalReference: "ext_123",
			Provider:          "MercadoPago",
			PaymentMethod:     "pix",
			QRCode:            "some-qr-code",
		}

		expected := dto.PaymentOutput{
			Amount:            1500,
			Status:            "authorized",
			ExternalReference: "ext_123",
			Provider:          "MercadoPago",
			PaymentMethod:     "pix",
			QRCode:            "some-qr-code",
		}

		result := mappers.ToPaymentDTO(payment)

		assert.Equal(t, expected, result)
	})
}
