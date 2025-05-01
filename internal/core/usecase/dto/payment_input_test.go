package dto_test

import (
	"strings"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase/dto"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePaymentCreate(t *testing.T) {
	t.Run("invalid input - all fields missing", func(t *testing.T) {
		input := dto.AuthorizePaymentInput{}

		err := dto.ValidatePaymentCreate(input)
		require.Error(t, err)

		validationErrors, ok := err.(validator.ValidationErrors)
		require.True(t, ok)

		var fieldNames []string
		for _, e := range validationErrors {
			ns := e.StructNamespace()
			parts := strings.SplitN(ns, ".", 2)
			if len(parts) == 2 {
				fieldNames = append(fieldNames, parts[1])
			} else {
				fieldNames = append(fieldNames, ns)
			}
		}

		assert.Contains(t, fieldNames, "Amount")
		assert.Contains(t, fieldNames, "ExternalReference")
		assert.Contains(t, fieldNames, "PaymentMethod")
	})

	t.Run("invalid input - zero amount", func(t *testing.T) {
		input := dto.AuthorizePaymentInput{
			Amount:            0,
			ExternalReference: "ref123",
			PaymentMethod:     "pix",
		}

		err := dto.ValidatePaymentCreate(input)
		require.Error(t, err)

		validationErrors := err.(validator.ValidationErrors)

		var fieldNames []string
		for _, e := range validationErrors {
			ns := e.StructNamespace()
			parts := strings.SplitN(ns, ".", 2)
			if len(parts) == 2 {
				fieldNames = append(fieldNames, parts[1])
			} else {
				fieldNames = append(fieldNames, ns)
			}
		}

		assert.Contains(t, fieldNames, "Amount")
	})

	t.Run("valid input", func(t *testing.T) {
		input := dto.AuthorizePaymentInput{
			Amount:            1000,
			ExternalReference: "ref123",
			PaymentMethod:     "pix",
		}

		err := dto.ValidatePaymentCreate(input)
		assert.NoError(t, err)
	})
}
