package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain/entities"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase/dto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

const updatePaymentTopicName = "fiap_sa_payment_service_payment_events"

type UpdatePaymentStatusUseCase interface {
	Run(ctx context.Context, input dto.UpdatePaymentStatusInput) (*entities.Payment, error)
}

type updatePaymentStatusUsecase struct {
	ds        domain.Datastore
	snsClient domain.SNSClient
}

func NewUpdatePaymentStatusUseCase(ds domain.Datastore, snsClient domain.SNSClient) UpdatePaymentStatusUseCase {
	return &updatePaymentStatusUsecase{
		ds:        ds,
		snsClient: snsClient,
	}
}

func (u *updatePaymentStatusUsecase) Run(ctx context.Context, input dto.UpdatePaymentStatusInput) (*entities.Payment, error) {
	payment, err := u.ds.FindPaymentByExternalReference(ctx, input.ExternalReference)
	if err != nil {
		return nil, err
	}

	payment.Status = input.Status

	if err = u.ds.UpdatePayment(ctx, payment); err != nil {
		return nil, err
	}

	topicArn, err := u.getTopicArnByName(ctx, updatePaymentTopicName)
	if err != nil {
		return nil, fmt.Errorf("failed to get SNS topic ARN: %w", err)
	}

	if err = u.sendSNSNotification(ctx, topicArn, payment.ExternalReference, payment.Status); err != nil {
		return nil, fmt.Errorf("failed to send SNS notification: %w", err)
	}

	return payment, nil
}

func (u *updatePaymentStatusUsecase) getTopicArnByName(ctx context.Context, topicName string) (string, error) {
	req := &sns.ListTopicsInput{}

	resp, err := u.snsClient.ListTopicsWithContext(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to list SNS topics: %w", err)
	}

	for _, topic := range resp.Topics {
		if strings.HasSuffix(*topic.TopicArn, topicName) {
			return *topic.TopicArn, nil
		}
	}

	return "", fmt.Errorf("SNS topic with name %s not found", topicName)
}

func (u *updatePaymentStatusUsecase) sendSNSNotification(ctx context.Context, topicArn, externalReference, status string) error {
	payload := map[string]string{
		"external_reference": externalReference,
		"status":             status,
	}

	message, _ := json.Marshal(payload)

	_, err := u.snsClient.PublishWithContext(ctx, &sns.PublishInput{
		Message:  aws.String(string(message)),
		TopicArn: aws.String(topicArn),
	})

	return err
}
