package vcluster

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/controllers/helm/types"
)

type Common struct {
	Init            Init            `json:"init,omitempty"`
	Syncer          Syncer          `json:"syncer,omitempty"`
	Sync            Sync            `json:"sync,omitempty"`
	Ingress         Ingress         `json:"ingress,omitempty"`
	FsGroup         int64           `json:"fsgroup,omitempty"`
	Isolation       Isolation       `json:"isolation,omitempty"`
	NodeSelector    NodeSelector    `json:"nodeSelector,omitempty"`
	SecurityContext SecurityContext `json:"securityContext,omitempty"`
	Tolerations     []Toleration    `json:"tolerations,omitempty"`
	MapServices     MapServices     `json:"mapServices,omitempty"`
	Plugin          Plugins         `json:"plugin,omitempty"`
	Storage         Storage         `json:"storage,omitempty"`
	EnableHA        bool            `json:"enableHA,omitempty"`
}

type K3S struct {
	VCluster K3SAPIServer `json:"vcluster,omitempty"`
	Common
}

type K8S struct {
	APIServer K8SAPIServer `json:"apiserver,omitempty"`
	Common
}

type K8SAPIServer struct {
	Image              string             `json:"image,omitempty"`
	ExtraArgs          []string           `json:"extraArgs,omitempty"`
	Replicas           int32              `json:"replicas,omitempty"`
	NodeSelector       NodeSelector       `json:"nodeSelector,omitempty"`
	Tolerations        []Toleration       `json:"tolerations,omitempty"`
	Labels             map[string]string  `json:"labels,omitempty"`
	Annotations        map[string]string  `json:"annotations,omitempty"`
	PodAnnotations     map[string]string  `json:"podAnnotations,omitempty"`
	PodLabels          map[string]string  `json:"podLabels,omitempty"`
	Resources          ContainerResources `json:"resources,omitempty"`
	PriorityClass      string             `json:"priorityClassName,omitempty"`
	SecurityContext    SecurityContext    `json:"securityContext,omitempty"`
	ServiceAnnotations map[string]string  `json:"serviceAnnotations,omitempty"`
}

// K3SAPIServer - parameters to create the vcluster container with
type K3SAPIServer struct {
	Image string         `json:"image,omitempty"`
	Env   []ContainerEnv `json:"env,omitempty"`
}

type ContainerResources struct {
	Limits   ContainerResourcesLimits `json:"limits,omitempty"`
	Requests types.ContainerMemoryCPU `json:"requests,omitempty"`
}

type ContainerResourcesLimits struct {
	Memory string `json:"memory,omitempty"`
}

type ContainerVolumeMounts struct {
	Name      string `json:"name,omitempty"`
	MountPath string `json:"mountPath,omitempty"`
}

// Init - resources which are created during the init phase of the vcluster
type Init struct {
	Manifests string               `json:"manifests,omitempty"`
	Helm      []v1alpha1.HelmChart `json:"helm,omitempty"`
}

// Syncer - parameters to create the syncer with
// https://www.vcluster.com/docs/architecture/basics#vcluster-syncer
type Syncer struct {
	KubeConfigContextName string                   `json:"kubeConfigContextName,omitempty"`
	ExtraArgs             []string                 `json:"extraArgs,omitempty"`
	Limits                types.ContainerMemoryCPU `json:"limits,omitempty"`
}

type Plugin struct {
	Env             []ContainerEnv `json:"env,omitempty"`
	Image           string         `json:"image,omitempty"`
	ImagePullPolicy string         `json:"imagePullPolicy,omitempty"`
	Rbac            Rbac           `json:"rbac,omitempty"`
}

type Plugins struct {
	UffizziClusterSyncPlugin Plugin `json:"ucluster-sync-plugin,omitempty"`
}

type ContainerEnv struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Rbac struct {
	Role        RbacRole        `json:"role,omitempty"`
	ClusterRole RbacClusterRole `json:"clusterRole,omitempty"`
}

type RbacRole struct {
	ExtraRules []RbacRule `json:"extraRules,omitempty"`
}

