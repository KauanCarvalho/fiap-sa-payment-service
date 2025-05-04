package usecase_test

import (
	"context"
	"errors"
	"testing"

	appmock "github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/mock"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase/dto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/test-go/testify/mock"
)

func TestUpdatePaymentStatusUseCase(t *testing.T) {
	t.Run("successfully update payment status", func(t *testing.T) {
		input := dto.AuthorizePaymentInput{
			Amount:            1500,
			ExternalReference: uuid.New().String(),
			PaymentMethod:     "pix",
		}

		payment, err := ap.Run(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, payment)

		updateInput := dto.UpdatePaymentStatusInput{
			ExternalReference: payment.ExternalReference,
			Status:            "completed",
		}

		_, err = up.Run(ctx, updateInput)
		require.NoError(t, err)

		found, err := ds.FindPaymentByExternalReference(ctx, payment.ExternalReference)
		require.NoError(t, err)
		require.NotNil(t, found)

		assert.Equal(t, "completed", found.Status)
	})

	t.Run("fail to update non-existent payment", func(t *testing.T) {
		updateInput := dto.UpdatePaymentStatusInput{
			ExternalReference: uuid.New().String(),
			Status:            "completed",
		}

		_, err := up.Run(ctx, updateInput)
		require.Error(t, err)
	})

	t.Run("fail to update with invalid context", func(t *testing.T) {
		timeoutCtx, cancel := context.WithTimeout(ctx, 0)
		defer cancel()

		updateInput := dto.UpdatePaymentStatusInput{
			ExternalReference: uuid.New().String(),
			Status:            "completed",
		}

		_, err := up.Run(timeoutCtx, updateInput)
		require.Error(t, err)
	})

	t.Run("should fail when SNS ListTopicsWithContext returns error", func(t *testing.T) {
		mockSNS := new(appmock.SNSClient)

		mockSNS.On("ListTopicsWithContext", mock.Anything, mock.Anything).
			Return((*sns.ListTopicsOutput)(nil), errors.New("SNS ListTopics error"))

		upWithFailingSNS := usecase.NewUpdatePaymentStatusUseCase(ds, mockSNS)

		input := dto.AuthorizePaymentInput{
			Amount:            1000,
			ExternalReference: uuid.New().String(),
			PaymentMethod:     "pix",
		}

		payment, err := ap.Run(ctx, input)
		require.NoError(t, err)

		updateInput := dto.UpdatePaymentStatusInput{
			ExternalReference: payment.ExternalReference,
			Status:            "failed",
		}

		_, err = upWithFailingSNS.Run(ctx, updateInput)
		require.ErrorContains(t, err, "failed to get SNS topic ARN")
	})

	t.Run("should fail when SNS PublishWithContext returns error", func(t *testing.T) {
		mockSNS := new(appmock.SNSClient)
		topicARN := "arn:aws:sns:us-east-1:123456789012:fiap_sa_payment_service_payment_events"

		mockSNS.On("ListTopicsWithContext", mock.Anything, mock.Anything).
			Return(&sns.ListTopicsOutput{
				Topics: []*sns.Topic{
					{TopicArn: aws.String(topicARN)},
				},
			}, nil)

		mockSNS.On("PublishWithContext", mock.Anything, mock.Anything).
			Return((*sns.PublishOutput)(nil), errors.New("SNS Publish error"))

		upWithFailingPublish := usecase.NewUpdatePaymentStatusUseCase(ds, mockSNS)

		input := dto.AuthorizePaymentInput{
			Amount:            1000,
			ExternalReference: uuid.New().String(),
			PaymentMethod:     "pix",
		}

		payment, err := ap.Run(ctx, input)
		require.NoError(t, err)

		updateInput := dto.UpdatePaymentStatusInput{
			ExternalReference: payment.ExternalReference,
			Status:            "failed",
		}

		_, err = upWithFailingPublish.Run(ctx, updateInput)
		require.ErrorContains(t, err, "failed to send SNS notification")
	})
}
