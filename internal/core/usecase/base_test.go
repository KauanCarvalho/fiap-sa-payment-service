package usecase_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/datastore"
	appmock "github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/mock"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/config"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/usecase"
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/di"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/test-go/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ctx     context.Context
	cfg     *config.Config
	mongoDB *mongo.Client
	ap      usecase.AuthorizePaymentUseCase
	up      usecase.UpdatePaymentStatusUseCase
	ds      domain.Datastore
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	cfg = config.Load()

	var err error
	mongoDB, err = di.NewMongoConnection(cfg)
	if err != nil {
		log.Fatalf("error when creating database connection pool: %v", err)
	}

	mockSNS := new(appmock.SNSClient)

	topicARN := "arn:aws:sns:us-east-1:123456789012:fiap_sa_payment_service_payment_events"

	mockSNS.On("ListTopicsWithContext", ctx, mock.Anything).
		Return(&sns.ListTopicsOutput{
			Topics: []*sns.Topic{
				{TopicArn: aws.String(topicARN)},
			},
		}, nil)

	mockSNS.On("PublishWithContext", ctx, mock.Anything).
		Return(&sns.PublishOutput{}, nil)

	ds = datastore.NewDatastore(mongoDB, cfg.MongoDatabaseName)
	ap = usecase.NewAuthorizePaymentUseCase(ds)
	up = usecase.NewUpdatePaymentStatusUseCase(ds, mockSNS)

	os.Exit(m.Run())
}
