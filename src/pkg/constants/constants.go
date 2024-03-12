package constants

const (
	UFFIZZI_APP_COMPONENT_LABEL = "app.uffizzi.com/component"
	VCLUSTER                    = "vcluster"
	ETCD                        = "etcd"
	UCLUSTER_NAME_PREFIX        = "uc-"
	LOFT_HELM_REPO              = "loft"
	BITNAMI_HELM_REPO           = "bitnami"
	VCLUSTER_CHART_K3S          = "vcluster"
	VCLUSTER_CHART_K8S          = "vcluster-k8s"
	VCLUSTER_CHART_K3S_VERSION  = "0.19.4"
	VCLUSTER_CHART_K8S_VERSION  = "0.19.4"
	ETCD_CHART                  = "etcd"
	ETCD_CHART_VERSION          = "9.5.6"
	BITNAMI_CHART_REPO_URL      = "oci://registry-1.docker.io/bitnamicharts"
	LOFT_CHART_REPO_URL         = "https://charts.loft.sh"
	VCLUSTER_K3S_DISTRO         = "k3s"
	VCLUSTER_K8S_DISTRO         = "k8s"
	NODESELECTOR_GKE            = "gvisor"
	K3S_DATASTORE_ENDPOINT      = "K3S_DATASTORE_ENDPOINT"
	VCLUSTER_INGRESS_HOSTNAME   = "VCLUSTER_INGRESS_HOST"
	DEFAULT_K3S_VERSION         = "rancher/k3s:v1.27.3-k3s1"
	UCLUSTER_SYNC_PLUGIN_TAG    = "uffizzi/ucluster-sync-plugin:v0.2.4"
	OCI_TYPE                    = "oci"
	PREMIUM_RWO_STORAGE_CLASS   = "premium-rwo"
	STANDARD_STORAGE_CLASS      = "storage"
	SANDBOX_GKE_IO_RUNTIME      = "sandbox.gke.io/runtime"
	GVISOR                      = "gvisor"
	VCLUSTER_MANAGED_BY_KEY     = "vcluster.loft.sh/managed-by"
	WORKLOAD_TYPE_DEPLOYMENT    = "deployment"
	WORKLOAD_TYPE_STATEFULSET   = "statefulset"
)

type LIFECYCLE_OP_TYPE string

const (
	LIFECYCLE_OP_TYPE_CREATE LIFECYCLE_OP_TYPE = "create"
	LIFECYCLE_OP_TYPE_UPDATE LIFECYCLE_OP_TYPE = "update"
	LIFECYCLE_OP_TYPE_DELETE LIFECYCLE_OP_TYPE = "delete"
)
