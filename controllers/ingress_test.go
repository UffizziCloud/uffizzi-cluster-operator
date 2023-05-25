package controllers

import (
	"fmt"
	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"testing"
)

func TestBuildVClusterIngress(t *testing.T) {
	// Test inputs
	helmReleaseName := "test-release"
	uCluster := &v1alpha1.UffizziCluster{
		// Populate fields as necessary for your test
	}

	// Expected outputs
	expectedName := helmReleaseName + "-ingress"
	expectedNamespace := uCluster.Namespace
	expectedAnnotations := map[string]string{
		"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS",
		"nginx.ingress.kubernetes.io/ssl-redirect":     "true",
		"nginx.ingress.kubernetes.io/ssl-passthrough":  "true",
		"cert-manager.io/cluster-issuer":               "my-uffizzi-letsencrypt",
	}

	// Call the function being tested
	ingress := BuildVClusterIngress(helmReleaseName, uCluster)

	// Assert that the outputs match the expected results
	if ingress.ObjectMeta.Name != expectedName {
		t.Errorf("expected %v, got %v", expectedName, ingress.ObjectMeta.Name)
	}

	if ingress.ObjectMeta.Namespace != expectedNamespace {
		t.Errorf("expected %v, got %v", expectedNamespace, ingress.ObjectMeta.Namespace)
	}

	for k, v := range expectedAnnotations {
		if ingress.Annotations[k] != v {
			t.Errorf("expected %v for key %v, got %v", v, k, ingress.Annotations[k])
		}
	}
}

func TestBuildVClusterInternalServiceIngress(t *testing.T) {
	// Test inputs
	helmReleaseName := "test-release"
	service := v1alpha1.ExposedVClusterService{
		// Populate fields as necessary for your test
	}
	uCluster := &v1alpha1.UffizziCluster{
		// Populate fields as necessary for your test
	}

	// Expected outputs
	expectedName := fmt.Sprintf("%s-x-%s-x-%s", service.Name, service.Namespace, helmReleaseName)
	expectedNamespace := uCluster.Namespace
	expectedAnnotations := map[string]string{} // Assuming no annotations for this example

	// Call the function being tested
	ingress := BuildVClusterInternalServiceIngress(service, uCluster, helmReleaseName)

	// Assert that the outputs match the expected results
	if ingress.ObjectMeta.Name != expectedName {
		t.Errorf("expected %v, got %v", expectedName, ingress.ObjectMeta.Name)
	}

	if ingress.ObjectMeta.Namespace != expectedNamespace {
		t.Errorf("expected %v, got %v", expectedNamespace, ingress.ObjectMeta.Namespace)
	}

	for k, v := range expectedAnnotations {
		if ingress.Annotations[k] != v {
			t.Errorf("expected %v for key %v, got %v", v, k, ingress.Annotations[k])
		}
	}

	// If your test includes a service with CertManagerTLSEnabled set to true, you should also test that the TLS field is populated correctly.
	if service.CertManagerTLSEnabled {
		expectedSecretName := fmt.Sprintf("%s-x-%s-x-%s-x-tls", service.Name, service.Namespace, helmReleaseName)
		if len(ingress.Spec.TLS) != 1 || ingress.Spec.TLS[0].SecretName != expectedSecretName {
			t.Errorf("expected TLS SecretName %v, got %v", expectedSecretName, ingress.Spec.TLS[0].SecretName)
		}
	}
}

func TestBuildVClusterIngressHost(t *testing.T) {
	// Test inputs
	uCluster := &v1alpha1.UffizziCluster{
		ObjectMeta: v1.ObjectMeta{
			Name: "cluster1",
		},
		Spec: v1alpha1.UffizziClusterSpec{
			Ingress: v1alpha1.UffizziClusterIngress{
				Host: "test.com",
			},
		},
	}

	// Expected output
	expectedHost := "cluster1.uc.test.com"

	// Call the function being tested
	host := BuildVClusterIngressHost(uCluster)

	// Assert that the output matches the expected result
	if host != expectedHost {
		t.Errorf("expected %v, got %v", expectedHost, host)
	}
}

func TestBuildVClusterInternalServiceIngressHost(t *testing.T) {
	// Test inputs
	rand.Seed(1) // set seed for reproducible test, you can use any constant number
	uCluster := &v1alpha1.UffizziCluster{
		Spec: v1alpha1.UffizziClusterSpec{
			Ingress: v1alpha1.UffizziClusterIngress{
				Host: "test.com",
			},
		},
	}

	// Expected output
	// Since we're using a fixed seed, the random string will be the same for each test run
	expectedHost := "xn8fg-svc.uc.test.com"

	// Call the function being tested
	host := BuildVClusterInternalServiceIngressHost(uCluster)

	// Assert that the output matches the expected result
	if host != expectedHost {
		t.Errorf("expected %v, got %v", expectedHost, host)
	}
}
