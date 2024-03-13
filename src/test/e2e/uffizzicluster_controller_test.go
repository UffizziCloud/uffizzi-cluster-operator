package e2e

import (
	"context"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/constants"
	. "github.com/onsi/ginkgo/v2"
	v1 "k8s.io/api/core/v1"
)

// Tests against k3s clusters

var _ = Describe("k3s", func() {
	BeforeEach(func() {
		if e2e.IsTainted {
			Skip("Skipping test because cluster is tainted")
		}
	})
	ctx := context.Background()
	testUffizziCluster := TestDefinition{
		Name: "k3s",
		Spec: v1alpha1.UffizziClusterSpec{},
	}
	testUffizziCluster.ExecLifecycleTest(ctx)
})

var _ = Describe("k3s: without persistence", func() {
	BeforeEach(func() {
		if e2e.IsTainted {
			Skip("Skipping test because cluster is tainted")
		}
	})
	ctx := context.Background()
	testUffizziCluster := TestDefinition{
		Name: "k3s-storage-non-persistent",
		Spec: v1alpha1.UffizziClusterSpec{
			Storage: &v1alpha1.UffizziClusterStorage{
				Persistence: false,
			},
		},
	}
	testUffizziCluster.ExecLifecycleTest(ctx)
})

var _ = Describe("k3s: with persistence", func() {
	BeforeEach(func() {
		if e2e.IsTainted {
			Skip("Skipping test because cluster is tainted")
		}
	})
	ctx := context.Background()
	testUffizziCluster := TestDefinition{
		Name: "k3s-storage-persistent",
		Spec: v1alpha1.UffizziClusterSpec{
			Storage: &v1alpha1.UffizziClusterStorage{
				Persistence: true,
				// test size - 5Gi is the default
				Size: "2Gi",
			},
		},
	}
	testUffizziCluster.ExecLifecycleTest(ctx)
})

var _ = Describe("k3s: with etcd", func() {
	BeforeEach(func() {
		if e2e.IsTainted {
			Skip("Skipping test because cluster is tainted")
		}
	})
	ctx := context.Background()
	testUffizziCluster := TestDefinition{
		Name: "k3s-etcd",
		Spec: v1alpha1.UffizziClusterSpec{
			ExternalDatastore: constants.ETCD,
		},
	}
	testUffizziCluster.ExecLifecycleTest(ctx)
})

// Test against cluster with tainted nodes - good for testing node affinities

var _ = Describe("k3s: nodeselector and tolerations", func() {
	BeforeEach(func() {
		if !e2e.IsTainted {
			Skip("Skipping test because cluster is not tainted")
		}
	})
	ctx := context.Background()
	testUffizziCluster := TestDefinition{
		Name: "k3s-nodeselector-tolerations",
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
	}
	testUffizziCluster.ExecLifecycleTest(ctx)
})

var _ = Describe("k8s", func() {
	BeforeEach(func() {
		if e2e.IsTainted {
			Skip("Skipping test because cluster is tainted")
		}
	})
	ctx := context.Background()
	testUffizziCluster := TestDefinition{
		Name: "k8s",
		Spec: v1alpha1.UffizziClusterSpec{
			Distro: "k8s",
		},
	}
	testUffizziCluster.ExecLifecycleTest(ctx)
})
