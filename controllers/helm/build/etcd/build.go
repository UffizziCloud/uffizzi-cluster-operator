package etcd

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/controllers/constants"
	helmtypes "github.com/UffizziCloud/uffizzi-cluster-operator/controllers/helm/types"
	etcdhelmtypes "github.com/UffizziCloud/uffizzi-cluster-operator/controllers/helm/types/etcd"
)

func BuildETCDHelmValues() etcdhelmtypes.Etcd {
	return etcdhelmtypes.Etcd{
		Global: etcdhelmtypes.Global{
			StorageClass: constants.PREMIUM_RWO_STORAGE_CLASS,
		},
		ReplicaCount: 1,
		ReadinessProbe: etcdhelmtypes.ReadinessProbe{
			InitialDelaySeconds: 30,
			PeriodSeconds:       5,
		},
		Persistence: etcdhelmtypes.Persistence{
			Size: "10Gi",
		},
		Tolerations: []etcdhelmtypes.Toleration{
			{
				Effect:   constants.NO_SCHEDULE,
				Key:      constants.SANDBOX_GKE_IO_RUNTIME,
				Operator: constants.EXISTS,
			},
		},
		NodeSelector: etcdhelmtypes.NodeSelector{
			SandboxGKEIORuntime: constants.GVISOR,
		},
		Auth: etcdhelmtypes.Auth{
			Rbac: etcdhelmtypes.Rbac{
				Create: false,
			},
		},
		Resources: etcdhelmtypes.Resources{
			Limits: helmtypes.ContainerMemoryCPU{
				CPU:    "500m",
				Memory: "800Mi",
			},
			Requests: helmtypes.ContainerMemoryCPU{
				CPU:    "100m",
				Memory: "200Mi",
			},
		},
	}
}
