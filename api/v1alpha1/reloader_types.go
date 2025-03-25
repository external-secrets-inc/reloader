/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigSpec defines the desired state of a Reloader Config
type ConfigSpec struct {
	// NotificationSources specifies the notification systems to listen to.
	// +required
	NotificationSources []NotificationSource `json:"notificationSources"`

	// DestinationsToWatch specifies which secrets the controller should monitor.
	// +required
	DestinationsToWatch []DestinationToWatch `json:"destinationsToWatch"`
}

// NotificationSource represents a notification system configuration.
type NotificationSource struct {
	// Type of the notification source (e.g., AwsSqs, AzureEventGrid, GooglePubSub, HashicorpVault, Webhook, TCPSocket).
	// +kubebuilder:validation:Enum=AwsSqs;AzureEventGrid;GooglePubSub;HashicorpVault;Webhook;TCPSocket
	// +required
	Type string `json:"type"`

	// AwsSqs configuration (required if Type is AwsSqs).
	// +optional
	AwsSqs *AWSSQSConfig `json:"awsSqs,omitempty"`

	AzureEventGrid *AzureEventGridConfig `json:"azureEventGrid,omitempty"`

	// GooglePubSub configuration (required if Type is GooglePubSub).
	// +optional
	GooglePubSub *GooglePubSubConfig `json:"googlePubSub,omitempty"`

	// Webhook configuration (required if Type is Webhook).
	// +optional
	Webhook *WebhookConfig `json:"webhook,omitempty"`

	// HashicorpVault configuration (required if Type is HashicorpVault).
	// +optional
	HashicorpVault *HashicorpVaultConfig `json:"hashicorpVault,omitempty"`

	// TCPSocket configuration (required if Type is TCPSocket).
	// +optional
	TCPSocket *TCPSocketConfig `json:"tcpSocket,omitempty"`

	// Mock configuration (optional field for testing purposes).
	Mock *MockConfig `json:"mock,omitempty"`
}

// AWSSQSConfig contains configuration for AWS SDK.
type AWSSQSConfig struct {
	// QueueURL is the URL of the AWS SDK queue.
	// +required
	QueueURL string `json:"queueURL"`

	// Authentication methods for AWS.
	// +required
	Auth AWSSDKAuth `json:"auth"`

	// MaxNumberOfMessages specifies the maximum number of messages to retrieve from the SDK queue in a single request.
	// +optional
	// +kubebuilder:default=10
	MaxNumberOfMessages int32 `json:"numberOfMessages"`

	// WaitTimeSeconds specifies the duration (in seconds) to wait for messages in the SDK queue before returning.
	// +optional
	// +kubebuilder:default=20
	WaitTimeSeconds int32 `json:"waitTimeSeconds"`

	// VisibilityTimeout specifies the duration (in seconds) that a message received from the SDK queue is hidden from subsequent retrievals.
	// +optional
	// +kubebuilder:default=30
	VisibilityTimeout int32 `json:"visibilityTimeout"`
}

// AWSSDKAuth contains authentication methods for AWS SDK.
type AWSSDKAuth struct {
	AuthMethod string `json:"authMethod"`

	Region string `json:"region"`

	ServiceAccount *ServiceAccountSelector `json:"serviceAccountRef,omitempty"`

	SecretRef *AWSSDKSecretRef `json:"secretRef,omitempty"`
}

type AWSSDKSecretRef struct {
	AccessKeyId     SecretKeySelector `json:"accessKeyIdSecretRef"`
	SecretAccessKey SecretKeySelector `json:"secretAccessKeySecretRef"`
}

// GooglePubSubConfig contains configuration for Google Pub/Sub.
type GooglePubSubConfig struct {
	// SubscriptionID is the ID of the Pub/Sub subscription.
	// +required
	SubscriptionID string `json:"subscriptionID"`

	// ProjectID is the GCP project ID where the subscription exists.
	// +required
	ProjectID string `json:"projectID"`

	// Authentication methods for Google Pub/Sub.
	// +optional
	Auth *GooglePubSubAuth `json:"auth,omitempty"`
}

