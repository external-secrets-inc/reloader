// internal/notification/listener.go

package listener

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	esov1alpha1 "github.com/external-secrets-inc/reloader/api/v1alpha1"
	"github.com/external-secrets-inc/reloader/internal/events"
)

const (
	AWS_SQS          = "AwsSqs"
	AZURE_EVENT_GRID = "AzureEventGrid"
	GOOGLE_PUB_SUB   = "GooglePubSub"
	WEBHOOK          = "Webhook"
	TCP_SOCKET       = "TCPSocket"
	HASHICORP_VAULT  = "HashicorpVault"
	MOCK             = "Mock"
)

var (
	instance *NotificationListenerFactory
)

// Listener defines the interface for starting and stopping a listener.
type Listener interface {
	Start() error
	Stop() error
}

// Factory is an interface for creating event listeners for secret rotation events.
type Factory interface {
	CreateListener(ctx context.Context, source interface{}, eventChan chan events.SecretRotationEvent, logger logr.Logger) (Listener, error)
}

// NotificationListenerFactory provides a factory for creating notification listeners.
type NotificationListenerFactory struct {
	client client.Client
}

func InitializeFactory(client client.Client) *NotificationListenerFactory {
	if instance != nil {
		return instance
	}
	instance = &NotificationListenerFactory{
		client: client,
	}
	return instance
}

func GetFactory() *NotificationListenerFactory {
	return instance
}

// CreateListener creates a listener based on the provided notification source configuration and event channel.
func (f *NotificationListenerFactory) CreateListener(ctx context.Context, source interface{}, eventChan chan events.SecretRotationEvent, logger logr.Logger) (Listener, error) {
	switch src := source.(type) {
	// TODO[gusfcarvalho]: change this to be a provider/interface pattern
	// Similar to Destinations approach
	case esov1alpha1.NotificationSource:
		switch src.Type {
		case AWS_SQS:
			return NewAWSSQSListener(ctx, src.AwsSqs, f.client, eventChan, logger)
		case AZURE_EVENT_GRID:
			return NewAzureEventGridListener(ctx, src.AzureEventGrid, f.client, eventChan, logger)
		case GOOGLE_PUB_SUB:
			return NewGooglePubSubListener(ctx, src.GooglePubSub, f.client, eventChan, logger)
		case WEBHOOK:
			return NewWebhookListener(ctx, src.Webhook, f.client, eventChan, logger)
		case TCP_SOCKET:
			return NewTCPSocketListener(ctx, src.TCPSocket, f.client, eventChan, logger)
		case HASHICORP_VAULT:
			return NewHashicorpVaultListener(ctx, src.HashicorpVault, f.client, eventChan, logger)
		default:
			return nil, fmt.Errorf("unsupported notification source type: %s", src.Type)
		}
	default:
		return nil, fmt.Errorf("invalid notification source type: %T", source)
	}
}
