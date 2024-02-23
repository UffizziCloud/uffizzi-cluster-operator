package etcd

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/controllers/helm/types"
)

type Etcd struct {
	Global         Global         `json:"global,omitempty"`
	ReplicaCount   int            `json:"replicaCount,omitempty"`
	ReadinessProbe ReadinessProbe `json:"readinessProbe,omitempty"`
	Persistence    Persistence    `json:"persistence,omitempty"`
	Tolerations    []Toleration   `json:"tolerations,omitempty"`
	NodeSelector   NodeSelector   `json:"nodeSelector,omitempty"`
	Auth           Auth           `json:"auth,omitempty"`
	Resources      Resources      `json:"resources,omitempty"`
}

type Global struct {
	StorageClass string `json:"storageClass,omitempty"`
}

type ReadinessProbe struct {
	InitialDelaySeconds int `json:"initialDelaySeconds,omitempty"`
	PeriodSeconds       int `json:"periodSeconds,omitempty"`
}

type Persistence struct {
	Size string `json:"size,omitempty"`
}

type Toleration struct {
	Effect   string `json:"effect,omitempty"`
	Key      string `json:"key,omitempty"`
	Operator string `json:"operator,omitempty"`
}

type NodeSelector struct {
	SandboxGKEIORuntime string `json:"sandbox.gke.io/runtime,omitempty"`
}

type Auth struct {
	Rbac Rbac `json:"rbac,omitempty"`
}

type Rbac struct {
	Create bool `json:"create"`
}

type Resources struct {
	Limits   types.ContainerMemoryCPU `json:"limits"`
	Requests types.ContainerMemoryCPU `json:"requests"`
}
