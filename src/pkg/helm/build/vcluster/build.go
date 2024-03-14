package vcluster

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/controllers/etcd"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/constants"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/helm/types"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/helm/types/vcluster"
	v1 "k8s.io/api/core/v1"
)

func BuildK3SHelmValues(uCluster *v1alpha1.UffizziCluster) (vcluster.K3S, string) {
	var (
		helmReleaseName, vclusterIngressHostname, outKubeConfigServerArgValue = configStrings(uCluster)
		nodeSelector                                                          = vcluster.NodeSelector{}
		tolerations                                                           = []v1.Toleration{}
	)

	if uCluster.Spec.NodeSelectorTemplate == constants.GVISOR {
		nodeSelector = vcluster.GvisorNodeSelector
		tolerations = []v1.Toleration{vcluster.GvisorToleration.ToV1()}
	}

	if len(uCluster.Spec.NodeSelector) > 0 {
		// merge nodeSelector and uCluster.Spec.NodeSelector
		for k, v := range uCluster.Spec.NodeSelector {
			map[string]string(nodeSelector)[k] = v
		}
	}

	if len(uCluster.Spec.Toleration) > 0 {
		// merge tolerations and uCluster.Spec.Toleration
		tolerations = append(tolerations, uCluster.Spec.Toleration...)
	}

	vclusterK3sHelmValues := vcluster.K3S{
		VCluster: k3SAPIServer(uCluster),
		Common:   common(helmReleaseName, vclusterIngressHostname, nodeSelector, tolerations),
	}

	if uCluster.Spec.ExternalDatastore == constants.ETCD {
		vclusterK3sHelmValues.VCluster.Env = []vcluster.ContainerEnv{
			{
				Name:  constants.K3S_DATASTORE_ENDPOINT,
				Value: "http://" + etcd.BuildEtcdHelmReleaseName(uCluster) + "." + uCluster.Namespace + ".svc.cluster.local:2379",
			},
		}
	}

	if uCluster.Spec.Ingress.Host != "" {
		vclusterK3sHelmValues.Plugin.UffizziClusterSyncPlugin.Env = []vcluster.ContainerEnv{
			{
				Name:  constants.VCLUSTER_INGRESS_HOSTNAME,
				Value: vclusterIngressHostname,
			},
		}
	}

	if uCluster.Spec.ResourceQuota != nil {
		// map uCluster.Spec.ResourceQuota to vclusterK3sHelmValues.Isolation.ResourceQuota
		q := *uCluster.Spec.ResourceQuota
		qHelmValues := vclusterK3sHelmValues.Isolation.ResourceQuota
		// enabled
		qHelmValues.Enabled = q.Enabled
		//requests
		qHelmValues.Quota.RequestsMemory = q.Requests.Memory
		qHelmValues.Quota.RequestsCpu = q.Requests.CPU
		qHelmValues.Quota.RequestsEphemeralStorage = q.Requests.EphemeralStorage
		qHelmValues.Quota.RequestsStorage = q.Requests.Storage
		// limits
		qHelmValues.Quota.LimitsMemory = q.Limits.Memory
		qHelmValues.Quota.LimitsCpu = q.Limits.CPU
		qHelmValues.Quota.LimitsEphemeralStorage = q.Limits.EphemeralStorage
		// services
		qHelmValues.Quota.ServicesNodePorts = q.Services.NodePorts
		qHelmValues.Quota.ServicesLoadbalancers = q.Services.LoadBalancers
		// count
		qHelmValues.Quota.CountPods = q.Count.Pods
		qHelmValues.Quota.CountServices = q.Count.Services
		qHelmValues.Quota.CountPersistentVolumeClaims = q.Count.PersistentVolumeClaims
		qHelmValues.Quota.CountConfigmaps = q.Count.ConfigMaps
		qHelmValues.Quota.CountSecrets = q.Count.Secrets
		qHelmValues.Quota.CountEndpoints = q.Count.Endpoints
		// set it back
		vclusterK3sHelmValues.Isolation.ResourceQuota = qHelmValues
	}

	if uCluster.Spec.LimitRange != nil {
		// same for limit range
		lr := uCluster.Spec.LimitRange
		lrHelmValues := vclusterK3sHelmValues.Isolation.LimitRange
		// enabled
		lrHelmValues.Enabled = lr.Enabled
		// default
		lrHelmValues.Default.Cpu = lr.Default.CPU
		lrHelmValues.Default.Memory = lr.Default.Memory
		lrHelmValues.Default.EphemeralStorage = lr.Default.EphemeralStorage
		// default requests
		lrHelmValues.DefaultRequest.Cpu = lr.DefaultRequest.CPU
		lrHelmValues.DefaultRequest.Memory = lr.DefaultRequest.Memory
		lrHelmValues.DefaultRequest.EphemeralStorage = lr.DefaultRequest.EphemeralStorage
		// set it back
		vclusterK3sHelmValues.Isolation.LimitRange = lrHelmValues
	}

	if vclusterIngressHostname != "" {
		vclusterK3sHelmValues.Syncer.ExtraArgs = append(vclusterK3sHelmValues.Syncer.ExtraArgs,
			"--tls-san="+vclusterIngressHostname,
		)
	}

	if outKubeConfigServerArgValue != "" {
		vclusterK3sHelmValues.Syncer.ExtraArgs = append(vclusterK3sHelmValues.Syncer.ExtraArgs,
			"--out-kube-config-server="+outKubeConfigServerArgValue,
		)
	}

	for _, t := range tolerations {
		vclusterK3sHelmValues.Syncer.ExtraArgs = append(vclusterK3sHelmValues.Syncer.ExtraArgs, "--enforce-toleration="+vcluster.Toleration(t).Notation())
		uCluster.Status.AddToleration(t)
	}

	if len(nodeSelector) > 0 {
		for k, v := range nodeSelector {
			vclusterK3sHelmValues.Syncer.ExtraArgs = append(vclusterK3sHelmValues.Syncer.ExtraArgs, "--node-selector="+k+"="+v)
			uCluster.Status.AddNodeSelector(k, v)
		}
		vclusterK3sHelmValues.Syncer.ExtraArgs = append(vclusterK3sHelmValues.Syncer.ExtraArgs, "--enforce-node-selector")
	}

	// keep cluster data intact in case the vcluster scales up or down
	vclusterK3sHelmValues.Syncer.Storage = vcluster.Storage{
		Persistence: true,
		Size:        "5Gi",
	}

	if uCluster.Spec.Storage != nil {
		storage := uCluster.Spec.Storage
		vclusterK3sHelmValues.Syncer.Storage.Persistence = storage.Persistence
		if len(uCluster.Spec.Storage.Size) > 0 {
			vclusterK3sHelmValues.Syncer.Storage.Size = uCluster.Spec.Storage.Size
		}
	}

	if len(uCluster.Spec.Helm) > 0 {
		vclusterK3sHelmValues.Init.Helm = uCluster.Spec.Helm
	}

	if uCluster.Spec.Manifests != nil {
		vclusterK3sHelmValues.Init.Manifests = *uCluster.Spec.Manifests
	}

	return vclusterK3sHelmValues, helmReleaseName
}

