package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoggerMiddleware(t *testing.T) {
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
	router.Use(middleware.Logger())

	router.GET("/ping", func(c *gin.Context) {
		time.Sleep(5 * time.Millisecond)
		c.String(http.StatusOK, "pong")
	})

	t.Run("logs request and returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ping?foo=bar", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "pong", w.Body.String())

		log := logBuffer.String()

		t.Run("includes method", func(t *testing.T) {
			assert.Contains(t, log, `"method":"GET"`)
		})

		t.Run("includes path", func(t *testing.T) {
			assert.Contains(t, log, `"path":"/ping"`)
		})

		t.Run("includes query", func(t *testing.T) {
			assert.Contains(t, log, `"query":"foo=bar"`)
		})

		t.Run("includes status", func(t *testing.T) {
			assert.Contains(t, log, `"status":200`)
		})

		t.Run("includes latency", func(t *testing.T) {
			assert.Contains(t, log, `"latency":`)
		})

		t.Run("includes client_ip", func(t *testing.T) {
			assert.Contains(t, log, `"client_ip":"`)
		})
	})
}
