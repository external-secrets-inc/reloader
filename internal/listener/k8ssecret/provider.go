package k8ssecret

import (
	"context"
	"errors"
	"sync"

	v1alpha1 "github.com/external-secrets-inc/reloader/api/v1alpha1"
	"github.com/external-secrets-inc/reloader/internal/events"
	"github.com/external-secrets-inc/reloader/internal/listener/schema"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Provider struct{}

// Creates a Kubernetes Secret Listener
func (p *Provider) CreateListener(ctx context.Context, config *v1alpha1.NotificationSource, client client.Client, eventChan chan events.SecretRotationEvent, logger logr.Logger) (schema.Listener, error) {
	if config == nil || config.KubernetesSecret == nil {
		return nil, errors.New("KubernetesSecret config is nil")
	}
	ctx, cancel := context.WithCancel(ctx)
	h := &Handler{
		config:     config.KubernetesSecret,
		context:    ctx,
		cancel:     cancel,
		client:     client,
		eventChan:  eventChan,
		logger:     logger,
		versionMap: sync.Map{},
	}

	return h, nil
}

func init() {
	schema.RegisterProvider(schema.KUBERNETES_SECRET, &Provider{})
}
