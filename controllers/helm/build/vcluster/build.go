package vcluster

import (
	"fmt"
	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/controllers/constants"
	"github.com/UffizziCloud/uffizzi-cluster-operator/controllers/etcd"
	"github.com/UffizziCloud/uffizzi-cluster-operator/controllers/helm/types"
	"github.com/UffizziCloud/uffizzi-cluster-operator/controllers/helm/types/vcluster"
	v1 "k8s.io/api/core/v1"
)

func BuildK3SHelmValues(uCluster *v1alpha1.UffizziCluster) (vcluster.K3S, string) {
	helmReleaseName, vclusterIngressHostname, outKubeConfigServerArgValue := configStrings(uCluster)

	vclusterK3sHelmValues := vcluster.K3S{
		VCluster: k3SAPIServer(uCluster),
		Common:   common(helmReleaseName, vclusterIngressHostname),
	}

	if uCluster.Spec.ExternalDatastore == constants.ETCD || uCluster.Spec.ExternalDatastore == "" {
		vclusterK3sHelmValues.VCluster.Env = []vcluster.ContainerEnv{
			{
				Name:  constants.K3S_DATASTORE_ENDPOINT,
				Value: etcd.BuildEtcdHelmReleaseName(uCluster) + "." + uCluster.Namespace + ".svc.cluster.local:2379",
			},
		}
		vclusterK3sHelmValues.Storage.Persistence = false
		vclusterK3sHelmValues.EnableHA = true
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

	vclusterK3sHelmValues.Syncer.ExtraArgs = append(vclusterK3sHelmValues.Syncer.ExtraArgs,
		"--tls-san="+vclusterIngressHostname,
		"--out-kube-config-server="+outKubeConfigServerArgValue,
	)

	if len(uCluster.Spec.Helm) > 0 {
		vclusterK3sHelmValues.Init.Helm = uCluster.Spec.Helm
	}

	if uCluster.Spec.Manifests != nil {
		vclusterK3sHelmValues.Init.Manifests = *uCluster.Spec.Manifests
	}

	return vclusterK3sHelmValues, helmReleaseName
}

func BuildK8SHelmValues(uCluster *v1alpha1.UffizziCluster) (vcluster.K8S, string) {
	helmReleaseName, vclusterIngressHostname, outKubeConfigServerArgValue := configStrings(uCluster)

	vclusterHelmValues := vcluster.K8S{
		APIServer: k8SAPIServer(),
		Common:    common(helmReleaseName, vclusterIngressHostname),
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

	vclusterHelmValues.Syncer.ExtraArgs = append(vclusterHelmValues.Syncer.ExtraArgs,
		"--tls-san="+vclusterIngressHostname,
		"--out-kube-config-server="+outKubeConfigServerArgValue,
	)

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
	return vcluster.Syncer{
		KubeConfigContextName: helmReleaseName,
		ExtraArgs: []string{
			fmt.Sprintf(
				"--enforce-toleration=%s:%s",
				constants.SANDBOX_GKE_IO_RUNTIME,
				string(v1.TaintEffectNoSchedule),
			),
			"--node-selector=sandbox.gke.io/runtime=gvisor",
			"--enforce-node-selector",
		},
		Limits: types.ContainerMemoryCPU{
			CPU:    "1000m",
			Memory: "1024Mi",
		},
	}
}

func syncConfig() vcluster.Sync {
	return vcluster.Sync{
		Ingresses: vcluster.EnabledBool{
			Enabled: false,
		},
	}
}

func tolerations() []vcluster.Toleration {
	return []vcluster.Toleration{
		{
			Key:      constants.SANDBOX_GKE_IO_RUNTIME,
			Effect:   string(v1.TaintEffectNoSchedule),
			Operator: string(v1.NodeSelectorOpExists),
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
				RequestsEphemeralStorage:    "15Gi",
				RequestsStorage:             "10Gi",
				LimitsCpu:                   "10",
				LimitsMemory:                "15Gi",
				LimitsEphemeralStorage:      "30Gi",
				ServicesLoadbalancers:       5,
				ServicesNodePorts:           0,
				CountEndpoints:              40,
				CountConfigmaps:             100,
				CountPersistentVolumeClaims: 40,
				CountPods:                   40,
				CountSecrets:                100,
				CountServices:               40,
			},
		},
		LimitRange: vcluster.LimitRange{
			Enabled: true,
			Default: vcluster.LimitRangeResources{
				Cpu:              "1",
				Memory:           "512Mi",
				EphemeralStorage: "8Gi",
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

func nodeSelector() vcluster.NodeSelector {
	return vcluster.NodeSelector{
		SandboxGKEIORuntime: "gvisor",
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

func common(helmReleaseName, vclusterIngressHostname string) vcluster.Common {
	return vcluster.Common{
		Init:            vcluster.Init{},
		FsGroup:         12345,
		Ingress:         ingress(vclusterIngressHostname),
		Isolation:       isolation(),
		NodeSelector:    nodeSelector(),
		SecurityContext: securityContext(),
		Tolerations:     tolerations(),
		Plugin:          pluginsConfig(),
		Syncer:          syncerConfig(helmReleaseName),
		Sync:            syncConfig(),
	}
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
