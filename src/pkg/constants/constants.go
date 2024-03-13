package constants

const (
	// versions
	ETCD_CHART_VERSION         = "9.5.6"
	VCLUSTER_CHART_K3S_VERSION = "0.19.4"
	VCLUSTER_CHART_K8S_VERSION = "0.19.4"

	// images
	DEFAULT_K3S_VERSION_TAG  = "rancher/k3s:v1.27.3-k3s1"
	UCLUSTER_SYNC_PLUGIN_TAG = "uffizzi/ucluster-sync-plugin:v0.2.4"

	// keys
	UFFIZZI_APP_COMPONENT_LABEL_KEY = "app.uffizzi.com/component"
	VCLUSTER_MANAGED_BY_KEY         = "vcluster.loft.sh/managed-by"

	// env keys
	K3S_DATASTORE_ENDPOINT_ENVKEY = "K3S_DATASTORE_ENDPOINT"
	VCLUSTER_INGRESS_HOST_ENVKEY  = "VCLUSTER_INGRESS_HOST"

	// chart repo urls
	BITNAMI_CHART_REPO_URL = "oci://registry-1.docker.io/bitnamicharts"
	LOFT_CHART_REPO_URL    = "https://charts.loft.sh"

	// other
	VCLUSTER                     = "vcluster"
	ETCD                         = "etcd"
	UCLUSTER_NAME_PREFIX         = "uc-"
	LOFT_HELM_REPO               = "loft"
	BITNAMI_HELM_REPO            = "bitnami"
	VCLUSTER_CHART_K3S           = "vcluster"
	VCLUSTER_CHART_K8S           = "vcluster-k8s"
	ETCD_CHART                   = "etcd"
	VCLUSTER_K3S_DISTRO          = "k3s"
	VCLUSTER_K8S_DISTRO          = "k8s"
	NODESELECTOR_TEMPLATE_GVISOR = "gvisor"
	OCI_TYPE                     = "oci"
	PREMIUM_RWO_STORAGE_CLASS    = "premium-rwo"
	STANDARD_STORAGE_CLASS       = "standard"
)

type LIFECYCLE_OP_TYPE string

const (
	LIFECYCLE_OP_TYPE_CREATE LIFECYCLE_OP_TYPE = "create"
	LIFECYCLE_OP_TYPE_UPDATE LIFECYCLE_OP_TYPE = "update"
	LIFECYCLE_OP_TYPE_DELETE LIFECYCLE_OP_TYPE = "delete"
)
