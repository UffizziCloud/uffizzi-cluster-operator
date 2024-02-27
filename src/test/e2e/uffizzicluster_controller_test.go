package e2e

import (
	"context"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/test/util/resources"
	. "github.com/onsi/ginkgo/v2"
)

type TestDefinition struct {
	Name string
	Spec v1alpha1.UffizziClusterSpec
}

func (td *TestDefinition) ExecLifecycleTest(ctx context.Context) {
	ns := resources.CreateTestNamespace(td.Name)
	uc := resources.CreateTestUffizziCluster(td.Name, ns.Name)
	wrapUffizziClusterLifecycleTest(ctx, ns, uc)
}

const (
	timeout        = "1m"
	pollingTimeout = "100ms"
)

var _ = Describe("Basic K3S UffizziCluster Lifecycle", func() {
	ctx := context.Background()
	testUffizziCluster := TestDefinition{
		Name: "basic",
		Spec: v1alpha1.UffizziClusterSpec{},
	}
	testUffizziCluster.ExecLifecycleTest(ctx)
})

//
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

var _ = Describe("Basic Vanila K8S UffizziCluster Lifecycle", func() {
	ctx := context.Background()
	testUffizziCluster := TestDefinition{
		Name: "basic-k8s",
		Spec: v1alpha1.UffizziClusterSpec{
			Distro: "k8s",
		},
	}
	testUffizziCluster.ExecLifecycleTest(ctx)
})
