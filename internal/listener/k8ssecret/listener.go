package k8ssecret

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"sync"
	"time"

	v1alpha1 "github.com/external-secrets-inc/reloader/api/v1alpha1"
	"github.com/external-secrets-inc/reloader/internal/events"
	"github.com/external-secrets-inc/reloader/internal/listener/schema"
	"github.com/external-secrets-inc/reloader/pkg/util"
	"github.com/go-logr/logr"
	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Handler creates a kubernetes Secret watch on a given cluster and sends event
// on any operation propagated to the watch (create, update, delete).
type Handler struct {
	config        *v1alpha1.KubernetesSecretConfig
	context       context.Context
	cancel        context.CancelFunc
	client        client.Client
	mgr           ctrl.Manager
	ctrlClientSet typedcorev1.CoreV1Interface
	eventChan     chan events.SecretRotationEvent
	logger        logr.Logger
	versionMap    sync.Map // map[types.NamespacedName]string
}

// Start initiates the Kubernetes Secret listener.
func (h *Handler) Start() error {
	log := ctrl.Log.WithName("secret-watcher")
	cfg, err := h.getKubeConfig(h.config)
	if err != nil {
		return fmt.Errorf("could not parse config: %w", err)
	}
	h.ctrlClientSet, err = typedcorev1.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("could not create client set: %w", err)
	}
	manager, err := ctrl.NewManager(cfg, ctrl.Options{})
	if err != nil {
		return fmt.Errorf("could not create manager: %w", err)
	}
	err = ctrl.
		NewControllerManagedBy(manager).
		Named("k8ssecret").
		Watches(&corev1.Secret{}, handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, s client.Object) []ctrl.Request {
			secret, ok := s.(*corev1.Secret)
			if !ok {
				log.Error(err, "while processing secret")
				return nil
			}
			if secret.DeletionTimestamp != nil {
				log.V(2).Info("skipping deleted secret", "namespace", secret.GetNamespace(), "name", secret.GetName())
				return nil
			}
			version := secret.GetResourceVersion()
			storedVersion, loaded := h.versionMap.LoadOrStore(types.NamespacedName{Namespace: secret.GetNamespace(), Name: secret.GetName()}, version)
			if !loaded {
				log.V(2).Info("secret not added to cache, skipping", "namespace", secret.GetNamespace(), "name", secret.GetName())
				return nil
			}
			if version == storedVersion {
				// Happens for some weird reason
				log.V(2).Info("skipping secret with same version", "namespace", secret.GetNamespace(), "name", secret.GetName())
				return nil
			}
			// Safe to send event
			h.versionMap.Store(types.NamespacedName{Namespace: secret.GetNamespace(), Name: secret.GetName()}, version)
			h.eventChan <- events.SecretRotationEvent{
				SecretIdentifier:  secret.GetName(),
				Namespace:         secret.GetNamespace(),
				RotationTimestamp: time.Now().Format(time.RFC3339),
				TriggerSource:     fmt.Sprintf("%s/%s", schema.KUBERNETES_SECRET, secret.GetName()),
			}
			return nil
		})).
		Complete(reconcile.Func(func(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
			// We dont need to reconcile anything, as we are sending this over another controller
			return reconcile.Result{}, nil
		}))
	if err != nil {
		return fmt.Errorf("could not create controller: %w", err)
	}
	h.mgr = manager
	go func() {
		if err := manager.Start(h.context); err != nil {
			h.logger.Error(err, "failed to start secrets watching")
		}
	}()
	return nil
}

// Stop stops the Watch by closing the stop channel.
func (h *Handler) Stop() error {
	h.cancel()
	return nil
}

