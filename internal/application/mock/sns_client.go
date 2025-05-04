package mock

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/stretchr/testify/mock"
)

type SNSClient struct {
	mock.Mock
}

func (m *SNSClient) ListTopicsWithContext(ctx aws.Context, input *sns.ListTopicsInput, _ ...request.Option) (*sns.ListTopicsOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*sns.ListTopicsOutput), args.Error(1) //nolint:errcheck // ignore error check.
}

func (m *SNSClient) PublishWithContext(ctx aws.Context, input *sns.PublishInput, _ ...request.Option) (*sns.PublishOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*sns.PublishOutput), args.Error(1) //nolint:errcheck // ignore error check.
}