type RbacClusterRole struct {
	ExtraRules []RbacRule `json:"extraRules,omitempty"`
}

type RbacRule struct {
	ApiGroups []string `json:"apiGroups,omitempty"`
	Resources []string `json:"resources,omitempty"`
	Verbs     []string `json:"verbs,omitempty"`
}

type Sync struct {
	Ingresses EnabledBool `json:"ingresses,omitempty"`
}

type EnabledBool struct {
	Enabled bool `json:"enabled"`
}

type Ingress struct {
	Enabled          bool              `json:"enabled,omitempty"`
	IngressClassName string            `json:"ingressClassName,omitempty"`
	Host             string            `json:"host,omitempty"`
	Annotations      map[string]string `json:"annotations,omitempty"`
}

type ResourceQuota struct {
	Enabled bool                   `json:"enabled,omitempty"`
	Quota   ResourceQuotaDefiniton `json:"quota,omitempty"`
}

type ResourceQuotaDefiniton struct {
	RequestsCpu                 string `json:"requests.cpu,omitempty"`
	RequestsMemory              string `json:"requests.memory,omitempty"`
	RequestsStorage             string `json:"requests.storage,omitempty"`
	RequestsEphemeralStorage    string `json:"requests.ephemeral-storage,omitempty"`
	LimitsCpu                   string `json:"limits.cpu,omitempty"`
	LimitsMemory                string `json:"limits.memory,omitempty"`
	LimitsEphemeralStorage      string `json:"limits.ephemeral-storage,omitempty"`
	ServicesLoadbalancers       int    `json:"services.loadbalancers,omitempty"`
	ServicesNodePorts           int    `json:"services.nodeports,omitempty"`
	CountEndpoints              int    `json:"count/endpoints,omitempty"`
	CountPods                   int    `json:"count/pods,omitempty"`
	CountServices               int    `json:"count/services,omitempty"`
	CountSecrets                int    `json:"count/secrets,omitempty"`
	CountConfigmaps             int    `json:"count/configmaps,omitempty"`
	CountPersistentVolumeClaims int    `json:"count/persistentvolumeclaims,omitempty"`
}

type LimitRangeResources struct {
	EphemeralStorage string `json:"ephemeral-storage,omitempty"`
	Memory           string `json:"memory,omitempty"`
	Cpu              string `json:"cpu,omitempty"`
}

type LimitRange struct {
	Enabled        bool                `json:"enabled,omitempty"`
	Default        LimitRangeResources `json:"default,omitempty"`
	DefaultRequest LimitRangeResources `json:"defaultRequest,omitempty"`
}

type NetworkPolicy struct {
	Enabled bool `json:"enabled,omitempty"`
}

type MapServicesFromVirtual struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type MapServices struct {
	FromVirtual []MapServicesFromVirtual `json:"fromVirtual"`
}

// Isolation - parameters to define the isolation of the cluster
type Isolation struct {
	Enabled             bool          `json:"enabled,omitempty"`
	PodSecurityStandard string        `json:"podSecurityStandard,omitempty"`
	ResourceQuota       ResourceQuota `json:"resourceQuota,omitempty"`
	LimitRange          LimitRange    `json:"limitRange,omitempty"`
	NetworkPolicy       NetworkPolicy `json:"networkPolicy,omitempty"`
}

// NodeSelector - parameters to define the node selector of the cluster
type NodeSelector struct {
	SandboxGKEIORuntime string `json:"sandbox.gke.io/runtime"`
}

type SecurityContextCapabilities struct {
	Drop []string `json:"drop"`
}

// SecurityContext - parameters to define the security context of the cluster
type SecurityContext struct {
	Capabilities           SecurityContextCapabilities `json:"capabilities"`
	ReadOnlyRootFilesystem bool                        `json:"readOnlyRootFilesystem"`
	RunAsNonRoot           bool                        `json:"runAsNonRoot"`
	RunAsUser              int64                       `json:"runAsUser"`
}

type Toleration struct {
	Effect   string `json:"effect"`
	Key      string `json:"key"`
	Operator string `json:"operator"`
}

type Storage struct {
	Persistence bool `json:"persistence"`
}
