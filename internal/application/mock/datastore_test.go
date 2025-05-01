package mock_test

import (
	"context"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/mock"

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
