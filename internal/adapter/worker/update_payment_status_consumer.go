package worker

import (
	"context"
	"encoding/json"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase"
	useCaseDTO "github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase/dto"
	"go.uber.org/zap"
)

type UpdatePaymentStatusConsumer struct {
	updatePaymentUseCase usecase.UpdatePaymentStatusUseCase
}

func NewUpdatePaymentStatusConsumer(updatePaymentUseCase usecase.UpdatePaymentStatusUseCase) *UpdatePaymentStatusConsumer {
	return &UpdatePaymentStatusConsumer{
		updatePaymentUseCase: updatePaymentUseCase,
	}
}

func (consumer UpdatePaymentStatusConsumer) Process(_ context.Context, processingMessage ProcessingMessage) error {
	parsedBody := dto.UpdatePaymentStatusInput{}

	envlope := SNSEnvelope{}
	err := json.Unmarshal([]byte(*processingMessage.Message.Body), &envlope)
	if err != nil {
		zap.L().Error(
			"error parsing sns envelope",
			zap.String("queueName", processingMessage.QueueName),
			zap.Error(err),
		)
		return err
	}

	err = json.Unmarshal([]byte(envlope.Body), &parsedBody)
	if err != nil {
		zap.L().Error(
			"error parsing message body",
			zap.String("queueName", processingMessage.QueueName),
			zap.Error(err),
		)
		return err
	}

	ctx := context.Background()

	payload := useCaseDTO.UpdatePaymentStatusInput{
		ExternalReference: parsedBody.ExternalRef,
		Status:            parsedBody.Status,
	}

	err = useCaseDTO.ValidatePaymentStatusUpdate(payload)
	if err != nil {
		zap.L().Error(
			"error validating payment status update",
			zap.String("queueName", processingMessage.QueueName),
			zap.Error(err),
		)
		return err
	}

	_, err = consumer.updatePaymentUseCase.Run(ctx, payload)
	if err != nil {
		zap.L().Error(
			"error updating payment status",
			zap.String("queueName", processingMessage.QueueName),
			zap.Error(err),
		)
		return err
	}

	zap.L().Info(
		"payment status updated successfully",
		zap.String("queueName", processingMessage.QueueName),
	)

	return nil
}
