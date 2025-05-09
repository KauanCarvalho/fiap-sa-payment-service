package config

import (
	"github.com/KauanCarvalho/fiap-sa-payment-service/internal/core/domain"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SNSConfig struct {
	Endpoint string
}

type SQSQueue struct {
	Name     string
	FullName string
}

type SQSConfig struct {
	WaitTime   int64
	BaseURL    string
	Queues     []SQSQueue
	Session    *session.Session
	Client     *sqs.SQS
	NumWorkers int
}

type AWSConfig struct {
	AccessKey    string
	SecretKey    string
	SessionToken string
	Region       string
	SNSConfig    SNSConfig
	SQSConfig    SQSConfig
}

const defaultSQSWaitTime = int64(5)
const defaultSQSNumWorkers = 5

func LoadAWSConfig(projectCfg *Config) *AWSConfig {
	return &AWSConfig{
		AccessKey:    fetchEnv("AWS_ACCESS_KEY_ID"),
		SecretKey:    fetchEnv("AWS_SECRET_ACCESS_KEY"),
		SessionToken: fetchEnv("AWS_SESSION_TOKEN"),
		Region:       getEnv("AWS_REGION", "us-east-1"),
		SNSConfig: SNSConfig{
			Endpoint: getEnv("AWS_SNS_ENDPOINT", ""),
		},
		SQSConfig: SQSConfig{
			WaitTime:   getEnvAsInt64("AWS_SQS_WAIT_TIME", defaultSQSWaitTime),
			BaseURL:    getEnv("AWS_BASE_URL", ""),
			Queues:     projectCfg.loadQueues(),
			NumWorkers: getEnvAsInt("AWS_SQS_NUM_WORKERS", defaultSQSNumWorkers),
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

func (appConfig Config) loadQueues() []SQSQueue {
	loadedSQSQueues := make([]SQSQueue, 0, len(sqsQueueNames()))

	for _, queueName := range sqsQueueNames() {
		sqsQueue := SQSQueue{
			Name:     queueName,
			FullName: appConfig.AppName + "_" + queueName,
		}

		loadedSQSQueues = append(loadedSQSQueues, sqsQueue)
	}

	return loadedSQSQueues
}

func InstantiateSQSClient(awsConfig *AWSConfig) {
	awsSession := session.Must(
		session.NewSession(&aws.Config{
			Endpoint: aws.String(awsConfig.SQSConfig.BaseURL),
			Region:   aws.String(awsConfig.Region)},
		),
	)

	awsConfig.SQSConfig.Session = awsSession
	awsConfig.SQSConfig.Client = sqs.New(awsConfig.SQSConfig.Session)
}

func sqsQueueNames() []string {
	return []string{"webhook_events"}
}
