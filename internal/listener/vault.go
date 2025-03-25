package listener

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	v1alpha1 "github.com/external-secrets-inc/reloader/api/v1alpha1"
	"github.com/external-secrets-inc/reloader/internal/events"
	vault "github.com/external-secrets-inc/reloader/pkg/models/vault"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HashicorpVault represents a TCP socket listener configured to parse hashicorp vault messages.
// It utilizes a stop channel to manage its lifecycle.
type HashicorpVault struct {
	config    *v1alpha1.HashicorpVaultConfig
	context   context.Context
	cancel    context.CancelFunc
	client    client.Client
	eventChan chan events.SecretRotationEvent
	logger    logr.Logger
	tcpSocket *TCPSocket
}

// NewTCPSocketListener initializes a new TCP socket listener using the provided configuration and event channel.
func NewHashicorpVaultListener(ctx context.Context, config *v1alpha1.HashicorpVaultConfig, client client.Client, eventChan chan events.SecretRotationEvent, logger logr.Logger) (Listener, error) {
	ctx, cancel := context.WithCancel(ctx)
	h := &HashicorpVault{
		config:    config,
		context:   ctx,
		cancel:    cancel,
		client:    client,
		eventChan: eventChan,
		logger:    logger,
	}
	sockConfig := &v1alpha1.TCPSocketConfig{
		Host: config.Host,
		Port: config.Port,
	}
	sock := &TCPSocket{
		config:    sockConfig,
		context:   ctx,
		cancel:    cancel,
		client:    client,
		eventChan: eventChan,
		logger:    logger,
	}
	sock.SetProcessFn(h.processFn)
	h.tcpSocket = sock
	return h, nil
}

func (h *HashicorpVault) processFn(message []byte) {
	msg := &vault.AuditLog{}
	err := json.Unmarshal(message, msg)
	if err != nil {
		h.logger.Error(err, "Failed to unmarshal message", "Message", message)
		return
	}
	if !vault.ValidMessage(msg) {
		h.logger.V(1).Info("Invalid message - ignoring")
		return
	}
	basePath := msg.AuthResponse.MountPoint
	// Removing "data" if any
	path := strings.TrimPrefix(strings.Split(msg.AuthRequest.Path, basePath)[1], "data/")
	switch msg.AuthRequest.Operation {
	case "create":
	case "update":
		h.logger.V(1).Info("Received Valid Message", "Message", msg)
		event := events.SecretRotationEvent{
			SecretIdentifier:  path,
			RotationTimestamp: time.Now().Format("2006-01-02-15-04-05.000"),
			TriggerSource:     "vault",
		}
		h.eventChan <- event
		h.logger.V(1).Info("Published event to eventChan", "Event", event)
	default:
		h.logger.V(2).Info("Non-Applicable Operation", "Operation", msg.AuthRequest.Operation)
	}
}

// Start initiates the HashicorpVault service, making it ready to accept incoming connections.
func (h *HashicorpVault) Start() error {
	return h.tcpSocket.Start()
}

// Stop stops the HashicorpVault  by closing the stop channel.
func (h *HashicorpVault) Stop() error {
	h.cancel()
	return h.tcpSocket.Stop()
}
