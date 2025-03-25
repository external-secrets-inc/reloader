package v1alpha1

// KubernetesConfig contains configuration for Kubernetes notifications.
type KubernetesConfig struct {
	// Server URL
	// +required
	ServerURL string `json:"serverURL"`

	// +required
	CABundle string `json:"caBundle"`

	// +required
	Auth KubernetesAuth `json:"auth"`
}

type KubernetesAuth struct {
	//+optional
	KubeConfigRef *KubeConfigRef `json:"kubeConfigRef,omitempty"`
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
