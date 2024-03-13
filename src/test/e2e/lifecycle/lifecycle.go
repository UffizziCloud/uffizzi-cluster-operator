package lifecycle

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/test/e2e"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type LifecycleTestDefinition struct {
	Name           string
	Spec           v1alpha1.UffizziClusterSpec
	ExpectedStatus ExpectedStatusThroughLifetime
	K8SClient      client.Client
}

type ExpectedStatusThroughLifetime struct {
	Initializing v1alpha1.UffizziClusterStatus
	Ready        v1alpha1.UffizziClusterStatus
	Sleeping     v1alpha1.UffizziClusterStatus
	Awoken       v1alpha1.UffizziClusterStatus
}

func deleteTestNamespace(name string) error {
	e2eObj := e2e.GetE2E()
	return e2e.GetE2E().K8SClient.Delete(e2eObj.Ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	})
}
