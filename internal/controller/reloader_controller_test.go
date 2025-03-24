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

package controller

import (
	"context"
	"fmt"

	"github.com/external-secrets-inc/reloader/internal/events"
	"github.com/external-secrets-inc/reloader/internal/handler"
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	esov1 "github.com/external-secrets-inc/reloader/api/v1alpha1"
	"github.com/external-secrets-inc/reloader/internal/listener"
)

var _ = Describe("Reloader Controller", func() {
	var (
		ctx            context.Context
		scheme         *runtime.Scheme
		fakeClient     client.Client
		reconciler     *ReloaderReconciler
		config         *esov1.Config
		externalSecret *esv1beta1.ExternalSecret
		eventChan      chan events.SecretRotationEvent
		mockFactory    *listener.MockListenerFactory
	)

	BeforeEach(func() {
		ctx = context.Background()
		scheme = runtime.NewScheme()
		Expect(esov1.AddToScheme(scheme)).To(Succeed())
		Expect(esv1beta1.AddToScheme(scheme)).To(Succeed())

		fakeClient = fake.NewClientBuilder().WithScheme(scheme).Build()

		eventChan = make(chan events.SecretRotationEvent, 10)
		mockFactory = listener.NewMockListenerFactory()
		manager := listener.NewListenerManager(ctx, mockFactory, eventChan, log.FromContext(ctx))
		eventHandler := handler.NewEventHandler(fakeClient)

		reconciler = &ReloaderReconciler{
			Client:          fakeClient,
			Scheme:          scheme,
			listenerFactory: mockFactory,
			listenerManager: manager,
			eventChan:       eventChan,
			eventHandler:    eventHandler,
		}

		go reconciler.processEvents(ctx)

		config = &esov1.Config{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-reloader",
				Namespace: "default",
			},
			Spec: esov1.ConfigSpec{
				NotificationSources: []esov1.NotificationSource{
					{
						Type: "Mock",
						Mock: &esov1.MockConfig{EmitInterval: 1000},
					},
				},
				DestinationsToWatch: []esov1.DestinationToWatch{
					{
						ExternalSecret: &esov1.ExternalSecretDestination{
							Names: []string{
								"test-external-secret-data",
								"test-external-secret-datafrom-extract",
								"test-external-secret-datafrom-find",
							},
						},
					},
				},
			},
		}
		Expect(fakeClient.Create(context.Background(), config)).To(Succeed())

		// Reconcile the controller to set up the listeners
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      config.Name,
				Namespace: config.Namespace,
			},
		})
		Expect(err).NotTo(HaveOccurred())
		// Reconcile the controller to set up the listeners
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      config.Name,
				Namespace: config.Namespace,
			},
		})
		Expect(err).NotTo(HaveOccurred())
	})

	Context("When a config is created/updated/deleted", func() {
		It("should add the processed annotation and push a Created/Updated/Deleted event for Reloader", func() {
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      config.Name,
					Namespace: config.Namespace,
				},
			}

			// Get update async rotator to check ProcessedAnnotation
			updatedconfig := &esov1.Config{}
			Expect(fakeClient.Get(ctx, req.NamespacedName, updatedconfig)).To(Succeed())
			Expect(updatedconfig.Annotations).To(HaveKey(ProcessedAnnotation))

			// Call reconcile to generate updated event
			_, err := reconciler.Reconcile(ctx, req)
			Expect(err).NotTo(HaveOccurred())

			// Deleting async rotator manifest to generate deleted event
			Expect(fakeClient.Delete(context.Background(), config)).To(Succeed())

			_, err = reconciler.Reconcile(ctx, req)
			Expect(err).NotTo(HaveOccurred())

		})
	})

	Context("When a secret rotation event is received, and the secret is not watched", func() {
		It("should not annotate any event out of the secrets to watch list", func() {
			// Create an ExternalSecret that references the secret not watched
			externalSecret = &esv1beta1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-not-watched",
					Namespace: "default",
				},
				Spec: esv1beta1.ExternalSecretSpec{
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "my-secret-store",
						Kind: "SecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "password",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key: "aws://secret/arn:aws:secretsmanager:us-east-1:123456789012:secret:mysecret",
							},
						},
					},
				},
			}

			Expect(fakeClient.Create(context.Background(), externalSecret)).To(Succeed())
			assertNotWatchedAnnotations(fakeClient, "test-not-watched")
		})
	})

	Context("When a secret rotation event is received", func() {
		It("should annotate the corresponding ExternalSecret using data field", func() {
			// Create an ExternalSecret that references the secret via data field
			externalSecret = &esv1beta1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-external-secret-data",
					Namespace: "default",
				},
				Spec: esv1beta1.ExternalSecretSpec{
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "my-secret-store",
						Kind: "SecretStore",
					},
					Data: []esv1beta1.ExternalSecretData{
						{
							SecretKey: "password",
							RemoteRef: esv1beta1.ExternalSecretDataRemoteRef{
								Key: "aws://secret/arn:aws:secretsmanager:us-east-1:123456789012:secret:mysecret",
							},
						},
					},
				},
			}

			Expect(fakeClient.Create(context.Background(), externalSecret)).To(Succeed())
			assertAnnotations(fakeClient, "test-external-secret-data")
		})
	})

	Context("When a secret rotation event is received and ExternalSecret uses dataFrom.extract", func() {
		It("should annotate the corresponding ExternalSecret using dataFrom.extract", func() {
			secretIdentifier := "aws://secret/arn:aws:secretsmanager:us-east-1:123456789012:secret:mysecret"

			// Create an ExternalSecret that references the secret via dataFrom.extract
			externalSecret = &esv1beta1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-external-secret-datafrom-extract",
					Namespace: "default",
				},
				Spec: esv1beta1.ExternalSecretSpec{
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "my-secret-store",
						Kind: "SecretStore",
					},
					DataFrom: []esv1beta1.ExternalSecretDataFromRemoteRef{
						{
							Extract: &esv1beta1.ExternalSecretDataRemoteRef{
								Key: secretIdentifier,
							},
						},
					},
				},
			}

			Expect(fakeClient.Create(context.Background(), externalSecret)).To(Succeed())
			assertAnnotations(fakeClient, "test-external-secret-datafrom-extract")
		})
	})

	Context("When a secret rotation event is received and ExternalSecret uses dataFrom.find", func() {
		It("should annotate the corresponding ExternalSecret using dataFrom.find", func() {
			secretIdentifier := "aws://secret/arn:aws:secretsmanager:us-east-1:123456789012:secret:mysecret"

			// Create an ExternalSecret that references the secret via dataFrom.find
			externalSecret = &esv1beta1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-external-secret-datafrom-find",
					Namespace: "default",
				},
				Spec: esv1beta1.ExternalSecretSpec{
					SecretStoreRef: esv1beta1.SecretStoreRef{
						Name: "my-secret-store",
						Kind: "SecretStore",
					},
					DataFrom: []esv1beta1.ExternalSecretDataFromRemoteRef{
						{
							Find: &esv1beta1.ExternalSecretFind{
								Name: &esv1beta1.FindName{
									RegExp: secretIdentifier,
								},
							},
						},
					},
				},
			}

			Expect(fakeClient.Create(context.Background(), externalSecret)).To(Succeed())
			assertAnnotations(fakeClient, "test-external-secret-datafrom-find")
		})
	})
})