func BuildK8SHelmValues(uCluster *v1alpha1.UffizziCluster) (vcluster.K8S, string) {
	var (
		helmReleaseName, vclusterIngressHostname, outKubeConfigServerArgValue = configStrings(uCluster)
		nodeSelector                                                          = vcluster.NodeSelector{}
		tolerations                                                           = []v1.Toleration{}
	)
	if uCluster.Spec.NodeSelectorTemplate == constants.GVISOR {
		nodeSelector = vcluster.GvisorNodeSelector
		tolerations = []v1.Toleration{vcluster.GvisorToleration.ToV1()}
	}

	if len(uCluster.Spec.NodeSelector) > 0 {
		// merge nodeSelector and uCluster.Spec.NodeSelector
		for k, v := range uCluster.Spec.NodeSelector {
			map[string]string(nodeSelector)[k] = v
		}
	}

	if len(uCluster.Spec.Toleration) > 0 {
		// merge tolerations and uCluster.Spec.Toleration
		tolerations = append(tolerations, uCluster.Spec.Toleration...)
	}

	vclusterHelmValues := vcluster.K8S{
		APIServer: k8SAPIServer(),
		Common:    common(helmReleaseName, vclusterIngressHostname, nodeSelector, tolerations),
	}

	if uCluster.Spec.APIServer.Image != "" {
		vclusterHelmValues.APIServer.Image = uCluster.Spec.APIServer.Image
	}

	if uCluster.Spec.Ingress.Host != "" {
		vclusterHelmValues.Plugin.UffizziClusterSyncPlugin.Env = []vcluster.ContainerEnv{
			{
				Name:  "VCLUSTER_INGRESS_HOST",
				Value: vclusterIngressHostname,
			},
		}
	}

	if uCluster.Spec.ResourceQuota != nil {
		// map uCluster.Spec.ResourceQuota to vclusterHelmValues.Isolation.ResourceQuota
		q := *uCluster.Spec.ResourceQuota
		qHelmValues := vclusterHelmValues.Isolation.ResourceQuota
		// enabled
		qHelmValues.Enabled = q.Enabled
		//requests
		qHelmValues.Quota.RequestsMemory = q.Requests.Memory
		qHelmValues.Quota.RequestsCpu = q.Requests.CPU
		qHelmValues.Quota.RequestsEphemeralStorage = q.Requests.EphemeralStorage
		qHelmValues.Quota.RequestsStorage = q.Requests.Storage
		// limits
		qHelmValues.Quota.LimitsMemory = q.Limits.Memory
		qHelmValues.Quota.LimitsCpu = q.Limits.CPU
		qHelmValues.Quota.LimitsEphemeralStorage = q.Limits.EphemeralStorage
		// services
		qHelmValues.Quota.ServicesNodePorts = q.Services.NodePorts
		qHelmValues.Quota.ServicesLoadbalancers = q.Services.LoadBalancers
		// count
		qHelmValues.Quota.CountPods = q.Count.Pods
		qHelmValues.Quota.CountServices = q.Count.Services
		qHelmValues.Quota.CountPersistentVolumeClaims = q.Count.PersistentVolumeClaims
		qHelmValues.Quota.CountConfigmaps = q.Count.ConfigMaps
		qHelmValues.Quota.CountSecrets = q.Count.Secrets
		qHelmValues.Quota.CountEndpoints = q.Count.Endpoints
		// set it back
		vclusterHelmValues.Isolation.ResourceQuota = qHelmValues
	}

	if uCluster.Spec.LimitRange != nil {
		// same for limit range
		lr := uCluster.Spec.LimitRange
		lrHelmValues := vclusterHelmValues.Isolation.LimitRange
		// enabled
		lrHelmValues.Enabled = lr.Enabled
		// default
		lrHelmValues.Default.Cpu = lr.Default.CPU
		lrHelmValues.Default.Memory = lr.Default.Memory
		lrHelmValues.Default.EphemeralStorage = lr.Default.EphemeralStorage
		// default requests
		lrHelmValues.DefaultRequest.Cpu = lr.DefaultRequest.CPU
		lrHelmValues.DefaultRequest.Memory = lr.DefaultRequest.Memory
		lrHelmValues.DefaultRequest.EphemeralStorage = lr.DefaultRequest.EphemeralStorage
		// set it back
		vclusterHelmValues.Isolation.LimitRange = lrHelmValues
	}

	if vclusterIngressHostname != "" {
		vclusterHelmValues.Syncer.ExtraArgs = append(vclusterHelmValues.Syncer.ExtraArgs,
			"--tls-san="+vclusterIngressHostname,
		)
	}

	if outKubeConfigServerArgValue != "" {
		vclusterHelmValues.Syncer.ExtraArgs = append(vclusterHelmValues.Syncer.ExtraArgs,
			"--out-kube-config-server="+outKubeConfigServerArgValue,
		)
	}

	for _, t := range tolerations {
		vclusterHelmValues.Syncer.ExtraArgs = append(vclusterHelmValues.Syncer.ExtraArgs, "--enforce-toleration="+vcluster.Toleration(t).Notation())
		uCluster.Status.AddToleration(t)
	}

	if len(nodeSelector) > 0 {
		for k, v := range nodeSelector {
			vclusterHelmValues.Syncer.ExtraArgs = append(vclusterHelmValues.Syncer.ExtraArgs, "--node-selector="+k+"="+v)
			uCluster.Status.AddNodeSelector(k, v)
		}
		vclusterHelmValues.Syncer.ExtraArgs = append(vclusterHelmValues.Syncer.ExtraArgs, "--enforce-node-selector")
	}

	if len(uCluster.Spec.Helm) > 0 {
		vclusterHelmValues.Init.Helm = uCluster.Spec.Helm
	}

	if uCluster.Spec.Manifests != nil {
		vclusterHelmValues.Init.Manifests = *uCluster.Spec.Manifests
	}
	return vclusterHelmValues, helmReleaseName
}

