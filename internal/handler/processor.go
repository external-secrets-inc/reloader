package handler

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	esov1alpha1 "github.com/external-secrets-inc/reloader/api/v1alpha1"
	"github.com/external-secrets-inc/reloader/internal/events"
	"github.com/external-secrets-inc/reloader/internal/handler/schema"
)

type EventHandler struct {
	ctx                 context.Context
	client              client.Client
	secretsToWatchCache []esov1alpha1.DestinationToWatch
}

func NewEventHandler(client client.Client) *EventHandler {
	ctx := context.Background()
	return &EventHandler{
		ctx:    ctx,
		client: client,
	}
}

func (h *EventHandler) UpdateSecretsToWatch(secretsToWatch []esov1alpha1.DestinationToWatch) {
	h.secretsToWatchCache = secretsToWatch
}

func (h *EventHandler) HandleSecretRotationEvent(ctx context.Context, event events.SecretRotationEvent) error {
	logger := log.FromContext(ctx)
	for _, watchCriteria := range h.secretsToWatchCache {

		prov := schema.GetProvider(watchCriteria.Type)
		if prov == nil {
			logger.Info("Provider not found", "destination type", watchCriteria.Type)
			continue
		}
		h := prov.NewHandler(ctx, h.client, watchCriteria)
		// Mutate Handler for different Update and Match Strategies
		if watchCriteria.UpdateStrategy != nil {
			logger.Info("Optional Update strategies are not implemented", "UpdateStrategy", watchCriteria.UpdateStrategy)
		}
		if watchCriteria.MatchStrategy != nil {
			logger.Info("Optional Match strategies are not implemented", "MatchStrategy", watchCriteria.MatchStrategy)
		}
		objs, err := h.Filter(&watchCriteria, event)
		if err != nil {
			return fmt.Errorf("failed to filter objects:%w", err)
		}
		// Use Handler methods to figure out and apply objects
		for _, obj := range objs {
			isReferenced, err := h.References(obj, event.SecretIdentifier)
			if err != nil {
				return fmt.Errorf("failed to check if object is referenced:%w", err)
			}
			if !isReferenced {
				logger.V(1).Info("skipping object as its not referenced", "name", obj.GetName(), "namespace", obj.GetNamespace())
				continue
			}
			// object is referenced - apply
			err = h.Apply(obj, event)
			if err != nil {
				return fmt.Errorf("failed to apply object:%w", err)
			}
		}
	}
	return nil
}
