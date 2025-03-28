package gcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	iam "cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	relv1 "github.com/external-secrets-inc/reloader/api/v1alpha1"
	"github.com/googleapis/gax-go/v2"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc.go4.org/credentials/oauth"
	authenticationv1 "k8s.io/api/authentication/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	ctrlcfg "sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	gcpSAAnnotation = "iam.gke.io/gcp-service-account"
)

// workloadIdentity holds all clients and generators needed
// to create a gcp oauth token.
type workloadIdentity struct {
	iamClient            IamClient
	idBindTokenGenerator idBindTokenGenerator
	saTokenGenerator     saTokenGenerator
	clusterProjectID     string
}

// interface to GCP IAM API.
type IamClient interface {
	GenerateAccessToken(ctx context.Context, req *credentialspb.GenerateAccessTokenRequest, opts ...gax.CallOption) (*credentialspb.GenerateAccessTokenResponse, error)
	Close() error
}

// interface to securetoken/identitybindingtoken API.
type idBindTokenGenerator interface {
	Generate(context.Context, *http.Client, string, string, string) (*oauth2.Token, error)
}

// interface to kubernetes serviceaccount token request API.
type saTokenGenerator interface {
	Generate(context.Context, []string, string, string) (*authenticationv1.TokenRequest, error)
}

// newWorkloadIdentity creates a new workloadIdentity object.
//
// This function creates new instances of the IamClient and the
// saTokenGenerator. The IamClient is used to generate an access token
// for the service account, and the saTokenGenerator is used to generate
// a token request for the service account.
//
// The returned workloadIdentity object is used to create a TokenSource
// which can be used to authenticate with GCP.
func newWorkloadIdentity(ctx context.Context, projectID string) (*workloadIdentity, error) {
	satg, err := newSATokenGenerator()
	if err != nil {
		return nil, err
	}
	iamc, err := newIAMClient(ctx)
	if err != nil {
		return nil, err
	}
	return &workloadIdentity{
		iamClient:            iamc,
		idBindTokenGenerator: newIDBindTokenGenerator(),
		saTokenGenerator:     satg,
		clusterProjectID:     projectID,
	}, nil
}

func (w *workloadIdentity) TokenSource(ctx context.Context, auth *relv1.GooglePubSubAuth, kube kclient.Client) (oauth2.TokenSource, error) {
	wi := auth.WorkloadIdentity
	if wi == nil {
		return nil, nil
	}
	saKey := types.NamespacedName{
		Name:      wi.ServiceAccountRef.Name,
		Namespace: wi.ServiceAccountRef.Namespace,
	}

	sa := &v1.ServiceAccount{}
	err := kube.Get(ctx, saKey, sa)
	if err != nil {
		return nil, err
	}

	idProvider := fmt.Sprintf("https://container.googleapis.com/v1/projects/%s/locations/%s/clusters/%s",
		w.clusterProjectID,
		wi.ClusterLocation,
		wi.ClusterName)
	idPool := fmt.Sprintf("%s.svc.id.goog", w.clusterProjectID)
	audiences := []string{idPool}

	if len(wi.ServiceAccountRef.Audiences) > 0 {
		audiences = append(audiences, wi.ServiceAccountRef.Audiences...)
	}
	gcpSA := sa.Annotations[gcpSAAnnotation]
	resp, err := w.saTokenGenerator.Generate(ctx, audiences, saKey.Name, saKey.Namespace)
	if err != nil {
		return nil, err
	}
	idBindToken, err := w.idBindTokenGenerator.Generate(ctx, http.DefaultClient, resp.Status.Token, idPool, idProvider)
	if err != nil {
		return nil, err
	}

	// If no `iam.gke.io/gcp-service-account` annotation is present the
	// identitybindingtoken will be used directly, allowing bindings on secrets
	// of the form "serviceAccount:<project>.svc.id.goog[<namespace>/<sa>]".
	if gcpSA == "" {
		return oauth2.StaticTokenSource(idBindToken), nil
	}
	gcpSAResp, err := w.iamClient.GenerateAccessToken(ctx, &credentialspb.GenerateAccessTokenRequest{
		Name:  fmt.Sprintf("projects/-/serviceAccounts/%s", gcpSA),
		Scope: secretmanager.DefaultAuthScopes(),
	}, gax.WithGRPCOptions(grpc.PerRPCCredentials(oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(idBindToken)})))
	if err != nil {
		return nil, err
	}

	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: gcpSAResp.GetAccessToken(),
	}), nil
}

func (w *workloadIdentity) Close() error {
	if w.iamClient != nil {
		return w.iamClient.Close()
	}
	return nil
}

func newIAMClient(ctx context.Context) (IamClient, error) {
	iamOpts := []option.ClientOption{
		option.WithUserAgent("external-secrets-operator"),
		// tell the secretmanager library to not add transport-level ADC since
		// we need to override on a per call basis
		option.WithoutAuthentication(),
		// grpc oauth TokenSource credentials require transport security, so
		// this must be set explicitly even though TLS is used
		option.WithGRPCDialOption(grpc.WithTransportCredentials(credentials.NewTLS(nil))),
		option.WithGRPCConnectionPool(5),
	}
	return iam.NewIamCredentialsClient(ctx, iamOpts...)
}

type k8sSATokenGenerator struct {
	corev1 clientcorev1.CoreV1Interface
}

func (g *k8sSATokenGenerator) Generate(ctx context.Context, audiences []string, name, namespace string) (*authenticationv1.TokenRequest, error) {
	// Request a serviceaccount token for the pod
	ttl := int64((15 * time.Minute).Seconds())
	return g.corev1.
		ServiceAccounts(namespace).
		CreateToken(ctx, name,
			&authenticationv1.TokenRequest{
				Spec: authenticationv1.TokenRequestSpec{
					ExpirationSeconds: &ttl,
					Audiences:         audiences,
				},
			},
			metav1.CreateOptions{},
		)
}

func newSATokenGenerator() (saTokenGenerator, error) {
	cfg, err := ctrlcfg.GetConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &k8sSATokenGenerator{
		corev1: clientset.CoreV1(),
	}, nil
}

// Trades the kubernetes token for an identitybindingtoken token.
type gcpIDBindTokenGenerator struct {
	targetURL string
}

func newIDBindTokenGenerator() idBindTokenGenerator {
	return &gcpIDBindTokenGenerator{
		targetURL: "https://securetoken.googleapis.com/v1/identitybindingtoken",
	}
}

func (g *gcpIDBindTokenGenerator) Generate(ctx context.Context, client *http.Client, k8sToken, idPool, idProvider string) (*oauth2.Token, error) {
	body, err := json.Marshal(map[string]string{
		"grant_type":           "urn:ietf:params:oauth:grant-type:token-exchange",
		"subject_token_type":   "urn:ietf:params:oauth:token-type:jwt",
		"requested_token_type": "urn:ietf:params:oauth:token-type:access_token",
		"subject_token":        k8sToken,
		"audience":             fmt.Sprintf("identitynamespace:%s:%s", idPool, idProvider),
		"scope":                "https://www.googleapis.com/auth/cloud-platform",
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", g.targetURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get idbindtoken token, status: %v", resp.StatusCode)
	}

	defer resp.Body.Close() //nolint
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	idBindToken := &oauth2.Token{}
	if err := json.Unmarshal(respBody, idBindToken); err != nil {
		return nil, err
	}
	return idBindToken, nil
}
