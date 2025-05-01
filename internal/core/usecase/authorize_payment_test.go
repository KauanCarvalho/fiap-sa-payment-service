package usecase_test

import (
	"context"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorizePaymentUseCase(t *testing.T) {
	t.Run("successfully authorize payment", func(t *testing.T) {
		input := dto.AuthorizePaymentInput{
			Amount:            1000,
			ExternalReference: uuid.New().String(),
			PaymentMethod:     "credit_card",
		}

		payment, err := ap.Run(ctx, input)
		require.NoError(t, err, "Run should not return error")

		assert.InEpsilon(t, input.Amount, payment.Amount, 0.01)
		assert.Equal(t, input.ExternalReference, payment.ExternalReference)
		assert.Equal(t, input.PaymentMethod, payment.PaymentMethod)
		assert.NotEmpty(t, payment.QRCode, "QRCode should not be empty")
		assert.Equal(t, "MercadoPago", payment.Provider, "Provider should be 'MercadoPago'")
		assert.Equal(t, entities.PaymentStatusPending, payment.Status, "Payment status should be 'pending'")
	})

	t.Run("fail to authorize payment due to duplicate external reference", func(t *testing.T) {
		input1 := dto.AuthorizePaymentInput{
			Amount:            1000,
			ExternalReference: uuid.New().String(),
			PaymentMethod:     "credit_card",
		}

		payment1, err := ap.Run(ctx, input1)
		require.NoError(t, err, "First payment insertion should not return error")
		assert.NotNil(t, payment1)

		input2 := dto.AuthorizePaymentInput{
			Amount:            2000,
			ExternalReference: payment1.ExternalReference,
			PaymentMethod:     "credit_card",
		}

		payment2, err := ap.Run(ctx, input2)
		require.Error(t, err, "Run should return error due to duplicate external reference")
		assert.Nil(t, payment2, "Second payment should be nil due to duplicate external reference")
	})

	t.Run("fail to authorize payment due to context timeout", func(t *testing.T) {
		timeoutCtx, cancel := context.WithTimeout(ctx, 0)
		defer cancel()

		input := dto.AuthorizePaymentInput{
			Amount:            1000,
			ExternalReference: uuid.New().String(),
			PaymentMethod:     "credit_card",
		}

		payment, err := ap.Run(timeoutCtx, input)
		require.Error(t, err, "Run should return error due to context timeout")
		assert.Nil(t, payment, "Payment should be nil due to context timeout")
	})
}
