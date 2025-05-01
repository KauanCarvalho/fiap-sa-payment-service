package handler

import (
	"net/http"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain"

	"github.com/gin-gonic/gin"
)

type HealthCheckHandler struct {
	datastore domain.Datastore
}

func NewHealthCheckHandler(datastore domain.Datastore) *HealthCheckHandler {
	return &HealthCheckHandler{
		datastore: datastore,
	}
}

// Ping checks the health of the application (connection to database).
// @Summary	    Health check
// @Description Checks the health of the application (connection to database)
// @Tags        HealthCheck
// @Accept      json
// @Produce     json
// @Success     200 {object} dto.HealthCheckOutput
// @Failure     500 {object} dto.APIErrorsOutput
// @Router      /healthcheck [get].
func (hch *HealthCheckHandler) Ping(c *gin.Context) {
	var apiErrors dto.APIErrorsOutput

	if err := hch.datastore.Ping(c.Request.Context()); err != nil {
		apiErrors.Errors = append(apiErrors.Errors, dto.APIErrorOutput{
			Details: "Database health check failed",
			Field:   "database",
			Message: err.Error(),
		})
	}

	if len(apiErrors.Errors) > 0 {
		c.JSON(http.StatusInternalServerError, apiErrors)
		return
	}

	c.JSON(http.StatusOK, dto.HealthCheckOutput{Status: "ok"})
}