// GooglePubSubAuth contains authentication methods for Google Pub/Sub.
type GooglePubSubAuth struct {
	// +optional
	SecretRef *GCPSMAuthSecretRef `json:"secretRef,omitempty"`
	// +optional
	WorkloadIdentity *GCPWorkloadIdentity `json:"workloadIdentity,omitempty"`
}

type GCPSMAuthSecretRef struct {
	// The SecretAccessKey is used for authentication
	// +optional
	SecretAccessKey SecretKeySelector `json:"secretAccessKeySecretRef,omitempty"`
}

type GCPWorkloadIdentity struct {
	ServiceAccountRef ServiceAccountSelector `json:"serviceAccountRef"`
	ClusterLocation   string                 `json:"clusterLocation"`
	ClusterName       string                 `json:"clusterName"`
	ClusterProjectID  string                 `json:"clusterProjectID,omitempty"`
}

type AzureEventGridConfig struct {
	Host string `json:"host"`

	// +required
	// +kubebuilder:default=8080
	Port int32 `json:"port"`

	Subscriptions []string `json:"subscriptions"`
}

// WebhookConfig contains configuration for Webhook notifications.
type WebhookConfig struct {
	// Path that the webhook will receive the notifications.
	// If not present `/webhook` will be used. The path always expects a POST and this is not configurable
	// +optional
	Path string `json:"path"`

	// Address is the address where the webhook will be served in your infrastructure.
	// If not present, defaults to `:8090`
	// +optional
	Address string `json:"address"`

	// SecretIdentifierOnPayload is the key that the rotator will look for in the payload.
	// The value of this key should be the same name as in the external secret. It will default to `0.data.ObjectName` if not set
	// +optional
	SecretIdentifierOnPayload string `json:"identifierPathOnPayload,omitempty"`

	// Auth is the authentication method for the webhook
	// +optional
	Auth *WebhookAuth `json:"webhookAuth,omitempty"`

	// RetryPolicy represents the policy to retry when a message fails.
	// If it's empty, reloader will return a 4xx and won't retry.
	// +optional
	RetryPolicy *RetryPolicy `json:"retryPolicy,omitempty"`
}

type RetryPolicy struct {
	// MaxRetries represents the maximum times the reloader should retry to process a message. Numbers greater than 10 will be ignored and 10 will be used instead
	// +optional
	MaxRetries int `json:"maxRetries"`

	// Algorithm represents how watiting time will change for each retry.
	// Currently supports "linear" and "exponential". If an invalid string or null is given, "exponential" will be used
	// +optional
	Algorithm string `json:"algorithm"`
}

// WebhookAuth contains authentication methods for webhooks.
type WebhookAuth struct {
	// BasicAuth contains basic authentication credentials.
	// +optional
	BasicAuth *BasicAuth `json:"basicAuth,omitempty"`

	// BearerToken references a Kubernetes Secret containing the bearer token.
	// +optional
	BearerToken *BearerToken `json:"bearerToken,omitempty"`
}

// TCPSocketConfig contains configuration for TCP Socket notifications.
type TCPSocketConfig struct {
	// Host is the hostname or IP address to listen on.
	// +required
	Host string `json:"host"`

	// Port is the port number to listen on.
	// +required
	// +kubebuilder:default=8000
	Port int32 `json:"port"`

	// SecretIdentifierOnPayload is the key that the rotator will look for in the payload.
	// The value of this key should be the same name as in the external secret. It will default to `0.data.ObjectName` if not set
	SecretIdentifierOnPayload string `json:"identifierPathOnPayload,omitempty"`
}

// HashicorpVault contains configuration for HashicorpVault notifications.
type HashicorpVaultConfig struct {
	// Host is the hostname or IP address to listen on.
	// +required
	Host string `json:"host"`

	// Port is the port number to listen on.
	// +required
	// +kubebuilder:default=8000
	Port int32 `json:"port"`
}

