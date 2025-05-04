package mock_test

import (
	"context"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/mock"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain/entities"

	"github.com/stretchr/testify/require"
)

func TestDatastoreMock_Ping(t *testing.T) {
	t.Run("when PingFn is defined, it returns nil", func(t *testing.T) {
		ds := &mock.DatastoreMock{
			PingFn: func(_ context.Context) error {
				return nil
			},
		}

		err := ds.Ping(ctx)
		require.NoError(t, err)
	})

	t.Run("when PingFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		err := ds.Ping(ctx)
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
	})
}

func TestDatastoreMock_AuthorizePayment(t *testing.T) {
	t.Run("when AuthorizePaymentFn is defined, it returns nil", func(t *testing.T) {
		ds := &mock.DatastoreMock{
			CreatePaymentFn: func(_ context.Context, _ *entities.Payment) error {
				return nil
			},
		}

		err := ds.CreatePayment(ctx, &entities.Payment{})
		require.NoError(t, err)
	})

	t.Run("when CreatePaymentFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		err := ds.CreatePayment(ctx, &entities.Payment{})
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
	})
}

func TestDatastoreMock_FindPaymentByExternalReference(t *testing.T) {
	t.Run("when FindPaymentByExternalReferenceFn is defined, it returns payment", func(t *testing.T) {
		ds := &mock.DatastoreMock{
			FindPaymentByExternalReferenceFn: func(_ context.Context, _ string) (*entities.Payment, error) {
				return &entities.Payment{}, nil
			},
		}

		payment, err := ds.FindPaymentByExternalReference(ctx, "external_reference")
		require.NoError(t, err)
		require.NotNil(t, payment)
	})

	t.Run("when FindPaymentByExternalReferenceFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		payment, err := ds.FindPaymentByExternalReference(ctx, "external_reference")
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
		require.Nil(t, payment)
	})
}

func TestDatastoreMock_UpdatePayment(t *testing.T) {
	t.Run("when UpdatePaymentFn is defined, it returns nil", func(t *testing.T) {
		ds := &mock.DatastoreMock{
			UpdatePaymentFn: func(_ context.Context, _ *entities.Payment) error {
				return nil
			},
		}

		err := ds.UpdatePayment(ctx, &entities.Payment{})
		require.NoError(t, err)
	})

	t.Run("when UpdatePaymentFn is not defined, it returns ErrFunctionNotImplemented", func(t *testing.T) {
		ds := &mock.DatastoreMock{}

		err := ds.UpdatePayment(ctx, &entities.Payment{})
		require.ErrorIs(t, err, mock.ErrFunctionNotImplemented)
	})
}
