package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/api"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/mock"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/healthcheck", nil)
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"ok"`)
	})

	t.Run("failure database", func(t *testing.T) {
		mockDatastore := &mock.DatastoreMock{
			PingFn: func(_ context.Context) error {
				return errors.New("database error")
			},
		}

		engineMocked := api.GenerateRouter(cfg, mockDatastore)

		req, _ := http.NewRequest(http.MethodGet, "/healthcheck", nil)
		w := httptest.NewRecorder()

		engineMocked.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"field":"database"`)
		assert.Contains(t, w.Body.String(), `"message":"database error"`)
	})
}