// MockConfig represents configuration settings for mock notifications.
type MockConfig struct {
	EmitInterval int32 `json:"emitInterval"`
}

// BasicAuth contains basic authentication credentials.
type BasicAuth struct {
	// UsernameSecretRef contains a secret reference for the username
	// +required
	UsernameSecretRef SecretKeySelector `json:"usernameSecretRef,omitempty"`

	// PasswordSecretRef contains a secret reference for the password
	// +required
	PasswordSecretRef SecretKeySelector `json:"passwordSecretRef,omitempty"`
}

// BearerToken contains the bearer token credentials.
type BearerToken struct {
	// BearerTokenSecretRef references a Kubernetes Secret containing the bearer token.
	// +required
	BearerTokenSecretRef SecretKeySelector `json:"bearerTokenSecretRef"`
}

// TLSClientAuth contains client certificate authentication details.
type TLSClientAuth struct {
	// ClientCertSecretRef references a Kubernetes Secret containing the client certificate.
	// +required
	ClientCertSecretRef SecretKeySelector `json:"clientCertSecretRef"`

	// ClientKeySecretRef references a Kubernetes Secret containing the client key.
	// +required
	ClientKeySecretRef SecretKeySelector `json:"clientKeySecretRef"`
}

type ServiceAccountSelector struct {

	// Name specifies the name of the service account to be selected.
	// +required
	Name string `json:"name"`

	// ServiceAccountSelector represents a Kubernetes service account with a name and namespace for selection purposes.
	// +required
	Namespace string `json:"namespace"`
	// Audience specifies the `aud` claim for the service account token
	// If the service account uses a well-known annotation for e.g. IRSA or GCP Workload Identity
	// then this audiences will be appended to the list
	// +optional
	Audiences []string `json:"audiences,omitempty"`
}

// SecretKeySelector is used to reference a specific secret within a Kubernetes namespace.
// It contains the name of the secret and the namespace where it resides.
type SecretKeySelector struct {

	// Name specifies the name of the referenced Kubernetes secret.
	// +required
	Name string `json:"name"`

	// Key specifies the key within the referenced Kubernetes secret.
	// +required
	Key string `json:"key"`

	// Namespace specifies the Kubernetes namespace where the referenced secret resides.
	// +required
	Namespace string `json:"namespace"`
}

// DestinationToWatch specifies the criteria for monitoring secrets in the cluster.
type DestinationToWatch struct {
	// Type specifies the type of destination to watch.
	// +required
	// +kubebuilder:validation:Enum=generic;externalsecret;deployment;secret;certificate
	Type string `json:"type"`
	// GenericDestination specifies the destination to watch
	// +optional
	Generic *GenericDestination `json:"generic,omitempty"`
	// +optional
	ExternalSecret *ExternalSecretDestination `json:"externalsecret,omitempty"`
	// +optional
	Deployment *DeploymentDestination `json:"deployment,omitempty"`
	// +optional
	Secret *SecretDestination `json:"secret,omitempty"`
	// +optional
	Certificate *CertificateDestination `json:"certificate,omitempty"`
	//UpdateStrategy. If not specified, will use each destinations' default update strategy.
	// +optional
	UpdateStrategy *UpdateStrategy `json:"updateStrategy,omitempty"`
	//MatchStrategy. If not specified, will use each destinations' default match strategy.
	// +optional
	MatchStrategy *MatchStrategy `json:"matchStrategy,omitempty"`
	//WaitStrategy. If not specified, will use each destinations's default wait strategy.
	// +optional
	WaitStrategy *WaitStrategy `json:"waitStrategy,omitempty"`
}

