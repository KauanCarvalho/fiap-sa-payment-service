package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPaymentHandler_Authorize(t *testing.T) {
	t.Run("should authorize payment successfully", func(t *testing.T) {
		body := map[string]interface{}{
			"amount":             50.0,
			"external_reference": uuid.New().String(),
			"payment_method":     "pix",
		}
		payload, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/payments/authorize", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/payments/authorize", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request body")
	})

	t.Run("should return 400 for validation error", func(t *testing.T) {
		body := map[string]interface{}{
			"amount":             -10.0,
			"external_reference": "",
			"payment_method":     "",
		}
		payload, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/payments/authorize", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "errors")
	})

	t.Run("should return 409 for duplicate external_reference", func(t *testing.T) {
		body := map[string]interface{}{
			"amount":             99.0,
			"external_reference": uuid.New().String(),
			"payment_method":     "pix",
		}
		payload, _ := json.Marshal(body)

		req1, _ := http.NewRequest(http.MethodPost, "/api/v1/payments/authorize", bytes.NewReader(payload))
		req1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		ginEngine.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Code)

		req2, _ := http.NewRequest(http.MethodPost, "/api/v1/payments/authorize", bytes.NewReader(payload))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		ginEngine.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusConflict, w2.Code)
		assert.Contains(t, w2.Body.String(), "External reference already exists")
	})
}
