package e2e

import (
	"context"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/constants"
	. "github.com/onsi/ginkgo/v2"
	v1 "k8s.io/api/core/v1"
)

var _ = Describe("k3s", func() {
	BeforeEach(func() {
		if e2e.IsTainted {
			Skip("Skipping test because cluster is tainted")
		}
	})
	ctx := context.Background()
	testUffizziCluster := TestDefinition{
		Name:           "basic-k3s-test",
		Spec:           v1alpha1.UffizziClusterSpec{},
		ExpectedStatus: initExpectedStatusOverLifetime(),
	}
	testUffizziCluster.Run(ctx)
})

var _ = Describe("Basic Vanilla K8S UffizziCluster Lifecycle", func() {
	BeforeEach(func() {
		if e2e.IsTainted {
			Skip("Skipping test because cluster is tainted")
		}
	})
	ctx := context.Background()
	testUffizziCluster := TestDefinition{
		Name: "basic-k8s-test",
		Spec: v1alpha1.UffizziClusterSpec{
			Distro: "k8s",
		},
		ExpectedStatus: initExpectedStatusOverLifetime(),
	}
	testUffizziCluster.Run(ctx)
})

var _ = Describe("Basic K3S UffizziCluster with ETCD Lifecycle", func() {
	BeforeEach(func() {
		if e2e.IsTainted {
			Skip("Skipping test because cluster is tainted")
		}
	})
	ctx := context.Background()
	testUffizziCluster := TestDefinition{
		Name: "k3s-etcd-test",
		Spec: v1alpha1.UffizziClusterSpec{
			ExternalDatastore: constants.ETCD,
		},
		ExpectedStatus: initExpectedStatusOverLifetime(),
	}
	testUffizziCluster.Run(ctx)
})

// Test against cluster with tainted nodes - good for testing node affinities

var _ = Describe("UffizziCluster NodeSelector and Tolerations", func() {
	BeforeEach(func() {
		if !e2e.IsTainted {
			Skip("Skipping test because cluster is not tainted")
		}
	})
	ctx := context.Background()
	testUffizziCluster := TestDefinition{
		Name: "k3s-nodeselector-toleration-test",
		Spec: v1alpha1.UffizziClusterSpec{
			NodeSelector: map[string]string{
				"testkey": "testvalue",
			},
			Toleration: []v1.Toleration{
				{
					Key:      "testkey",
					Operator: "Equal",
					Value:    "testvalue",
					Effect:   "NoSchedule",
				},
			},
		},
		ExpectedStatus: initExpectedStatusOverLifetime(),
	}
	testUffizziCluster.Run(ctx)
})