func pluginsConfig() vcluster.Plugins {
	return vcluster.Plugins{
		UffizziClusterSyncPlugin: vcluster.Plugin{
			Image: constants.UCLUSTER_SYNC_PLUGIN_TAG,
			Rbac: vcluster.Rbac{
				Role: vcluster.RbacRole{
					ExtraRules: []vcluster.RbacRule{
						{
							ApiGroups: []string{"networking.k8s.io"},
							Resources: []string{"ingresses"},
							Verbs:     []string{"create", "delete", "patch", "update", "get", "list", "watch"},
						},
					},
				},
				ClusterRole: vcluster.RbacClusterRole{
					ExtraRules: []vcluster.RbacRule{
						{
							ApiGroups: []string{"apiextensions.k8s.io"},
							Resources: []string{"customresourcedefinitions"},
							Verbs:     []string{"patch", "update", "get", "list", "watch"},
						},
					},
				},
			},
		},
	}
}

func syncerConfig(helmReleaseName string) vcluster.Syncer {
	syncer := vcluster.Syncer{
		KubeConfigContextName: helmReleaseName,
		Limits: types.ContainerMemoryCPU{
			CPU:    "1000m",
			Memory: "1024Mi",
		},
	}

	return syncer
}

func syncConfig() vcluster.Sync {
	return vcluster.Sync{
		Ingresses: vcluster.EnabledBool{
			Enabled: false,
		},
	}
}