func assertAnnotations(fakeClient client.Client, secretName string) {
	updatedES := &esv1beta1.ExternalSecret{}
	key := types.NamespacedName{
		Namespace: "default",
		Name:      secretName,
	}
	// Wait for the controller to process the event by polling
	Eventually(func() error {
		updatedES = &esv1beta1.ExternalSecret{}
		err := fakeClient.Get(context.Background(), key, updatedES)
		if err != nil {
			return err
		}
		annotations := updatedES.GetAnnotations()
		if annotations == nil {
			return fmt.Errorf("ExternalSecret annotations should not be nil")
		}
		if annotations["async-rotation/last-rotated"] != "2024-09-19T12:00:00Z" {
			return fmt.Errorf("async-rotation/last-rotated annotation should be set to 2024-09-19T12:00:00Z")
		}
		if annotations["async-rotation/trigger-source"] != "aws-secretsmanager" {
			return fmt.Errorf("async-rotation/trigger-source annotation should be set to aws-secretsmanager")
		}
		return nil
	}, "5s", "500ms").Should(Succeed())
}

func assertNotWatchedAnnotations(fakeClient client.Client, secretName string) {
	updatedES := &esv1beta1.ExternalSecret{}
	key := types.NamespacedName{
		Namespace: "default",
		Name:      secretName,
	}
	// Wait for the controller to process the event by polling
	Eventually(func() error {
		updatedES = &esv1beta1.ExternalSecret{}
		err := fakeClient.Get(context.Background(), key, updatedES)
		if err != nil {
			return err
		}
		annotations := updatedES.GetAnnotations()
		if annotations != nil {
			return fmt.Errorf("ExternalSecret annotations should not be nil")
		}
		return nil
	}, "5s", "500ms").Should(Succeed())
}
