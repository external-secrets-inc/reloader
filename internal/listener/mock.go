package listener

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-logr/logr"

	v1alpha1 "github.com/external-secrets-inc/reloader/api/v1alpha1"
	"github.com/external-secrets-inc/reloader/internal/events"
)

// MockNotificationListener is a mock implementation of a notification listener for secret rotation events.
type MockNotificationListener struct {
	events       []events.SecretRotationEvent
	emitInterval time.Duration
	mu           sync.Mutex
	stopped      bool
	eventChan    chan events.SecretRotationEvent
}

// NewMockListener creates a new MockNotificationListener with specified events, emit interval, and event channel.
func NewMockListener(events []events.SecretRotationEvent, emitInterval time.Duration, eventChan chan events.SecretRotationEvent) *MockNotificationListener {
	return &MockNotificationListener{
		events:       events,
		emitInterval: emitInterval,
		eventChan:    eventChan,
	}
}

// Start initiates the emission of events from the MockNotificationListener. Returns an error if the listener has been stopped.
func (m *MockNotificationListener) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.stopped {
		return fmt.Errorf("listener has been stopped")
	}

	go func() {
		for _, event := range m.events {
			time.Sleep(m.emitInterval)
			m.eventChan <- event
		}
	}()

	return nil
}

// Stop signals the MockNotificationListener to stop emitting events.
func (m *MockNotificationListener) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopped = true
	return nil
}

// MockListenerFactory is a factory for creating mock listeners.
type MockListenerFactory struct{}

// NewMockListenerFactory creates and returns a new instance of MockListenerFactory.
func NewMockListenerFactory() *MockListenerFactory {
	return &MockListenerFactory{}
}

// CreateListener creates a mock listener for simulated secret rotation events.
func (f *MockListenerFactory) CreateListener(ctx context.Context, source interface{}, eventChan chan events.SecretRotationEvent, logger logr.Logger) (Listener, error) {
	// Convert source to NotificationSource
	n, ok := source.(v1alpha1.NotificationSource)
	if !ok {
		return nil, fmt.Errorf("invalid source type")
	}

	mockEvents := []events.SecretRotationEvent{
		{
			SecretIdentifier:  "aws://secret/arn:aws:secretsmanager:us-east-1:123456789012:secret:mysecret",
			RotationTimestamp: "2024-09-19T12:00:00Z",
			TriggerSource:     "aws-secretsmanager",
		},
	}
	return NewMockListener(mockEvents, time.Duration(n.Mock.EmitInterval)*time.Millisecond, eventChan), nil
}
