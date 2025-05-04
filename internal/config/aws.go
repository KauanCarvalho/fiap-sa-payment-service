package config

import (
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SNSConfig struct {
	Endpoint string
}

type AWSConfig struct {
	AccessKey string
	SecretKey string
	Region    string
	SNSConfig SNSConfig
}

func LoadAWSConfig() *AWSConfig {
	return &AWSConfig{
		AccessKey: fetchEnv("AWS_ACCESS_KEY_ID"),
		SecretKey: fetchEnv("AWS_SECRET_ACCESS_KEY"),
		Region:    getEnv("AWS_REGION", "us-east-1"),
		SNSConfig: SNSConfig{
			Endpoint: getEnv("AWS_SNS_ENDPOINT", ""),
		},
	}
}

func NewSNSClient(awsConfig *AWSConfig) (domain.SNSClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(awsConfig.Region),
		Endpoint: aws.String(awsConfig.SNSConfig.Endpoint),
	})
	if err != nil {
		return nil, err
	}

	return sns.New(sess), nil
}
