package handler

import (
	"errors"
	"net/http"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase"
	useCaseDTO "github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase/dto"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase/mappers"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/shared/validation"

	"github.com/gin-gonic/gin"
)

type PaymentHandler interface {
	Authorize(c *gin.Context)
	UpdateStatus(c *gin.Context)
}

type paymentHandler struct {
	authorizePaymentUseCase usecase.AuthorizePaymentUseCase
	updatePaymentUseCase    usecase.UpdatePaymentStatusUseCase
}

func NewPaymentHandler(authorizePaymentUseCase usecase.AuthorizePaymentUseCase, updatePaymentUseCase usecase.UpdatePaymentStatusUseCase) PaymentHandler {
	return &paymentHandler{
		authorizePaymentUseCase: authorizePaymentUseCase,
		updatePaymentUseCase:    updatePaymentUseCase,
	}
}

// Authorize a payment.
// @Summary	    Authorize a payment.
// @Description Authorize a payment.
// @Tags        Payment
// @Accept      json
// @Produce     json
// @Param       payment body useCaseDTO.AuthorizePaymentInput true "request body"
// @Success     201 {object} dto.PaymentOutput
// @Failure     400 {object} dto.APIErrorsOutput
// @Failure     409 {object} dto.APIErrorsOutput
// @Failure     500 {object} dto.APIErrorsOutput
// @Router      /api/v1/payments/authorize [post].
func (h *paymentHandler) Authorize(c *gin.Context) {
	ctx := c.Request.Context()

	var input useCaseDTO.AuthorizePaymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput(
			"",
			"body",
			"Invalid request body",
		))
		return
	}

	if err := useCaseDTO.ValidatePaymentCreate(input); err != nil {
		errors := validation.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, dto.ErrorsFromValidationErrors(errors))
		return
	}

	payment, err := h.authorizePaymentUseCase.Run(ctx, input)
	if err != nil {
		if errors.Is(err, datastore.ErrDuplicateExternalReference) {
			c.JSON(http.StatusConflict, dto.SimpleAPIErrorsOutput(
				"",
				"external_reference",
				"External reference already exists",
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.SimpleAPIErrorsOutput("", "", "failed to create payment"))
		return
	}

	c.JSON(http.StatusCreated, mappers.ToPaymentDTO(*payment))
}

// UpdateStatus updates the status of a payment.
// @Summary	    Update payment status.
// @Description Update the status of a payment by external reference.
// @Tags        Payment
// @Accept      json
// @Produce     json
// @Param       update body useCaseDTO.UpdatePaymentStatusInput true "request body"
// @Param       external_reference path string true "external reference"
// @Success     204
// @Failure     400 {object} dto.APIErrorsOutput
// @Failure     404 {object} dto.APIErrorsOutput
// @Failure     500 {object} dto.APIErrorsOutput
// @Router      /api/v1/payments/{exteral_reference}/update-status [patch].
func (h *paymentHandler) UpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()

	externalRef := c.Param("external_reference")
	if externalRef == "" {
		c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput("", "external_reference", "Missing external reference in path"))
		return
	}

	var input useCaseDTO.UpdatePaymentStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput(
			"",
			"body",
			"Invalid request body",
		))
		return
	}

	input.ExternalReference = externalRef

	if err := useCaseDTO.ValidatePaymentStatusUpdate(input); err != nil {
		errors := validation.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, dto.ErrorsFromValidationErrors(errors))
		return
	}

	_, err := h.updatePaymentUseCase.Run(ctx, input)
	if err != nil {
		if errors.Is(err, datastore.ErrPaymentNotFound) {
			c.JSON(http.StatusNotFound, dto.SimpleAPIErrorsOutput(
				"",
				"external_reference",
				"Payment not found",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.SimpleAPIErrorsOutput(
			"",
			"",
			"Failed to update payment status",
		))
		return
	}

	c.Status(http.StatusNoContent)
}
