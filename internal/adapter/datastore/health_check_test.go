package datastore_test

import (
	"context"
	"errors"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/mock"
	"github.com/test-go/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	t.Run("successfully ping the database", func(t *testing.T) {
		assert.NoError(t, ds.Ping(ctx), "The Ping method should not return an error when the database is reachable")
	})

	t.Run("fail to ping the database", func(t *testing.T) {
		mockDatastore := &mock.DatastoreMock{
			PingFn: func(_ context.Context) error {
				return errors.New("database error")
			},
		}

		require.Error(t, mockDatastore.Ping(ctx),
			"The Ping method should return an error when the database is unreachable",
		)
	})
}
