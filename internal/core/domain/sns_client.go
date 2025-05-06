package domain

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SNSClient interface {
	ListTopicsWithContext(ctx aws.Context, input *sns.ListTopicsInput, opts ...request.Option) (*sns.ListTopicsOutput, error)
	PublishWithContext(ctx aws.Context, input *sns.PublishInput, opts ...request.Option) (*sns.PublishOutput, error)
}
