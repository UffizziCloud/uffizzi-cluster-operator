package controllers

import "github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"

type VCluster struct {
	Init            VClusterInit            `json:"init,omitempty"`
	Syncer          VClusterSyncer          `json:"syncer,omitempty"`
	Sync            VClusterSync            `json:"sync,omitempty"`
	Ingress         EnabledBool             `json:"ingress,omitempty"`
	FsGroup         int64                   `json:"fsgroup,omitempty"`
	Isolation       VClusterIsolation       `json:"isolation,omitempty"`
	NodeSelector    VClusterNodeSelector    `json:"nodeSelector,omitempty"`
	SecurityContext VClusterSecurityContext `json:"securityContext,omitempty"`
	Tolerations     []VClusterToleration    `json:"tolerations,omitempty"`
	MapServices     VClusterMapServices     `json:"mapServices,omitempty"`
}

// VClusterInit - resources which are created during the init phase of the vcluster
type VClusterInit struct {
	Manifests string               `json:"manifests"`
	Helm      []v1alpha1.HelmChart `json:"helm"`
}

// VClusterSyncer - parameters to create the syncer with
// https://www.vcluster.com/docs/architecture/basics#vcluster-syncer
type VClusterSyncer struct {
	KubeConfigContextName string   `json:"kubeConfigContextName"`
	ExtraArgs             []string `json:"extraArgs"`
}

type VClusterSync struct {
	Ingresses              EnabledBool `json:"ingresses,omitempty"`
	PersistentVolumes      EnabledBool `json:"persistentvolumes,omitempty"`
	PersistentVolumeClaims EnabledBool `json:"persistentvolumeclaims,omitempty"`
	StorageClasses         EnabledBool `json:"storageclasses,omitempty"`
}

type EnabledBool struct {
	Enabled bool `json:"enabled"`
}

type VClusterResourceQuota struct {
	Quota VClusterResourceQuotaDefiniton `json:"quota"`
}

type VClusterResourceQuotaDefiniton struct {
	ServicesNodePorts int `json:"services.nodeports"`
}

type VClusterMapServicesFromVirtual struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type VClusterMapServices struct {
	FromVirtual []VClusterMapServicesFromVirtual `json:"fromVirtual"`
}

// VClusterIsolation - parameters to define the isolation of the cluster
type VClusterIsolation struct {
	Enabled             bool                  `json:"enabled"`
	PodSecurityStandard string                `json:"podSecurityStandard"`
	ResourceQuota       VClusterResourceQuota `json:"resourceQuota"`
}

// VClusterNodeSelector - parameters to define the node selector of the cluster
type VClusterNodeSelector struct {
	SandboxGKEIORuntime string `json:"sandbox.gke.io/runtime"`
}

type VClusterSecurityContextCapabilities struct {
	Drop []string `json:"drop"`
}

// VClusterSecurityContext - parameters to define the security context of the cluster
type VClusterSecurityContext struct {
	Capabilities           VClusterSecurityContextCapabilities `json:"capabilities"`
	ReadOnlyRootFilesystem bool                                `json:"readOnlyRootFilesystem"`
	RunAsNonRoot           bool                                `json:"runAsNonRoot"`
	RunAsUser              int64                               `json:"runAsUser"`
}

type VClusterToleration struct {
	Effect   string `json:"effect"`
	Key      string `json:"key"`
	Operator string `json:"operator"`
}

func BuildVClusterHelmReleaseName(uCluster *v1alpha1.UffizziCluster) string {
	helmReleaseName := UCLUSTER_NAME_PREFIX + uCluster.Name
	return helmReleaseName
}
