package worker_test

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/worker"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/dto"
	useCaseDTO "github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase/dto"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdatePaymentStatusConsumer(t *testing.T) {
	consumer := worker.NewUpdatePaymentStatusConsumer(up)

	t.Run("successfully process message", func(t *testing.T) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		num := r.Intn(999_999_999) + 1

		input := useCaseDTO.AuthorizePaymentInput{
			Amount:            1500,
			ExternalReference: strconv.Itoa(num),
			PaymentMethod:     "pix",
		}

		payment, err := ap.Run(ctx, input)
		require.NoError(t, err)

		updatePayload := dto.UpdatePaymentStatusInput{
			ExternalRef: payment.ExternalReference,
			Status:      "completed",
		}

		payloadBytes, err := json.Marshal(updatePayload)
		require.NoError(t, err)

		snsEnvelope := worker.SNSEnvelope{
			Body: string(payloadBytes),
		}
		snsEnvelopeBytes, err := json.Marshal(snsEnvelope)
		require.NoError(t, err)

		body := string(snsEnvelopeBytes)
		msg := worker.ProcessingMessage{
			QueueName: "test-queue",
			Message:   &sqs.Message{Body: &body},
		}

		err = consumer.Process(ctx, msg)
		require.NoError(t, err)

		updated, err := ds.FindPaymentByExternalReference(ctx, payment.ExternalReference)
		require.NoError(t, err)
		assert.Equal(t, "completed", updated.Status)
	})

	t.Run("whe payment does not exist", func(t *testing.T) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		num := r.Intn(999_999_999) + 1

		updatePayload := dto.UpdatePaymentStatusInput{
			ExternalRef: strconv.Itoa(num),
			Status:      "completed",
		}

		payloadBytes, err := json.Marshal(updatePayload)
		require.NoError(t, err)

		snsEnvelope := worker.SNSEnvelope{
			Body: string(payloadBytes),
		}
		snsEnvelopeBytes, err := json.Marshal(snsEnvelope)
		require.NoError(t, err)

		body := string(snsEnvelopeBytes)
		msg := worker.ProcessingMessage{
			QueueName: "test-queue",
			Message:   &sqs.Message{Body: &body},
		}

		err = consumer.Process(ctx, msg)
		require.NoError(t, err)
	})

	t.Run("invalid JSON in envelope", func(t *testing.T) {
		body := "{ invalid json }"

		msg := worker.ProcessingMessage{
			QueueName: "test-queue",
			Message:   &sqs.Message{Body: &body},
		}

		err := consumer.Process(ctx, msg)
		require.Error(t, err)
	})

	t.Run("invalid JSON in SNS message", func(t *testing.T) {
		invalidEnvelope := worker.SNSEnvelope{
			Body: "{{ invalid }}",
		}
		body, _ := json.Marshal(invalidEnvelope)
		bodyStr := string(body)

		msg := worker.ProcessingMessage{
			QueueName: "test-queue",
			Message:   &sqs.Message{Body: &bodyStr},
		}

		err := consumer.Process(ctx, msg)
		require.Error(t, err)
	})

	t.Run("invalid status", func(t *testing.T) {
		input := useCaseDTO.AuthorizePaymentInput{
			Amount:            900,
			ExternalReference: uuid.New().String(),
			PaymentMethod:     "pix",
		}
		payment, err := ap.Run(ctx, input)
		require.NoError(t, err)

		invalidPayload := dto.UpdatePaymentStatusInput{
			ExternalRef: payment.ExternalReference,
			Status:      "invalid-status",
		}
		payloadBytes, _ := json.Marshal(invalidPayload)

		env := worker.SNSEnvelope{
			Body: string(payloadBytes),
		}
		envBytes, _ := json.Marshal(env)
		body := string(envBytes)

		msg := worker.ProcessingMessage{
			QueueName: "test-queue",
			Message:   &sqs.Message{Body: &body},
		}

		err = consumer.Process(ctx, msg)
		require.Error(t, err)
	})
}
