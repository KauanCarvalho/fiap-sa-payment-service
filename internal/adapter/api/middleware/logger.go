package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		raw := c.Request.URL.RawQuery

		c.Next()

		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		status := c.Writer.Status()
		duration := time.Since(start)

		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", raw),
			zap.Int("status", status),
			zap.String("client_ip", clientIP),
			zap.Duration("latency", duration),
		}

		zap.L().Info("request", fields...)
	}
}
