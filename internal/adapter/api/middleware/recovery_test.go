package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/api/middleware"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/dto"
	"github.com/test-go/testify/require"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestRecoveryMiddleware(t *testing.T) {
	var logBuffer bytes.Buffer

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(&logBuffer),
		zap.DebugLevel,
	)

	logger := zap.New(core)
	zap.ReplaceGlobals(logger)

	router := gin.New()
	router.Use(middleware.Recovery())

	router.GET("/panic", func(_ *gin.Context) {
		panic("simulated panic")
	})

	t.Run("recovers from panic and returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		t.Run("returns status 500", func(t *testing.T) {
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("returns valid JSON body", func(t *testing.T) {
			var output dto.APIErrorsOutput
			err := json.Unmarshal(w.Body.Bytes(), &output)
			require.NoError(t, err)
			assert.Len(t, output.Errors, 1)
			assert.Equal(t, "internal server error", output.Errors[0].Message)
			assert.Equal(t, "A unexpected internal error occurred", output.Errors[0].Details)
		})

		log := logBuffer.String()

		t.Run("logs recovered panic", func(t *testing.T) {
			assert.Contains(t, log, `"msg":"Panic recovered"`)
			assert.Contains(t, log, `"error":"simulated panic"`)
			assert.Contains(t, log, `"stacktrace":"`)
		})
	})
}
