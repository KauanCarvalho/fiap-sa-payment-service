package worker

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.uber.org/zap"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/config"
)

type ProcessingMessage struct {
	Message   *sqs.Message
	QueueName string
	QueueURL  *string
}

type SNSEnvelope struct {
	Body string `json:"body"`
}

type Processor interface {
	Process(context.Context, ProcessingMessage) error
}

func Consumer(
	ctx context.Context,
	awsConfig config.AWSConfig,
	sqsQueue config.SQSQueue,
	chn chan<- ProcessingMessage,
	wgConsumer *sync.WaitGroup,
) {
	defer wgConsumer.Done()

	zap.L().Info("Starting consumer", zap.String("queueName", sqsQueue.FullName))

	queueOut, err := awsConfig.SQSConfig.Client.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: &sqsQueue.FullName})
	if err != nil {
		zap.L().Panic("Error getting queue URL", zap.String("queueName", sqsQueue.FullName), zap.Error(err))
	}

	for {
		select {
		case <-ctx.Done():
			zap.L().Panic("Context canceled, stopping consumer", zap.String("queueName", sqsQueue.FullName))
			return
		default:
			receiveParams := &sqs.ReceiveMessageInput{
				QueueUrl:        queueOut.QueueUrl,
				WaitTimeSeconds: aws.Int64(awsConfig.SQSConfig.WaitTime),
			}

			var result *sqs.ReceiveMessageOutput
			result, err = awsConfig.SQSConfig.Client.ReceiveMessage(receiveParams)
			if err != nil {
				zap.L().Error("Error receiving message", zap.String("queueName", sqsQueue.FullName), zap.Error(err))
				continue
			}

			for _, message := range result.Messages {
				chn <- buildProcessingMessage(message, sqsQueue.Name, queueOut.QueueUrl)
			}
		}
	}
}

func buildProcessingMessage(message *sqs.Message, queueName string, queueURL *string) ProcessingMessage {
	return ProcessingMessage{
		Message:   message,
		QueueName: queueName,
		QueueURL:  queueURL,
	}
}

func DeleteMessage(awsConfig config.AWSConfig, processingMessage ProcessingMessage) error {
	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      processingMessage.QueueURL,
		ReceiptHandle: processingMessage.Message.ReceiptHandle,
	}

	_, err := awsConfig.SQSConfig.Client.DeleteMessage(deleteParams)
	if err != nil {
		zap.L().Error("Error deleting message", zap.String("queueName", processingMessage.QueueName), zap.Error(err))
		return err
	}

	return nil
}
