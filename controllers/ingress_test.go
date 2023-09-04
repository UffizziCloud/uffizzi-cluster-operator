package controllers

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	expectedAnnotations := map[string]string{}

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
	expectedHost := "cluster1-test.com"

	// Call the function being tested
	host := BuildVClusterIngressHost(uCluster)

	// Assert that the output matches the expected result
	if host != expectedHost {
		t.Errorf("expected %v, got %v", expectedHost, host)
	}
}
