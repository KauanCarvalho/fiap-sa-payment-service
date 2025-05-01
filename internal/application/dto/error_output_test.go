package dto_test

import (
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/dto"

	"github.com/stretchr/testify/assert"
)

func TestSimpleAPIErrorsOutput(t *testing.T) {
	details := "Error details"
	field := "field"
	message := "Error message"

	t.Run("should return APIErrorsOutput with one error", func(t *testing.T) {
		result := dto.SimpleAPIErrorsOutput(details, field, message)

		t.Run("should contain exactly one error", func(t *testing.T) {
			assert.Len(t, result.Errors, 1, "Must have only one error")
		})

		t.Run("should match error details", func(t *testing.T) {
			assert.Equal(t, details, result.Errors[0].Details, "Details not match")
		})

		t.Run("should match error field", func(t *testing.T) {
			assert.Equal(t, field, result.Errors[0].Field, "Field not match")
		})

		t.Run("should match error message", func(t *testing.T) {
			assert.Equal(t, message, result.Errors[0].Message, "Message not match")
		})
	})
}
