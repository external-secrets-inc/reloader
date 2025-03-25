package aws_test

import (
	"context"
	"testing"

	"github.com/external-secrets-inc/providers-listeners/mocks"
	"github.com/external-secrets-inc/providers-listeners/pkg"
	authAWS "github.com/external-secrets-inc/providers-listeners/pkg/auth/aws"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAWSSQSConfig_StaticAuth(t *testing.T) {
	ctx := context.Background()
	logger := testr.New(t) // Mock logger
	k8sClient := mocks.CreateFakeK8sClient().Build()

	authConfig := pkg.AWSSDKAuth{
		AuthMethod: authAWS.AuthMethodStatic,
		Region:     "us-east-1",
		SecretRef: &pkg.AWSSDKSecretRef{
			AccessKeyId: pkg.SecretKeySelector{
				Name:      "aws-secret",
				Key:       "access-key-id",
				Namespace: "default",
			},
			SecretAccessKey: pkg.SecretKeySelector{
				Name:      "aws-secret",
				Key:       "secret-access-key",
				Namespace: "default",
			},
		},
	}

	awsConfig, err := authAWS.CreateAWSSDKConfig(ctx, k8sClient, authConfig, logger)
	require.NoError(t, err)
	assert.NotEmpty(t, awsConfig.Region)
}

func TestCreateAWSSQSConfig_IRSA(t *testing.T) {
	ctx := context.Background()
	logger := testr.New(t)
	k8sClient := mocks.CreateFakeK8sClient().Build()

	authConfig := pkg.AWSSDKAuth{
		AuthMethod: authAWS.AuthMethodIRSA,
		Region:     "us-east-1",
		ServiceAccount: &pkg.ServiceAccountSelector{
			Name:      "aws-sa",
			Namespace: "default",
		},
	}

	awsConfig, err := authAWS.CreateAWSSDKConfig(ctx, k8sClient, authConfig, logger)
	require.NoError(t, err)
	assert.NotEmpty(t, awsConfig.Region)
}

func TestCreateAWSSQSConfig_InvalidAuthMethod(t *testing.T) {
	ctx := context.Background()
	logger := testr.New(t)
	k8sClient := mocks.CreateFakeK8sClient().Build()

	authConfig := pkg.AWSSDKAuth{
		AuthMethod: "invalid-auth",
		Region:     "us-east-1",
	}

	_, err := authAWS.CreateAWSSDKConfig(ctx, k8sClient, authConfig, logger)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported authentication method")
}
