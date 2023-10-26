package etcd

import "github.com/UffizziCloud/uffizzi-cluster-operator/controllers/helm/types"

type Etcd struct {
	Global         Global         `json:"global"`
	ReplicaCount   int            `json:"replicaCount"`
	ReadinessProbe ReadinessProbe `json:"readinessProbe"`
	Persistence    Persistence    `json:"persistence"`
	Tolerations    []Toleration   `json:"tolerations"`
	NodeSelector   NodeSelector   `json:"nodeSelector"`
	Auth           Auth           `json:"auth"`
	Resources      Resources      `json:"resources"`
}

type Global struct {
	StorageClass string `json:"storageClass"`
}

type ReadinessProbe struct {
	InitialDelaySeconds int `json:"initialDelaySeconds"`
	PeriodSeconds       int `json:"periodSeconds"`
}

type Persistence struct {
	Size string `json:"size"`
}

type Toleration struct {
	Effect   string `json:"effect"`
	Key      string `json:"key"`
	Operator string `json:"operator"`
}

type NodeSelector struct {
	SandboxGKEIORuntime string `json:"sandbox.gke.io/runtime"`
}

type Auth struct {
	Rbac Rbac `json:"rbac"`
}

type Rbac struct {
	Create bool `json:"create"`
}

type Resources struct {
	Limits   types.ContainerMemoryCPU `json:"limits"`
	Requests types.ContainerMemoryCPU `json:"requests"`
}