func (h *Handler) getKubeConfig(config *v1alpha1.KubernetesSecretConfig) (*rest.Config, error) {
	if config.Auth == nil {
		h.logger.V(1).Info("no auth specified - using default config")
		cfg, err := ctrl.GetConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get default config: %w", err)
		}
		return cfg, nil
	}
	if config.Auth.KubeConfigRef != nil {
		cfg, err := h.fetchSecretKey(h.context, &config.Auth.KubeConfigRef.SecretRef)
		if err != nil {
			return nil, err
		}

		return clientcmd.RESTConfigFromKubeConfig(cfg)
	}

	if h.config.ServerURL == "" {
		return nil, errors.New("no server URL provided")
	}

	cfg := &rest.Config{
		Host: h.config.ServerURL,
	}

	ca, err := h.fetchCACertFromSource(h.context, certOpts{
		CABundle: []byte(config.Auth.CABundle),
	})
	if err != nil {
		return nil, err
	}

	cfg.TLSClientConfig = rest.TLSClientConfig{
		Insecure: false,
		CAData:   ca,
	}

	switch {
	case config.Auth.TokenRef != nil:
		token, err := h.fetchSecretKey(h.context, &config.Auth.TokenRef.SecretRef)
		if err != nil {
			return nil, fmt.Errorf("could not fetch Auth.Token.BearerToken: %w", err)
		}
		cfg.BearerToken = string(token)
	case config.Auth.ServiceAccountRef != nil:
		token, err := h.serviceAccountToken(h.context, config.Auth.ServiceAccountRef)
		if err != nil {
			return nil, fmt.Errorf("could not fetch Auth.ServiceAccount: %w", err)
		}
		cfg.BearerToken = string(token)
	default:
		return nil, errors.New("no auth provider given")
	}

	return cfg, nil

}

func (h *Handler) fetchSecretKey(ctx context.Context, secretKeySelector *v1alpha1.SecretKeySelector) ([]byte, error) {
	if secretKeySelector == nil {
		return nil, errors.New("secret key selector is nil")
	}
	secret, err := util.GetSecret(ctx, h.client, secretKeySelector.Name, secretKeySelector.Namespace, h.logger)
	if err != nil {
		return nil, fmt.Errorf("could not fetch secret key: %w", err)
	}
	data, ok := secret.Data[secretKeySelector.Key]
	if !ok {
		return nil, fmt.Errorf("key %s not found in secret %s", secretKeySelector.Key, secretKeySelector.Name)
	}
	return data, nil
}

func (h *Handler) serviceAccountToken(ctx context.Context, serviceAccountSelector *v1alpha1.ServiceAccountSelector) ([]byte, error) {
	if serviceAccountSelector == nil {
		return nil, errors.New("service account selector is nil")
	}
	namespace := serviceAccountSelector.Namespace
	expirationSeconds := int64(3600)
	tr, err := h.ctrlClientSet.ServiceAccounts(namespace).CreateToken(ctx, serviceAccountSelector.Name, &authenticationv1.TokenRequest{
		Spec: authenticationv1.TokenRequestSpec{
			Audiences:         serviceAccountSelector.Audiences,
			ExpirationSeconds: &expirationSeconds,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not create token: %w", err)
	}
	return []byte(tr.Status.Token), nil
}

func (h *Handler) fetchCACertFromSource(ctx context.Context, certOpts certOpts) ([]byte, error) {
	if len(certOpts.CABundle) > 0 {
		pem, err := base64decode(certOpts.CABundle)
		if err != nil {
			return nil, fmt.Errorf("failed to decode ca bundle: %w", err)
		}

		return pem, nil
	}
	return nil, nil
}

type certOpts struct {
	CABundle []byte
}

func base64decode(cert []byte) ([]byte, error) {
	if c, err := parseCertificateBytes(cert); err == nil {
		return c, nil
	}

	// try b64 decoding and test for validity again...
	certificate, err := base64.StdEncoding.DecodeString(string(cert))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return parseCertificateBytes(certificate)
}

func parseCertificateBytes(certBytes []byte) ([]byte, error) {
	block, _ := pem.Decode(certBytes)
	if block == nil {
		return nil, errors.New("failed to parse the new certificate, not valid pem data")
	}

	if _, err := x509.ParseCertificate(block.Bytes); err != nil {
		return nil, fmt.Errorf("failed to validate certificate: %w", err)
	}

	return certBytes, nil
}
