package validation_test

import (
	"errors"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/shared/validation"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Product struct {
	Name  string `validate:"required"`
	Price int    `validate:"gt=0"`
}

func TestHandleValidationError(t *testing.T) {
	t.Run("validates required field error", func(t *testing.T) {
		validate := validator.New()
		product := Product{
			Name:  "",
			Price: 10,
		}

		err := validate.Struct(product)

		validationErrors := validation.HandleValidationError(err)

		require.Len(t, validationErrors, 1)
		assert.Equal(t, "Name", validationErrors[0].Field)
		assert.Equal(t, "This field is required", validationErrors[0].Message)
	})

	t.Run("validates gt error", func(t *testing.T) {
		validate := validator.New()
		product := Product{
			Price: -1,
			Name:  "Product Name",
		}

		err := validate.Struct(product)

		validationErrors := validation.HandleValidationError(err)

		require.Len(t, validationErrors, 1)
		assert.Equal(t, "Price", validationErrors[0].Field)
		assert.Equal(t, "This field must be greater than 0", validationErrors[0].Message)
	})

	t.Run("validates generic error", func(t *testing.T) {
		err := validation.HandleValidationError(errors.New("Something went wrong"))

		require.Len(t, err, 1)
		assert.Equal(t, "generic error response", err[0].Field)
		assert.Equal(t, "Something went wrong", err[0].Message)
	})

	t.Run("validates multiple errors", func(t *testing.T) {
		validate := validator.New()
		product := Product{
			Name:  "",
			Price: -1,
		}

		err := validate.Struct(product)

		validationErrors := validation.HandleValidationError(err)

		require.Len(t, validationErrors, 2)

		assert.Equal(t, "Name", validationErrors[0].Field)
		assert.Equal(t, "This field is required", validationErrors[0].Message)

		assert.Equal(t, "Price", validationErrors[1].Field)
		assert.Equal(t, "This field must be greater than 0", validationErrors[1].Message)
	})

	t.Run("validates missing parameter in gt", func(t *testing.T) {
		validate := validator.New()
		product := Product{
			Price: 0,
			Name:  "Product Name",
		}

		err := validate.Struct(product)

		validationErrors := validation.HandleValidationError(err)

		require.Len(t, validationErrors, 1)
		assert.Equal(t, "Price", validationErrors[0].Field)
		assert.Equal(t, "This field must be greater than 0", validationErrors[0].Message)
	})

	t.Run("validates latitude error", func(t *testing.T) {
		type Product struct {
			Latitude float64 `validate:"latitude"`
		}

		product := Product{
			Latitude: -100,
		}

		validate := validator.New()

		err := validate.Struct(product)

		validationErrors := validation.HandleValidationError(err)

		require.Len(t, validationErrors, 1)
		assert.Equal(t, "Latitude", validationErrors[0].Field)
		assert.NotEmpty(t, validationErrors[0].Message)
	})
}
