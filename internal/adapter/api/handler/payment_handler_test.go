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

func TestPaymentHandler_UpdateStatus(t *testing.T) {
	t.Run("should update payment status successfully", func(t *testing.T) {
		externalRef := uuid.New().String()

		createBody := map[string]interface{}{
			"amount":             100.0,
			"external_reference": externalRef,
			"payment_method":     "pix",
		}
		createPayload, _ := json.Marshal(createBody)
		reqCreate, _ := http.NewRequest(http.MethodPost, "/api/v1/payments/authorize", bytes.NewReader(createPayload))
		reqCreate.Header.Set("Content-Type", "application/json")
		wCreate := httptest.NewRecorder()
		ginEngine.ServeHTTP(wCreate, reqCreate)
		assert.Equal(t, http.StatusCreated, wCreate.Code)

		updateBody := map[string]interface{}{
			"status": "completed",
		}
		updatePayload, _ := json.Marshal(updateBody)
		url := "/api/v1/payments/" + externalRef + "/update-status"
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(updatePayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("should return 400 for missing path param", func(t *testing.T) {
		updateBody := map[string]interface{}{
			"status": "completed",
		}
		updatePayload, _ := json.Marshal(updateBody)
		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/payments//update-status", bytes.NewReader(updatePayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Missing external reference")
	})

	t.Run("should return 400 for invalid body", func(t *testing.T) {
		externalRef := uuid.New().String()
		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/payments/"+externalRef+"/update-status", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request body")
	})

	t.Run("should return 404 if payment not found", func(t *testing.T) {
		nonExistent := uuid.New().String()
		updateBody := map[string]interface{}{
			"status": "completed",
		}
		updatePayload, _ := json.Marshal(updateBody)
		url := "/api/v1/payments/" + nonExistent + "/update-status"
		req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(updatePayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		ginEngine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Payment not found")
	})
}
