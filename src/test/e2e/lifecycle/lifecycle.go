package lifecycle

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
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
