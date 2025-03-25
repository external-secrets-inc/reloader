package listener_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/external-secrets-inc/providers-listeners/mocks"
	"github.com/external-secrets-inc/providers-listeners/pkg"
	authAWS "github.com/external-secrets-inc/providers-listeners/pkg/auth/aws"
	"github.com/external-secrets-inc/providers-listeners/pkg/listener"
	"github.com/go-logr/logr/testr"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Mock setup
func createMockListener(t *testing.T, mockSQS *mocks.MockSQSClientInterface) (*listener.AWSSQSListener, error) {
	ctx := context.Background()
	k8sClient := mocks.CreateFakeK8sClient().Build()
	authConfig := pkg.AWSSDKAuth{
		AuthMethod: authAWS.AuthMethodIRSA,
		Region:     "us-east-1",
		ServiceAccount: &pkg.ServiceAccountSelector{
			Name:      "aws-sa",
			Namespace: "default",
		},
	}

	config := &pkg.AWSSQSConfig{
		Auth: authConfig,
	}
	logger := testr.New(t)

	// Create the listener using the constructor function
	listenerInstance, err := listener.NewAWSSQSListener(ctx, config, k8sClient, logger)
	if err != nil {
		return nil, err
	}

	// Override the SQS client with the mock
	if err := listenerInstance.SetSQSClient(mockSQS); err != nil {
		return nil, err
	}
	return listenerInstance, nil
}

// Test polling with a successful message retrieval
func TestPollMessages_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSQS := mocks.NewMockSQSClientInterface(ctrl)
	listener, err := createMockListener(t, mockSQS)
	if err != nil {
		t.Fatal(err)
	}

	// Mock SQS response
	mockSQS.EXPECT().ReceiveMessage(gomock.Any(), gomock.Any()).Return(&sqs.ReceiveMessageOutput{
		Messages: []types.Message{
			{MessageId: aws.String("msg-1"), Body: aws.String("Hello")},
			{MessageId: aws.String("msg-2"), Body: aws.String("World")},
		},
	}, nil)

	mockSQS.EXPECT().DeleteMessage(gomock.Any(), gomock.Any()).Times(2).Return(nil, nil) // Expect delete calls

	messages, err := listener.PollMessages()
	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, "Hello", *messages[0].Body)
	assert.Equal(t, "World", *messages[1].Body)
}

// Test polling when SQS returns an error
func TestPollMessages_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSQS := mocks.NewMockSQSClientInterface(ctrl)
	listener, err := createMockListener(t, mockSQS)
	if err != nil {
		t.Fatal(err)
	}

	// Mock failure
	mockSQS.EXPECT().ReceiveMessage(gomock.Any(), gomock.Any()).Return(nil, errors.New("SQS error"))

	messages, err := listener.PollMessages()
	assert.Nil(t, messages)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SQS error")
}

// Test Start function yielding messages
func TestStart_YieldsMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSQS := mocks.NewMockSQSClientInterface(ctrl)
	listener, err := createMockListener(t, mockSQS)
	if err != nil {
		t.Fatal(err)
	}

	// Mock SQS response
	mockSQS.EXPECT().ReceiveMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(&sqs.ReceiveMessageOutput{
		Messages: []types.Message{
			{MessageId: aws.String("msg-1"), Body: aws.String("Hello")},
		},
	}, nil)

	mockSQS.EXPECT().DeleteMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)

	msgCh, _ := listener.Start()

	// Wait for a message
	select {
	case messages := <-msgCh:
		assert.Len(t, messages, 1)
		assert.Equal(t, "Hello", *messages[0].Body)
	case <-time.After(2 * time.Second):
		t.Fatal("Expected message but timed out")
	}

	// Stop listener
	if err := listener.Stop(); err != nil {
		t.Fatal(err)
	}
}

// Test Start handling polling errors
func TestStart_HandlesErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSQS := mocks.NewMockSQSClientInterface(ctrl)
	listener, err := createMockListener(t, mockSQS)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate an error from SQS
	mockSQS.EXPECT().ReceiveMessage(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, errors.New("SQS failure"))

	_, errCh := listener.Start()

	// Expect an error
	select {
	case err := <-errCh:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SQS failure")
	case <-time.After(2 * time.Second):
		t.Fatal("Expected error but timed out")
	}

	// Stop listener
	if err := listener.Stop(); err != nil {
		t.Fatal(err)
	}
}

// Test Stop function ensures graceful shutdown
func TestStop_EnsuresShutdown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSQS := mocks.NewMockSQSClientInterface(ctrl)
	listener, err := createMockListener(t, mockSQS)
	if err != nil {
		t.Fatal(err)
	}

	msgCh, errCh := listener.Start()

	// Stop the listener
	if err := listener.Stop(); err != nil {
		t.Fatal(err)
	}

	// Ensure channels are closed
	_, msgOpen := <-msgCh
	_, errOpen := <-errCh

	assert.False(t, msgOpen, "Message channel should be closed")
	assert.False(t, errOpen, "Error channel should be closed")
}
