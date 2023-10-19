package uffizzicluster

import "github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"

type VClusterK3S struct {
	VCluster        VClusterContainer       `json:"vcluster,omitempty"`
	Init            VClusterInit            `json:"init,omitempty"`
	Syncer          VClusterSyncer          `json:"syncer,omitempty"`
	Sync            VClusterSync            `json:"sync,omitempty"`
	Ingress         VClusterIngress         `json:"ingress,omitempty"`
	FsGroup         int64                   `json:"fsgroup,omitempty"`
	Isolation       VClusterIsolation       `json:"isolation,omitempty"`
	NodeSelector    VClusterNodeSelector    `json:"nodeSelector,omitempty"`
	SecurityContext VClusterSecurityContext `json:"securityContext,omitempty"`
	Tolerations     []VClusterToleration    `json:"tolerations,omitempty"`
	MapServices     VClusterMapServices     `json:"mapServices,omitempty"`
	Plugin          VClusterPlugins         `json:"plugin,omitempty"`
}

type VClusterK8S struct {
	APIServer       VClusterK8SAPIServer    `json:"apiserver,omitempty"`
	Init            VClusterInit            `json:"init,omitempty"`
	Syncer          VClusterSyncer          `json:"syncer,omitempty"`
	Sync            VClusterSync            `json:"sync,omitempty"`
	Ingress         VClusterIngress         `json:"ingress,omitempty"`
	FsGroup         int64                   `json:"fsgroup,omitempty"`
	Isolation       VClusterIsolation       `json:"isolation,omitempty"`
	NodeSelector    VClusterNodeSelector    `json:"nodeSelector,omitempty"`
	SecurityContext VClusterSecurityContext `json:"securityContext,omitempty"`
	Tolerations     []VClusterToleration    `json:"tolerations,omitempty"`
	MapServices     VClusterMapServices     `json:"mapServices,omitempty"`
	Plugin          VClusterPlugins         `json:"plugin,omitempty"`
}

type VClusterK8SAPIServer struct {
	Image              string                     `json:"image,omitempty"`
	ExtraArgs          []string                   `json:"extraArgs,omitempty"`
	Replicas           int32                      `json:"replicas,omitempty"`
	NodeSelector       VClusterNodeSelector       `json:"nodeSelector,omitempty"`
	Tolerations        []VClusterToleration       `json:"tolerations,omitempty"`
	Labels             map[string]string          `json:"labels,omitempty"`
	Annotations        map[string]string          `json:"annotations,omitempty"`
	PodAnnotations     map[string]string          `json:"podAnnotations,omitempty"`
	PodLabels          map[string]string          `json:"podLabels,omitempty"`
	Resources          VClusterContainerResources `json:"resources,omitempty"`
	PriorityClass      string                     `json:"priorityClassName,omitempty"`
	SecurityContext    VClusterSecurityContext    `json:"securityContext,omitempty"`
	ServiceAnnotations map[string]string          `json:"serviceAnnotations,omitempty"`
}

// VClusterContainer - parameters to create the vcluster container with
type VClusterContainer struct {
	Image string `json:"image,omitempty"`
}

type VClusterContainerResources struct {
	Limits   VClusterContainerResourcesLimits `json:"limits,omitempty"`
	Requests ContainerMemoryCPU               `json:"requests,omitempty"`
}

type VClusterContainerResourcesLimits struct {
	Memory string `json:"memory,omitempty"`
}

type ContainerMemoryCPU struct {
	Memory string `json:"memory,omitempty"`
	CPU    string `json:"cpu,omitempty"`
}

type VClusterContainerVolumeMounts struct {
	Name      string `json:"name,omitempty"`
	MountPath string `json:"mountPath,omitempty"`
}

// VClusterInit - resources which are created during the init phase of the vcluster
type VClusterInit struct {
	Manifests string               `json:"manifests,omitempty"`
	Helm      []v1alpha1.HelmChart `json:"helm,omitempty"`
}

// VClusterSyncer - parameters to create the syncer with
// https://www.vcluster.com/docs/architecture/basics#vcluster-syncer
type VClusterSyncer struct {
	KubeConfigContextName string             `json:"kubeConfigContextName,omitempty"`
	ExtraArgs             []string           `json:"extraArgs,omitempty"`
	Limits                ContainerMemoryCPU `json:"limits,omitempty"`
}

type VClusterPlugins struct {
	UffizziClusterSyncPlugin VClusterPlugin `json:"ucluster-sync-plugin,omitempty"`
}

type VClusterPlugin struct {
	Env             []VClusterContainerEnv `json:"env,omitempty"`
	Image           string                 `json:"image,omitempty"`
	ImagePullPolicy string                 `json:"imagePullPolicy,omitempty"`
	Rbac            VClusterRbac           `json:"rbac,omitempty"`
}

type VClusterContainerEnv struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type VClusterRbac struct {
	Role        VClusterRbacRole        `json:"role,omitempty"`
	ClusterRole VClusterRbacClusterRole `json:"clusterRole,omitempty"`
}

type VClusterRbacRole struct {
	ExtraRules []VClusterRbacRule `json:"extraRules,omitempty"`
}

type VClusterRbacClusterRole struct {
	ExtraRules []VClusterRbacRule `json:"extraRules,omitempty"`
}

type VClusterRbacRule struct {
	ApiGroups []string `json:"apiGroups,omitempty"`
	Resources []string `json:"resources,omitempty"`
	Verbs     []string `json:"verbs,omitempty"`
}

type VClusterSync struct {
	Ingresses EnabledBool `json:"ingresses,omitempty"`
}

type EnabledBool struct {
	Enabled bool `json:"enabled"`
}

type VClusterIngress struct {
	Enabled          bool              `json:"enabled,omitempty"`
	IngressClassName string            `json:"ingressClassName,omitempty"`
	Host             string            `json:"host,omitempty"`
	Annotations      map[string]string `json:"annotations,omitempty"`
}

type VClusterResourceQuota struct {
	Enabled bool                           `json:"enabled,omitempty"`
	Quota   VClusterResourceQuotaDefiniton `json:"quota,omitempty"`
}

type VClusterResourceQuotaDefiniton struct {
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

type VClusterLimitRange struct {
	Enabled        bool                `json:"enabled,omitempty"`
	Default        LimitRangeResources `json:"default,omitempty"`
	DefaultRequest LimitRangeResources `json:"defaultRequest,omitempty"`
}

type VClusterNetworkPolicy struct {
	Enabled bool `json:"enabled,omitempty"`
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
	Enabled             bool                  `json:"enabled,omitempty"`
	PodSecurityStandard string                `json:"podSecurityStandard,omitempty"`
	ResourceQuota       VClusterResourceQuota `json:"resourceQuota,omitempty"`
	LimitRange          VClusterLimitRange    `json:"limitRange,omitempty"`
	NetworkPolicy       VClusterNetworkPolicy `json:"networkPolicy,omitempty"`
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