func securityContext() vcluster.SecurityContext {
	return vcluster.SecurityContext{
		Capabilities: vcluster.SecurityContextCapabilities{
			Drop: []string{"all"},
		},
	}
}

func isolation() vcluster.Isolation {
	return vcluster.Isolation{
		Enabled:             true,
		PodSecurityStandard: "baseline",
		ResourceQuota: vcluster.ResourceQuota{
			Enabled: true,
			Quota: vcluster.ResourceQuotaDefiniton{
				RequestsCpu:                 "2.5",
				RequestsMemory:              "10Gi",
				RequestsEphemeralStorage:    "50Gi",
				RequestsStorage:             "20Gi",
				LimitsCpu:                   "20",
				LimitsMemory:                "30Gi",
				LimitsEphemeralStorage:      "80Gi",
				ServicesLoadbalancers:       100,
				ServicesNodePorts:           0,
				CountEndpoints:              100,
				CountConfigmaps:             100,
				CountPersistentVolumeClaims: 100,
				CountPods:                   100,
				CountSecrets:                100,
				CountServices:               100,
			},
		},
		LimitRange: vcluster.LimitRange{
			Enabled: true,
			Default: vcluster.LimitRangeResources{
				Cpu:              "1",
				Memory:           "1Gi",
				EphemeralStorage: "16Gi",
			},
			DefaultRequest: vcluster.LimitRangeResources{
				Cpu:              "100m",
				Memory:           "128Mi",
				EphemeralStorage: "3Gi",
			},
		},
		NetworkPolicy: vcluster.NetworkPolicy{
			Enabled: true,
		},
	}
}

func ingress(VClusterIngressHostname string) vcluster.Ingress {
	return vcluster.Ingress{
		Enabled: true,
		Host:    VClusterIngressHostname,
		Annotations: map[string]string{
			"app.uffizzi.com/ingress-sync": "true",
		},
	}
}

func configStrings(uCluster *v1alpha1.UffizziCluster) (string, string, string) {
	helmReleaseName := BuildVClusterHelmReleaseName(uCluster)
	var (
		VClusterIngressHostname     = BuildVClusterIngressHost(uCluster)
		OutKubeConfigServerArgValue = ""
	)

	if VClusterIngressHostname != "" {
		OutKubeConfigServerArgValue = "https://" + VClusterIngressHostname
	}
	return helmReleaseName, VClusterIngressHostname, OutKubeConfigServerArgValue
}

func k8SAPIServer() vcluster.K8SAPIServer {
	return vcluster.K8SAPIServer{
		Image: "registry.k8s.io/kube-apiserver:v1.26.1",
		Resources: vcluster.ContainerResources{
			Requests: types.ContainerMemoryCPU{
				CPU:    "40m",
				Memory: "300Mi",
			},
		},
	}
}

func common(helmReleaseName, vclusterIngressHostname string, nodeSelector map[string]string, tolerations []v1.Toleration) vcluster.Common {
	c := vcluster.Common{
		Init:            vcluster.Init{},
		FsGroup:         12345,
		Ingress:         ingress(vclusterIngressHostname),
		Isolation:       isolation(),
		NodeSelector:    nodeSelector,
		Tolerations:     tolerations,
		SecurityContext: securityContext(),
		Plugin:          pluginsConfig(),
		Syncer:          syncerConfig(helmReleaseName),
		Sync:            syncConfig(),
	}

	return c
}

func k3SAPIServer(uCluster *v1alpha1.UffizziCluster) vcluster.K3SAPIServer {
	apiserver := vcluster.K3SAPIServer{
		Image: constants.DEFAULT_K3S_VERSION,
	}
	if uCluster.Spec.APIServer.Image != "" {
		apiserver.Image = uCluster.Spec.APIServer.Image
	}
	return apiserver
}

func BuildVClusterIngressHost(uCluster *v1alpha1.UffizziCluster) string {
	host := ""
	if uCluster.Spec.Ingress.Host != "" {
		host = uCluster.Name + "-" + uCluster.Spec.Ingress.Host
	}
	return host
}

func BuildVClusterHelmReleaseName(uCluster *v1alpha1.UffizziCluster) string {
	helmReleaseName := constants.UCLUSTER_NAME_PREFIX + uCluster.Name
	return helmReleaseName
}