// Defines a DeploymentDestination. Behavior is a pod templates annotations patch.
// Default UpdateStrategy is pod template annotations patch to trigger a new rollout.
// Default MatchStrategy is matching secret-key with any of:
// * Equality against `spec.template.spec.containers[*].env[*].valueFrom.secretKeyRef.name`
// * Equality against `spec.template.spec.containers[*].envFrom.secretRef.name`
// Default WaitStrategy is to wait for the rollout to be completed with 3 minutes of grace period before
// moving to the next matched deployment.
type DeploymentDestination struct {
	// NamespaceSelectors selects namespaces based on labels.
	// The manifest must reside in a namespace that matches at least one of these selectors.
	// +optional
	NamespaceSelectors []metav1.LabelSelector `json:"namespaceSelectors,omitempty"`

	// LabelSelectors selects resources based on their labels.
	// The resource must satisfy all conditions defined in this selector.
	// Supports both matchLabels and matchExpressions for advanced filtering.
	// +optional
	LabelSelectors *metav1.LabelSelector `json:"labelSelectors,omitempty"`

	// Names specifies a list of resource names to watch.
	// The resource must have a name that matches one of these entries.
	// +optional
	Names []string `json:"names,omitempty"`
}

// Defines an ExternalSecretDestination. Behavior is an annotations patch.
// Default UpdateStrategy is annotations patch to trigger externalSecret reconcile.
// Default MatchStrategy is matching secret-key with any of:
// * Equality against `spec.data.remoteRef.key`
// * Equality against `spec.dataFrom.remoteRef.key`
// * Regexp against `spec.dataFrom.find.name.regexp`
type ExternalSecretDestination struct {
	// NamespaceSelectors selects namespaces based on labels.
	// The manifest must reside in a namespace that matches at least one of these selectors.
	// +optional
	NamespaceSelectors []metav1.LabelSelector `json:"namespaceSelectors,omitempty"`

	// LabelSelectors selects resources based on their labels.
	// The resource must satisfy all conditions defined in this selector.
	// Supports both matchLabels and matchExpressions for advanced filtering.
	// +optional
	LabelSelectors *metav1.LabelSelector `json:"labelSelectors,omitempty"`

	// Names specifies a list of resource names to watch.
	// The resource must have a name that matches one of these entries.
	// +optional
	Names []string `json:"names,omitempty"`
}

// Defines an SecretDestination.
// Default UpdateStrategy is a annotations patch operation to trigger reconcile of any Owners of that Secret.
// Default MatchStrategy is matching secret-key with any of:
// * Equality against `spec.data`
type SecretDestination struct {
	// NamespaceSelectors selects namespaces based on labels.
	// The manifest must reside in a namespace that matches at least one of these selectors.
	// +optional
	NamespaceSelectors []metav1.LabelSelector `json:"namespaceSelectors,omitempty"`

	// LabelSelectors selects resources based on their labels.
	// The resource must satisfy all conditions defined in this selector.
	// Supports both matchLabels and matchExpressions for advanced filtering.
	// +optional
	LabelSelectors *metav1.LabelSelector `json:"labelSelectors,omitempty"`

	// Names specifies a list of resource names to watch.
	// The resource must have a name that matches one of these entries.
	// +optional
	Names []string `json:"names,omitempty"`
}

// Defines a Certificate (cert-manager) Destination.
// Default UpdateStrategy is status/patch to trigger cert renewal.
// Default MatchStrategy is matching secret-key with `spec.secretName`
type CertificateDestination struct {
	// NamespaceSelectors selects namespaces based on labels.
	// The manifest must reside in a namespace that matches at least one of these selectors.
	// +optional
	NamespaceSelectors []metav1.LabelSelector `json:"namespaceSelectors,omitempty"`

	// LabelSelectors selects resources based on their labels.
	// The resource must satisfy all conditions defined in this selector.
	// Supports both matchLabels and matchExpressions for advanced filtering.
	// +optional
	LabelSelectors *metav1.LabelSelector `json:"labelSelectors,omitempty"`

	// Names specifies a list of resource names to watch.
	// The resource must have a name that matches one of these entries.
	// +optional
	Names []string `json:"names,omitempty"`
}
type GenericDestination struct {
	// ResourceRef specifies the resource destination to watch
	ResourceRef ResourceRef `json:"resourceRef"`
}

