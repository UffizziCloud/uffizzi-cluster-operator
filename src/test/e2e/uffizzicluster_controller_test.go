package e2e

import (
	"context"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/test/util/resources"
	. "github.com/onsi/ginkgo/v2"
	v1 "k8s.io/api/core/v1"
)

type TestDefinition struct {
	Name string
	Spec v1alpha1.UffizziClusterSpec
}

func (td *TestDefinition) ExecLifecycleTest(ctx context.Context) {
	ns := resources.CreateTestNamespace(td.Name)
	uc := resources.CreateTestUffizziClusterWithSpec(td.Name, ns.Name, td.Spec)
	wrapUffizziClusterLifecycleTest(ctx, ns, uc)
}

const (
	timeout        = "5m"
	pollingTimeout = "100ms"
)

// basic clusters with no configuration

var _ = Describe("Basic Vanilla K3S UffizziCluster Lifecycle", func() {
	BeforeEach(func() {
		if e2e.IsTainted {
			Skip("Skipping test because cluster is tainted")
		}
	})
	ctx := context.Background()
	testUffizziCluster := TestDefinition{
		Name: "basic-test",
		Spec: v1alpha1.UffizziClusterSpec{},
	}
	testUffizziCluster.ExecLifecycleTest(ctx)
})

//
//var _ = Describe("Basic Vanilla K8S UffizziCluster Lifecycle", func() {
//	ctx := context.Background()
//	testUffizziCluster := TestDefinition{
//		Name: "basic-k8s-test",
//		Spec: v1alpha1.UffizziClusterSpec{
//			Distro: "k8s",
//		},
//	}
//	testUffizziCluster.ExecLifecycleTest(ctx, !e2e.IsTainted)
//})

// clusters with tainted nodes

var _ = Describe("UffizziCluster NodeSelector and Tolerations", func() {
	BeforeEach(func() {
		if !e2e.IsTainted {
			Skip("Skipping test because cluster is not tainted")
		}
	})
	ctx := context.Background()
	testUffizziCluster := TestDefinition{
		Name: "nodeselector-toleration-test",
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

// TODO: k3s with etcd

//var _ = Describe("Basic K3S UffizziCluster with ETCD Lifecycle", func() {
//	ctx := context.Background()
//	testUffizziCluster := TestDefinition{
//		Name: "basic-etcd",
//		Spec: v1alpha1.UffizziClusterSpec{
//			ExternalDatastore: constants.ETCD,
//		},
//	}
//	testUffizziCluster.ExecLifecycleTest(ctx)
//})
