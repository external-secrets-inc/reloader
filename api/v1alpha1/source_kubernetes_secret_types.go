package v1alpha1

// KubernetesConfig contains configuration for Kubernetes notifications.
type KubernetesSecretConfig struct {
	// Server URL
	// +required
	ServerURL string `json:"serverURL"`

	// How to authenticate with Kubernetes cluster. If not specified, the default config is used.
	// +optional
	Auth *KubernetesAuth `json:"auth,omitempty"`
}

type KubernetesAuth struct {
	//+optional
	KubeConfigRef *KubeConfigRef `json:"kubeConfigRef,omitempty"`
	// Defines a CABundle if either TokenRef or ServiceAccountRef are used.
	// +optional
	CABundle string `json:"caBundle,omitempty"`
	//+optional
	TokenRef *TokenRef `json:"tokenRef,omitempty"`
	//+optional
	ServiceAccountRef *ServiceAccountSelector `json:"serviceAccountRef,omitempty"`
}

type KubeConfigRef struct {
	SecretRef SecretKeySelector `json:"secretRef"`
}

type TokenRef struct {
	SecretRef SecretKeySelector `json:"secretRef"`
}
