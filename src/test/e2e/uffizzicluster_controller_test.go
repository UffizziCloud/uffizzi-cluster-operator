package e2e

import (
	"context"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/constants"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/helm/types/vcluster"
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
		ExpectedStatus: initExpectedStatusOverLifetime(),
	}
	testUffizziCluster.Run(ctx)
})

var _ = Describe("k3s: w/ etcd", func() {
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
		ExpectedStatus: initExpectedStatusOverLifetime(),
	}
	testUffizziCluster.Run(ctx)
})

// Test against cluster with tainted nodes - good for testing node affinities

var _ = Describe("k3s: explicit nodeselector and toleration", func() {
	BeforeEach(func() {
		if !e2e.IsTainted {
			Skip("Skipping test because cluster is not tainted")
		}
	})
	ctx := context.Background()
	tolerations := append([]v1.Toleration{}, vcluster.GvisorToleration.ToV1())
	testUffizziCluster := TestDefinition{
		Name: "k3s-nds-tlrtn",
		Spec: v1alpha1.UffizziClusterSpec{
			NodeSelector: vcluster.GvisorNodeSelector,
			Toleration:   tolerations,
		},
		ExpectedStatus: initExpectedStatusOverLifetime(),
	}
	testUffizziCluster.ExpectedStatus.Ready.NodeSelector = vcluster.GvisorNodeSelector
	testUffizziCluster.ExpectedStatus.Ready.Tolerations = tolerations
	testUffizziCluster.Run(ctx)
})

var _ = Describe("k3s: nodeselector template - gvisor", func() {
	BeforeEach(func() {
		if !e2e.IsTainted {
			Skip("Skipping test because cluster is not tainted")
		}
	})
	ctx := context.Background()
	tolerations := append([]v1.Toleration{}, vcluster.GvisorToleration.ToV1())
	testUffizziCluster := TestDefinition{
		Name: "k3s-nds-template-gvisor",
		Spec: v1alpha1.UffizziClusterSpec{
			NodeSelectorTemplate: constants.GVISOR,
		},
		ExpectedStatus: initExpectedStatusOverLifetime(),
	}
	testUffizziCluster.ExpectedStatus.Ready.NodeSelector = vcluster.GvisorNodeSelector
	testUffizziCluster.ExpectedStatus.Ready.Tolerations = tolerations
	testUffizziCluster.Run(ctx)
})
