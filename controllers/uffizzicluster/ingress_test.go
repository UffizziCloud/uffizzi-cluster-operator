package uffizzicluster

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/controllers/helm/build/vcluster"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

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
	host := vcluster.BuildVClusterIngressHost(uCluster)

	// Assert that the output matches the expected result
	if host != expectedHost {
		t.Errorf("expected %v, got %v", expectedHost, host)
	}
}