type WaitStrategy struct {
	// Waits for a given time interval to reconcile the next object
	//+optional
	Time *metav1.Duration `json:"time,omitempty"`
	// Waits for a given status condition to be met
	//+optional
	Condition *WaitForCondition `json:"condition,omitempty"`
}

type WaitForCondition struct {
	// Period to wait before each retry
	//+optional
	RetryTimeout *metav1.Duration `json:"retryTimeout,omitempty"`
	// Maximum retries to check for a condition
	//+optional
	MaxRetries *int32 `json:"maxRetries,omitempty"`
	// The name of the condition to wait for
	//+required
	Type string `json:"type"`
	// The status of the condition to wait for
	//+optional
	Status string `json:"status"`
	// Optional message to match
	//+optional
	Message string `json:"message,omitempty"`
	// Optional reason to match
	//+optional
	Reason string `json:"reason,omitempty"`
	// Only accept this condition after a given period from the transition time
	//+optional
	TransitionedAfter *metav1.Duration `json:"transitionedAfter,omitempty"`
	// Only accept this condition after a given period from the update time
	//+optional
	UpdatedAfter *metav1.Duration `json:"updatedAfter,omitempty"`
}

type MatchStrategy struct {
	Path       string      `json:"path"`
	Conditions []Condition `json:"conditions"`
}

type Condition struct {
	Value     string             `json:"value"`
	Operation ConditionOperation `json:"operation"`
}

type ConditionOperation string

const (
	ConditionOperationEqual       ConditionOperation = "Equal"
	ConditionOperationNotEqual    ConditionOperation = "NotEqual"
	ConditionOperationContains    ConditionOperation = "Contains"
	ConditionOperationNotContains ConditionOperation = "NotContains"
	ConditionOperationIn          ConditionOperation = "RegularExpression"
)

type UpdateStrategy struct {
	Operation UpdateStrategyOperation `json:"operation"`
	// Required if Operation == Patch
	PatchOperationConfig *PatchOperationConfig `json:"patchOperationConfig,omitempty"`
}

type PatchOperationConfig struct {
	Path     string `json:"path"`
	Template string `json:"template"`
}

type UpdateStrategyOperation string

const (
	UpdateStrategyOperationPatchStatus UpdateStrategyOperation = "PatchStatus"
	UpdateStrategyOperationPatch       UpdateStrategyOperation = "Patch"
	UpdateStrategyOperationDelete      UpdateStrategyOperation = "Delete"
)

type ResourceRef struct {
	// API Version of the resource destination to watch
	APIVersion string `json:"apiVersion"`
	// Kind of the resource destination to watch
	Kind string `json:"kind"`
	// NamespaceSelectors selects namespaces based on labels.
	// The manifest must reside in a namespace that matches at least one of these selectors.
	// +optional
	NamespaceSelectors []metav1.LabelSelector `json:"namespaceSelectors,omitempty"`

	// LabelSelectors selects resources based on their labels.
	// The resource must satisfy all conditions defined in this selector.
	// Supports both matchLabels and matchExpressions for advanced filtering.
	// +optional
	LabelSelectors *metav1.LabelSelector `json:"labelSelectors,omitempty"`

	// Names specifies a list of resource names to watch.
	// The resource must have a name that matches one of these entries.
	// +optional
	Names []string `json:"names,omitempty"`
}

// ConfigStatus defines the observed state of Reloader
type ConfigStatus struct {
	// Conditions represent the latest available observations of the resource's state.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:subresource:status

// Config is the Schema for the reloader config API
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigSpec   `json:"spec,omitempty"`
	Status ConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ConfigList contains a list of Config
type ConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Config `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Config{}, &ConfigList{})
}
