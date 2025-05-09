package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/worker"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/di"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	awsConfig := config.LoadAWSConfig(cfg)
	config.InstantiateSQSClient(awsConfig)

	mongoDB, err := di.NewMongoConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("cannot initialize zap logger: %v", err)
	}
	defer logger.Sync() //nolint:errcheck // It is not necessary to check for errors at this moment.

	zap.ReplaceGlobals(logger.With(zap.String("app", cfg.AppName), zap.String("env", cfg.AppEnv)))

	// SNS client.
	snsClient, err := config.NewSNSClient(awsConfig)
	if err != nil {
		zap.L().Fatal("Failed to create SNS client",
			zap.String("aws_region", awsConfig.Region),
			zap.String("aws_endpoint", awsConfig.SQSConfig.Client.Endpoint),
			zap.Error(err),
		)
	}

	// Datastore.
	ds := datastore.NewDatastore(mongoDB, cfg.MongoDatabaseName)

	// Use cases.
	updatePaymentStatusUseCase := usecase.NewUpdatePaymentStatusUseCase(ds, snsClient)

	// Consumers.
	updatePaymentStatusConsumer := worker.NewUpdatePaymentStatusConsumer(updatePaymentStatusUseCase)

	chnProcessingMessages := make(chan worker.ProcessingMessage, awsConfig.SQSConfig.NumWorkers)
	ctx, cancel := context.WithCancel(context.Background())

	var wgConsumer, wgProcessing sync.WaitGroup

	for range awsConfig.SQSConfig.NumWorkers {
		wgProcessing.Add(1)
		go startWorker(ctx, chnProcessingMessages, &wgProcessing, *awsConfig, updatePaymentStatusConsumer)
	}

	for _, sqsQueue := range awsConfig.SQSConfig.Queues {
		wgConsumer.Add(1)
		go worker.Consumer(ctx, *awsConfig, sqsQueue, chnProcessingMessages, &wgConsumer)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		select {
		case signalReceived := <-signalCh:
			zap.L().Info("Received signal, shutting down gracefully",
				zap.String("signal", signalReceived.String()),
			)
			zap.L().Info("Shutting down gracefully...")

			cancel()
		case <-ctx.Done():
			zap.L().Info("Context canceled, shutting down gracefully...")
		}

		wgConsumer.Wait()

		zap.L().Info("All consumers have been finalized. Closing channel...")

		close(chnProcessingMessages)
	}()

	zap.L().Info("Starting consumers...")

	wgProcessing.Wait()

	zap.L().Info("All channel processors have been finalized. Shutting down completely.")
}
func startWorker(
	ctx context.Context,
	chnProcessingMessages <-chan worker.ProcessingMessage,
	wg *sync.WaitGroup,
	awsConfig config.AWSConfig,
	updatePaymentStatusConsumer worker.Processor,
) {
	defer wg.Done()

	for processingMessage := range chnProcessingMessages {
		zap.L().Info("Processing message from queue",
			zap.String("queue_name", processingMessage.QueueName),
		)

		var processingError error

		func() {
			defer func() {
				if r := recover(); r != nil {
					zap.L().Error("Recovered from panic",
						zap.String("queue_name", processingMessage.QueueName),
						zap.Any("error", r),
					)
				}
			}()

			if processingMessage.QueueName == "webhook_events" {
				processingError = updatePaymentStatusConsumer.Process(ctx, processingMessage)
			}
		}()

		if processingError == nil {
			deleteMessageErr := worker.DeleteMessage(awsConfig, processingMessage)
			if deleteMessageErr != nil {
				zap.L().Error("Error deleting message",
					zap.String("queue_name", processingMessage.QueueName),
					zap.Error(deleteMessageErr),
				)
			}
		}
	}
}
