package v1alpha1

// KubernetesConfigMapConfig contains configuration for Kubernetes notifications.
type KubernetesConfigMapConfig struct {
	// Server URL
	// +required
	ServerURL string `json:"serverURL"`

	// How to authenticate with Kubernetes cluster. If not specified, the default config is used.
	// +optional
	Auth *KubernetesAuth `json:"auth,omitempty"`
}
