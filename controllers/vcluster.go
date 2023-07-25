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
	Ingresses EnabledBool `json:"ingresses,omitempty"`
}

type EnabledBool struct {
	Enabled bool `json:"enabled"`
}

type VClusterResourceQuota struct {
	Enabled bool                           `json:"enabled"`
	Quota   VClusterResourceQuotaDefiniton `json:"quota"`
}

type VClusterResourceQuotaDefiniton struct {
	RequestsCpu                 string `json:"requests.cpu"`
	RequestsMemory              string `json:"requests.memory"`
	RequestsStorage             string `json:"requests.storage"`
	RequestsEphemeralStorage    string `json:"requests.ephemeral-storage"`
	LimitsCpu                   string `json:"limits.cpu"`
	LimitsMemory                string `json:"limits.memory"`
	LimitsEphemeralStorage      string `json:"limits.ephemeral-storage"`
	ServicesLoadbalancers       int    `json:"services.loadbalancers"`
	ServicesNodePorts           int    `json:"services.nodeports"`
	CountEndpoints              int    `json:"count/endpoints"`
	CountPods                   int    `json:"count/pods"`
	CountServices               int    `json:"count/services"`
	CountSecrets                int    `json:"count/secrets"`
	CountConfigmaps             int    `json:"count/configmaps"`
	CountPersistentVolumeClaims int    `json:"count/persistentvolumeclaims"`
}

type LimitRangeResources struct {
	EphemeralStorage string `json:"ephemeral-storage"`
	Memory           string `json:"memory"`
	Cpu              string `json:"cpu"`
}

type VClusterLimitRange struct {
	Enabled        bool                `json:"enabled"`
	Default        LimitRangeResources `json:"default"`
	DefaultRequest LimitRangeResources `json:"defaultRequest"`
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
	LimitRange          VClusterLimitRange    `json:"limitRange"`
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
