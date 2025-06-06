//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSSDKAuth) DeepCopyInto(out *AWSSDKAuth) {
	*out = *in
	if in.ServiceAccount != nil {
		in, out := &in.ServiceAccount, &out.ServiceAccount
		*out = new(ServiceAccountSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(AWSSDKSecretRef)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSSDKAuth.
func (in *AWSSDKAuth) DeepCopy() *AWSSDKAuth {
	if in == nil {
		return nil
	}
	out := new(AWSSDKAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSSDKSecretRef) DeepCopyInto(out *AWSSDKSecretRef) {
	*out = *in
	out.AccessKeyId = in.AccessKeyId
	out.SecretAccessKey = in.SecretAccessKey
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSSDKSecretRef.
func (in *AWSSDKSecretRef) DeepCopy() *AWSSDKSecretRef {
	if in == nil {
		return nil
	}
	out := new(AWSSDKSecretRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSSQSConfig) DeepCopyInto(out *AWSSQSConfig) {
	*out = *in
	in.Auth.DeepCopyInto(&out.Auth)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSSQSConfig.
func (in *AWSSQSConfig) DeepCopy() *AWSSQSConfig {
	if in == nil {
		return nil
	}
	out := new(AWSSQSConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureEventGridConfig) DeepCopyInto(out *AzureEventGridConfig) {
	*out = *in
	if in.Subscriptions != nil {
		in, out := &in.Subscriptions, &out.Subscriptions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureEventGridConfig.
func (in *AzureEventGridConfig) DeepCopy() *AzureEventGridConfig {
	if in == nil {
		return nil
	}
	out := new(AzureEventGridConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BasicAuth) DeepCopyInto(out *BasicAuth) {
	*out = *in
	out.UsernameSecretRef = in.UsernameSecretRef
	out.PasswordSecretRef = in.PasswordSecretRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BasicAuth.
func (in *BasicAuth) DeepCopy() *BasicAuth {
	if in == nil {
		return nil
	}
	out := new(BasicAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BearerToken) DeepCopyInto(out *BearerToken) {
	*out = *in
	out.BearerTokenSecretRef = in.BearerTokenSecretRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BearerToken.
func (in *BearerToken) DeepCopy() *BearerToken {
	if in == nil {
		return nil
	}
	out := new(BearerToken)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Condition) DeepCopyInto(out *Condition) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Condition.
func (in *Condition) DeepCopy() *Condition {
	if in == nil {
		return nil
	}
	out := new(Condition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Config) DeepCopyInto(out *Config) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Config.
func (in *Config) DeepCopy() *Config {
	if in == nil {
		return nil
	}
	out := new(Config)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Config) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigList) DeepCopyInto(out *ConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Config, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigList.
func (in *ConfigList) DeepCopy() *ConfigList {
	if in == nil {
		return nil
	}
	out := new(ConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigSpec) DeepCopyInto(out *ConfigSpec) {
	*out = *in
	if in.NotificationSources != nil {
		in, out := &in.NotificationSources, &out.NotificationSources
		*out = make([]NotificationSource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.DestinationsToWatch != nil {
		in, out := &in.DestinationsToWatch, &out.DestinationsToWatch
		*out = make([]DestinationToWatch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigSpec.
func (in *ConfigSpec) DeepCopy() *ConfigSpec {
	if in == nil {
		return nil
	}
	out := new(ConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigStatus) DeepCopyInto(out *ConfigStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigStatus.
func (in *ConfigStatus) DeepCopy() *ConfigStatus {
	if in == nil {
		return nil
	}
	out := new(ConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeploymentDestination) DeepCopyInto(out *DeploymentDestination) {
	*out = *in
	if in.NamespaceSelectors != nil {
		in, out := &in.NamespaceSelectors, &out.NamespaceSelectors
		*out = make([]v1.LabelSelector, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.LabelSelectors != nil {
		in, out := &in.LabelSelectors, &out.LabelSelectors
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Names != nil {
		in, out := &in.Names, &out.Names
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeploymentDestination.
func (in *DeploymentDestination) DeepCopy() *DeploymentDestination {
	if in == nil {
		return nil
	}
	out := new(DeploymentDestination)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DestinationToWatch) DeepCopyInto(out *DestinationToWatch) {
	*out = *in
	if in.ExternalSecret != nil {
		in, out := &in.ExternalSecret, &out.ExternalSecret
		*out = new(ExternalSecretDestination)
		(*in).DeepCopyInto(*out)
	}
	if in.PushSecret != nil {
		in, out := &in.PushSecret, &out.PushSecret
		*out = new(PushSecretDestination)
		(*in).DeepCopyInto(*out)
	}
	if in.Deployment != nil {
		in, out := &in.Deployment, &out.Deployment
		*out = new(DeploymentDestination)
		(*in).DeepCopyInto(*out)
	}
	if in.UpdateStrategy != nil {
		in, out := &in.UpdateStrategy, &out.UpdateStrategy
		*out = new(UpdateStrategy)
		(*in).DeepCopyInto(*out)
	}
	if in.MatchStrategy != nil {
		in, out := &in.MatchStrategy, &out.MatchStrategy
		*out = new(MatchStrategy)
		(*in).DeepCopyInto(*out)
	}
	if in.WaitStrategy != nil {
		in, out := &in.WaitStrategy, &out.WaitStrategy
		*out = new(WaitStrategy)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DestinationToWatch.
func (in *DestinationToWatch) DeepCopy() *DestinationToWatch {
	if in == nil {
		return nil
	}
	out := new(DestinationToWatch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalSecretDestination) DeepCopyInto(out *ExternalSecretDestination) {
	*out = *in
	if in.NamespaceSelectors != nil {
		in, out := &in.NamespaceSelectors, &out.NamespaceSelectors
		*out = make([]v1.LabelSelector, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.LabelSelectors != nil {
		in, out := &in.LabelSelectors, &out.LabelSelectors
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Names != nil {
		in, out := &in.Names, &out.Names
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalSecretDestination.
func (in *ExternalSecretDestination) DeepCopy() *ExternalSecretDestination {
	if in == nil {
		return nil
	}
	out := new(ExternalSecretDestination)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCPSMAuthSecretRef) DeepCopyInto(out *GCPSMAuthSecretRef) {
	*out = *in
	out.SecretAccessKey = in.SecretAccessKey
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCPSMAuthSecretRef.
func (in *GCPSMAuthSecretRef) DeepCopy() *GCPSMAuthSecretRef {
	if in == nil {
		return nil
	}
	out := new(GCPSMAuthSecretRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCPWorkloadIdentity) DeepCopyInto(out *GCPWorkloadIdentity) {
	*out = *in
	in.ServiceAccountRef.DeepCopyInto(&out.ServiceAccountRef)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCPWorkloadIdentity.
func (in *GCPWorkloadIdentity) DeepCopy() *GCPWorkloadIdentity {
	if in == nil {
		return nil
	}
	out := new(GCPWorkloadIdentity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GooglePubSubAuth) DeepCopyInto(out *GooglePubSubAuth) {
	*out = *in
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(GCPSMAuthSecretRef)
		**out = **in
	}
	if in.WorkloadIdentity != nil {
		in, out := &in.WorkloadIdentity, &out.WorkloadIdentity
		*out = new(GCPWorkloadIdentity)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GooglePubSubAuth.
func (in *GooglePubSubAuth) DeepCopy() *GooglePubSubAuth {
	if in == nil {
		return nil
	}
	out := new(GooglePubSubAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GooglePubSubConfig) DeepCopyInto(out *GooglePubSubConfig) {
	*out = *in
	if in.Auth != nil {
		in, out := &in.Auth, &out.Auth
		*out = new(GooglePubSubAuth)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GooglePubSubConfig.
func (in *GooglePubSubConfig) DeepCopy() *GooglePubSubConfig {
	if in == nil {
		return nil
	}
	out := new(GooglePubSubConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HashicorpVaultConfig) DeepCopyInto(out *HashicorpVaultConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HashicorpVaultConfig.
func (in *HashicorpVaultConfig) DeepCopy() *HashicorpVaultConfig {
	if in == nil {
		return nil
	}
	out := new(HashicorpVaultConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeConfigRef) DeepCopyInto(out *KubeConfigRef) {
	*out = *in
	out.SecretRef = in.SecretRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeConfigRef.
func (in *KubeConfigRef) DeepCopy() *KubeConfigRef {
	if in == nil {
		return nil
	}
	out := new(KubeConfigRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubernetesAuth) DeepCopyInto(out *KubernetesAuth) {
	*out = *in
	if in.KubeConfigRef != nil {
		in, out := &in.KubeConfigRef, &out.KubeConfigRef
		*out = new(KubeConfigRef)
		**out = **in
	}
	if in.TokenRef != nil {
		in, out := &in.TokenRef, &out.TokenRef
		*out = new(TokenRef)
		**out = **in
	}
	if in.ServiceAccountRef != nil {
		in, out := &in.ServiceAccountRef, &out.ServiceAccountRef
		*out = new(ServiceAccountSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubernetesAuth.
func (in *KubernetesAuth) DeepCopy() *KubernetesAuth {
	if in == nil {
		return nil
	}
	out := new(KubernetesAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubernetesConfigMapConfig) DeepCopyInto(out *KubernetesConfigMapConfig) {
	*out = *in
	if in.Auth != nil {
		in, out := &in.Auth, &out.Auth
		*out = new(KubernetesAuth)
		(*in).DeepCopyInto(*out)
	}
	if in.LabelSelector != nil {
		in, out := &in.LabelSelector, &out.LabelSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubernetesConfigMapConfig.
func (in *KubernetesConfigMapConfig) DeepCopy() *KubernetesConfigMapConfig {
	if in == nil {
		return nil
	}
	out := new(KubernetesConfigMapConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubernetesObjectConfig) DeepCopyInto(out *KubernetesObjectConfig) {
	*out = *in
	if in.Auth != nil {
		in, out := &in.Auth, &out.Auth
		*out = new(KubernetesAuth)
		(*in).DeepCopyInto(*out)
	}
	if in.LabelSelector != nil {
		in, out := &in.LabelSelector, &out.LabelSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubernetesObjectConfig.
func (in *KubernetesObjectConfig) DeepCopy() *KubernetesObjectConfig {
	if in == nil {
		return nil
	}
	out := new(KubernetesObjectConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubernetesSecretConfig) DeepCopyInto(out *KubernetesSecretConfig) {
	*out = *in
	if in.Auth != nil {
		in, out := &in.Auth, &out.Auth
		*out = new(KubernetesAuth)
		(*in).DeepCopyInto(*out)
	}
	if in.LabelSelector != nil {
		in, out := &in.LabelSelector, &out.LabelSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubernetesSecretConfig.
func (in *KubernetesSecretConfig) DeepCopy() *KubernetesSecretConfig {
	if in == nil {
		return nil
	}
	out := new(KubernetesSecretConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MatchStrategy) DeepCopyInto(out *MatchStrategy) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]Condition, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchStrategy.
func (in *MatchStrategy) DeepCopy() *MatchStrategy {
	if in == nil {
		return nil
	}
	out := new(MatchStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MockConfig) DeepCopyInto(out *MockConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MockConfig.
func (in *MockConfig) DeepCopy() *MockConfig {
	if in == nil {
		return nil
	}
	out := new(MockConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NotificationSource) DeepCopyInto(out *NotificationSource) {
	*out = *in
	if in.AwsSqs != nil {
		in, out := &in.AwsSqs, &out.AwsSqs
		*out = new(AWSSQSConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.AzureEventGrid != nil {
		in, out := &in.AzureEventGrid, &out.AzureEventGrid
		*out = new(AzureEventGridConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.GooglePubSub != nil {
		in, out := &in.GooglePubSub, &out.GooglePubSub
		*out = new(GooglePubSubConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.Webhook != nil {
		in, out := &in.Webhook, &out.Webhook
		*out = new(WebhookConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.HashicorpVault != nil {
		in, out := &in.HashicorpVault, &out.HashicorpVault
		*out = new(HashicorpVaultConfig)
		**out = **in
	}
	if in.KubernetesSecret != nil {
		in, out := &in.KubernetesSecret, &out.KubernetesSecret
		*out = new(KubernetesSecretConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.KubernetesConfigMap != nil {
		in, out := &in.KubernetesConfigMap, &out.KubernetesConfigMap
		*out = new(KubernetesConfigMapConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.TCPSocket != nil {
		in, out := &in.TCPSocket, &out.TCPSocket
		*out = new(TCPSocketConfig)
		**out = **in
	}
	if in.Mock != nil {
		in, out := &in.Mock, &out.Mock
		*out = new(MockConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NotificationSource.
func (in *NotificationSource) DeepCopy() *NotificationSource {
	if in == nil {
		return nil
	}
	out := new(NotificationSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PatchOperationConfig) DeepCopyInto(out *PatchOperationConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PatchOperationConfig.
func (in *PatchOperationConfig) DeepCopy() *PatchOperationConfig {
	if in == nil {
		return nil
	}
	out := new(PatchOperationConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PushSecretDestination) DeepCopyInto(out *PushSecretDestination) {
	*out = *in
	if in.NamespaceSelectors != nil {
		in, out := &in.NamespaceSelectors, &out.NamespaceSelectors
		*out = make([]v1.LabelSelector, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.LabelSelectors != nil {
		in, out := &in.LabelSelectors, &out.LabelSelectors
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Names != nil {
		in, out := &in.Names, &out.Names
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PushSecretDestination.
func (in *PushSecretDestination) DeepCopy() *PushSecretDestination {
	if in == nil {
		return nil
	}
	out := new(PushSecretDestination)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RetryPolicy) DeepCopyInto(out *RetryPolicy) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RetryPolicy.
func (in *RetryPolicy) DeepCopy() *RetryPolicy {
	if in == nil {
		return nil
	}
	out := new(RetryPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretKeySelector) DeepCopyInto(out *SecretKeySelector) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretKeySelector.
func (in *SecretKeySelector) DeepCopy() *SecretKeySelector {
	if in == nil {
		return nil
	}
	out := new(SecretKeySelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceAccountSelector) DeepCopyInto(out *ServiceAccountSelector) {
	*out = *in
	if in.Audiences != nil {
		in, out := &in.Audiences, &out.Audiences
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceAccountSelector.
func (in *ServiceAccountSelector) DeepCopy() *ServiceAccountSelector {
	if in == nil {
		return nil
	}
	out := new(ServiceAccountSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TCPSocketConfig) DeepCopyInto(out *TCPSocketConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TCPSocketConfig.
func (in *TCPSocketConfig) DeepCopy() *TCPSocketConfig {
	if in == nil {
		return nil
	}
	out := new(TCPSocketConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TokenRef) DeepCopyInto(out *TokenRef) {
	*out = *in
	out.SecretRef = in.SecretRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TokenRef.
func (in *TokenRef) DeepCopy() *TokenRef {
	if in == nil {
		return nil
	}
	out := new(TokenRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UpdateStrategy) DeepCopyInto(out *UpdateStrategy) {
	*out = *in
	if in.PatchOperationConfig != nil {
		in, out := &in.PatchOperationConfig, &out.PatchOperationConfig
		*out = new(PatchOperationConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UpdateStrategy.
func (in *UpdateStrategy) DeepCopy() *UpdateStrategy {
	if in == nil {
		return nil
	}
	out := new(UpdateStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WaitForCondition) DeepCopyInto(out *WaitForCondition) {
	*out = *in
	if in.RetryTimeout != nil {
		in, out := &in.RetryTimeout, &out.RetryTimeout
		*out = new(v1.Duration)
		**out = **in
	}
	if in.MaxRetries != nil {
		in, out := &in.MaxRetries, &out.MaxRetries
		*out = new(int32)
		**out = **in
	}
	if in.TransitionedAfter != nil {
		in, out := &in.TransitionedAfter, &out.TransitionedAfter
		*out = new(v1.Duration)
		**out = **in
	}
	if in.UpdatedAfter != nil {
		in, out := &in.UpdatedAfter, &out.UpdatedAfter
		*out = new(v1.Duration)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WaitForCondition.
func (in *WaitForCondition) DeepCopy() *WaitForCondition {
	if in == nil {
		return nil
	}
	out := new(WaitForCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WaitStrategy) DeepCopyInto(out *WaitStrategy) {
	*out = *in
	if in.Time != nil {
		in, out := &in.Time, &out.Time
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Condition != nil {
		in, out := &in.Condition, &out.Condition
		*out = new(WaitForCondition)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WaitStrategy.
func (in *WaitStrategy) DeepCopy() *WaitStrategy {
	if in == nil {
		return nil
	}
	out := new(WaitStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookAuth) DeepCopyInto(out *WebhookAuth) {
	*out = *in
	if in.BasicAuth != nil {
		in, out := &in.BasicAuth, &out.BasicAuth
		*out = new(BasicAuth)
		**out = **in
	}
	if in.BearerToken != nil {
		in, out := &in.BearerToken, &out.BearerToken
		*out = new(BearerToken)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookAuth.
func (in *WebhookAuth) DeepCopy() *WebhookAuth {
	if in == nil {
		return nil
	}
	out := new(WebhookAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebhookConfig) DeepCopyInto(out *WebhookConfig) {
	*out = *in
	if in.Auth != nil {
		in, out := &in.Auth, &out.Auth
		*out = new(WebhookAuth)
		(*in).DeepCopyInto(*out)
	}
	if in.RetryPolicy != nil {
		in, out := &in.RetryPolicy, &out.RetryPolicy
		*out = new(RetryPolicy)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebhookConfig.
func (in *WebhookConfig) DeepCopy() *WebhookConfig {
	if in == nil {
		return nil
	}
	out := new(WebhookConfig)
	in.DeepCopyInto(out)
	return out
}
