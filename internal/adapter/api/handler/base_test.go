package handler_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/adapter/api"
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

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	ctx                 context.Context
	cfg                 *config.Config
	mongoDB             *mongo.Client
	ds                  domain.Datastore
	authorizePayment    usecase.AuthorizePaymentUseCase
	updatePaymentStatus usecase.UpdatePaymentStatusUseCase
	ginEngine           *gin.Engine
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	logger := zap.NewNop()
	zap.ReplaceGlobals(logger)

	ctx = context.Background()
	cfg = config.Load()

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

	var err error
	mongoDB, err = di.NewMongoConnection(cfg)
	if err != nil {
		log.Fatalf("error when initializing database connection: %v", err)
	}

	ds = datastore.NewDatastore(mongoDB, cfg.MongoDatabaseName)
	authorizePayment = usecase.NewAuthorizePaymentUseCase(ds)
	updatePaymentStatus = usecase.NewUpdatePaymentStatusUseCase(ds, mockSNS)
	ginEngine = api.GenerateRouter(cfg, ds, authorizePayment, updatePaymentStatus)

	os.Exit(m.Run())
}
