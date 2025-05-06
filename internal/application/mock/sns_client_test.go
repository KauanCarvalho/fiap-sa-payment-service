package mock_test

import (
	"testing"

	sdkmock "github.com/KauanCarvalho/fiap-sa-payment-service/internal/application/mock"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSNSClient_ListTopicsWithContext(t *testing.T) {
	t.Run("should return expected ListTopicsOutput without error", func(t *testing.T) {
		client := &sdkmock.SNSClient{}
		expectedOutput := &sns.ListTopicsOutput{}

		client.On("ListTopicsWithContext", ctx, &sns.ListTopicsInput{}).
			Return(expectedOutput, nil)

		output, err := client.ListTopicsWithContext(ctx, &sns.ListTopicsInput{})

		require.NoError(t, err)
		assert.Equal(t, expectedOutput, output)
		client.AssertExpectations(t)
	})

	t.Run("should return error when mock is configured to fail", func(t *testing.T) {
		client := &sdkmock.SNSClient{}
		mockErr := assert.AnError

		client.On("ListTopicsWithContext", ctx, &sns.ListTopicsInput{}).
			Return((*sns.ListTopicsOutput)(nil), mockErr)

		output, err := client.ListTopicsWithContext(ctx, &sns.ListTopicsInput{})

		require.ErrorIs(t, err, mockErr)
		assert.Nil(t, output)
		client.AssertExpectations(t)
	})
}

func TestSNSClient_PublishWithContext(t *testing.T) {
	t.Run("should return expected PublishOutput without error", func(t *testing.T) {
		client := &sdkmock.SNSClient{}
		expectedOutput := &sns.PublishOutput{}

		client.On("PublishWithContext", ctx, &sns.PublishInput{}).
			Return(expectedOutput, nil)

		output, err := client.PublishWithContext(ctx, &sns.PublishInput{})

		require.NoError(t, err)
		assert.Equal(t, expectedOutput, output)
		client.AssertExpectations(t)
	})

	t.Run("should return error when mock is configured to fail", func(t *testing.T) {
		client := &sdkmock.SNSClient{}
		mockErr := assert.AnError

		client.On("PublishWithContext", ctx, &sns.PublishInput{}).
			Return((*sns.PublishOutput)(nil), mockErr)

		output, err := client.PublishWithContext(ctx, &sns.PublishInput{})

		require.ErrorIs(t, err, mockErr)
		assert.Nil(t, output)
		client.AssertExpectations(t)
	})
}
