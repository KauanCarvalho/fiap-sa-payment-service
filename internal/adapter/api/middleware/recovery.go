package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/dto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				req := c.Request

				zap.L().Error("Panic recovered",
					zap.Any("error", rec),
					zap.String("route", c.FullPath()),
					zap.String("method", req.Method),
					zap.String("client_ip", c.ClientIP()),
					zap.String("user_agent", req.UserAgent()),
					zap.ByteString("stacktrace", debug.Stack()),
				)

				errorOutput := dto.SimpleAPIErrorsOutput(
					"A unexpected internal error occurred",
					"",
					"internal server error",
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, errorOutput)
			}
		}()

		c.Next()
	}
}
