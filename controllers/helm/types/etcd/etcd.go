package etcd

import "github.com/UffizziCloud/uffizzi-cluster-operator/controllers/helm/types"

type Etcd struct {
	Global         EtcdGlobal         `json:"global"`
	ReplicaCount   int                `json:"replicaCount"`
	ReadinessProbe EtcdReadinessProbe `json:"readinessProbe"`
	Persistence    EtcdPersistence    `json:"persistence"`
	Tolerations    []EtcdToleration   `json:"tolerations"`
	NodeSelector   EtcdNodeSelector   `json:"nodeSelector"`
	Auth           EtcdAuth           `json:"auth"`
	Resources      EtcdResources      `json:"resources"`
}

type EtcdGlobal struct {
	StorageClass string `json:"storageClass"`
}

type EtcdReadinessProbe struct {
	InitialDelaySeconds int `json:"initialDelaySeconds"`
	PeriodSeconds       int `json:"periodSeconds"`
}

type EtcdPersistence struct {
	Size string `json:"size"`
}

type EtcdToleration struct {
	Effect   string `json:"effect"`
	Key      string `json:"key"`
	Operator string `json:"operator"`
}

type EtcdNodeSelector struct {
	SandboxGKEIORuntime string `json:"sandbox.gke.io/runtime"`
}

type EtcdAuth struct {
	Rbac EtcdRbac `json:"rbac"`
}

type EtcdRbac struct {
	Create bool `json:"create"`
}

type EtcdResources struct {
	Limits   types.ContainerMemoryCPU `json:"limits"`
	Requests types.ContainerMemoryCPU `json:"requests"`
}
