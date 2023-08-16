/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"
	"github.com/fluxcd/pkg/apis/meta"
	"github.com/pkg/errors"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	controllerruntimesource "sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
	"time"

	uclusteruffizzicomv1alpha1 "github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	fluxhelmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	fluxsourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// UffizziClusterReconciler reconciles a UffizziCluster object
type UffizziClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	UCLUSTER_NAME_PREFIX   = "uc-"
	LOFT_HELM_REPO         = "loft"
	VCLUSTER_CHART         = "vcluster"
	VCLUSTER_CHART_VERSION = "0.15.5"
	LOFT_CHART_REPO_URL    = "https://charts.loft.sh"
)

type LIFECYCLE_OP_TYPE string

var (
	LIFECYCLE_OP_TYPE_CREATE LIFECYCLE_OP_TYPE = "create"
	LIFECYCLE_OP_TYPE_UPDATE LIFECYCLE_OP_TYPE = "update"
	LIFECYCLE_OP_TYPE_DELETE LIFECYCLE_OP_TYPE = "delete"
)

//+kubebuilder:rbac:groups=uffizzi.com,resources=uffizziclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=uffizzi.com,resources=uffizziclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=uffizzi.com,resources=uffizziclusters/finalizers,verbs=update

// add the helm controller rbac
//+kubebuilder:rbac:groups=helm.toolkit.fluxcd.io,resources=helmreleases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=helm.toolkit.fluxcd.io,resources=helmreleases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=helm.toolkit.fluxcd.io,resources=helmreleases/finalizers,verbs=update

// add the source controller rbac
//+kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=helmrepositories,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=helmrepositories/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=helmrepositories/finalizers,verbs=update

// add the ingress rbac
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete

// add services rbac
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// add secret rbac
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the UffizziCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *UffizziClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		lifecycleOpType LIFECYCLE_OP_TYPE
		currentSpec     string
		lastAppliedSpec string
	)
	// default lifecycle operation
	lifecycleOpType = LIFECYCLE_OP_TYPE_CREATE
	logger := log.FromContext(ctx)

	// Fetch the UffizziCluster instance in question
	uCluster := &uclusteruffizzicomv1alpha1.UffizziCluster{}
	if err := r.Get(ctx, req.NamespacedName, uCluster); err != nil {
		// possibly a delete event
		if k8serrors.IsNotFound(err) {
			lifecycleOpType = LIFECYCLE_OP_TYPE_DELETE
			logger.Info("UffizziCluster deleted", "NamespacedName", req.NamespacedName)
		} else {
			logger.Info("Failed to get UffizziCluster", "NamespacedName", req.NamespacedName, "Error", err)
		}
		// Handle error
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if lifecycleOpType == LIFECYCLE_OP_TYPE_CREATE {
		currentSpecBytes, err := json.Marshal(uCluster.Spec)
		if err != nil {
			logger.Error(err, "Failed to marshal current spec")
			return ctrl.Result{}, err
		}
		currentSpec = string(currentSpecBytes)
	}

	// Set default values for the status if it is not already set
	if len(uCluster.Status.Conditions) == 0 {
		var (
			intialConditions = []metav1.Condition{
				buildInitializingCondition(),
			}
			helmReleaseRef  = ""
			exposedServices = []uclusteruffizzicomv1alpha1.ExposedVClusterServiceStatus{}
			host            = ""
			kubeConfig      = uclusteruffizzicomv1alpha1.VClusterKubeConfig{
				SecretRef: &meta.SecretKeyReference{},
			}
		)
		uCluster.Status = uclusteruffizzicomv1alpha1.UffizziClusterStatus{
			Conditions:      intialConditions,
			HelmReleaseRef:  &helmReleaseRef,
			ExposedServices: exposedServices,
			Host:            &host,
			KubeConfig:      kubeConfig,
		}
		if err := r.Status().Update(ctx, uCluster); err != nil {
			logger.Error(err, "Failed to update the default UffizziCluster status")
			return ctrl.Result{}, err
		}
	}

	// Check if there is already exists a VCluster HelmRelease for this UCluster, if not create one
	helmReleaseName := BuildVClusterHelmReleaseName(uCluster)
	helmRelease := &fluxhelmv2beta1.HelmRelease{}
	helmReleaseNamespacedName := client.ObjectKey{
		Namespace: uCluster.Namespace,
		Name:      helmReleaseName,
	}

	err := r.Get(ctx, helmReleaseNamespacedName, helmRelease)
	if err != nil && k8serrors.IsNotFound(err) {
		// helm release does not exist so let's create one
		lifecycleOpType = LIFECYCLE_OP_TYPE_CREATE
		newHelmRelease, err := r.upsertVClusterHelmRelease(false, ctx, uCluster)
		if err != nil {
			logger.Error(err, "Failed to create HelmRelease")
			return ctrl.Result{}, err
		}
		// create the ingress for the vcluster
		vclusterIngressHost := BuildVClusterIngressHost(uCluster) // r.createVClusterIngress(ctx, uCluster)
		if err != nil {
			logger.Error(err, "Failed to create ingress to vcluster internal service")
			return ctrl.Result{}, err
		}

		uCluster.Status.Host = &vclusterIngressHost

		// check if there are any services that need to be exposed from the vcluster
		if uCluster.Spec.Ingress.Services != nil {
			// Create the ingress for each service that is specified in the UffizziCluster services config
			for _, service := range uCluster.Spec.Ingress.Services {
				vclusterInternalServiceIngressStatus, err := r.createVClusterInternalServiceIngress(uCluster, service, ctx)
				if err != nil {
					logger.Error(err, "Failed to create ingress to vcluster internal service")
					return ctrl.Result{}, err
				}

				// add the exposed service to the status
				uCluster.Status.ExposedServices = append(uCluster.Status.ExposedServices,
					*vclusterInternalServiceIngressStatus,
				)
			}
		}

		// reference the HelmRelease in the status
		uCluster.Status.HelmReleaseRef = &helmReleaseName
		uCluster.Status.KubeConfig.SecretRef = &meta.SecretKeyReference{
			Name: "vc-" + helmReleaseName,
		}

		if err := r.Status().Update(ctx, uCluster); err != nil {
			logger.Error(err, "Failed to update UffizziCluster status")
			return ctrl.Result{RequeueAfter: time.Second * 5}, err
		}

		logger.Info("Created HelmRelease", "HelmRelease", newHelmRelease.Name)
	} else if err != nil {
		logger.Error(err, "Failed to create HelmRelease, unknown error")
		return ctrl.Result{}, err
	} else {
		lifecycleOpType = LIFECYCLE_OP_TYPE_UPDATE
		// if helm release already exists then replicate the status conditions onto the uffizzicluster object
		uClusterConditions := []metav1.Condition{}
		for _, c := range helmRelease.Status.Conditions {
			helmMessage := "[HelmRelease] " + c.Message
			uClusterCondition := c
			uClusterCondition.Message = helmMessage
			uClusterConditions = append(uClusterConditions, uClusterCondition)
		}

		uCluster.Status.Conditions = uClusterConditions
		if err := r.Status().Update(ctx, uCluster); err != nil {
			logger.Error(err, "Failed to update UffizziCluster status")
			return ctrl.Result{RequeueAfter: time.Second * 5}, err
		}
	}

	if lifecycleOpType == LIFECYCLE_OP_TYPE_CREATE {
		err := r.createLoftHelmRepo(ctx, req)
		if err != nil && k8serrors.IsAlreadyExists(err) {
			logger.Info("Loft Helm Repo for UffizziCluster already exists", "NamespacedName", req.NamespacedName)
		} else {
			logger.Info("Loft Helm Repo for UffizziCluster created", "NamespacedName", req.NamespacedName)
		}
		uCluster.Status.LastAppliedConfiguration = &currentSpec
		if err := r.Status().Update(ctx, uCluster); err != nil {
			logger.Error(err, "Failed to update the default UffizziCluster lastAppliedConfig")
			return ctrl.Result{}, err
		}

		logger.Info("UffizziCluster lastAppliedConfig has been set")
	} else if lifecycleOpType == LIFECYCLE_OP_TYPE_DELETE {
		err := r.deleteLoftHelmRepo(ctx, req)
		if err != nil && k8serrors.IsNotFound(err) {
			logger.Info("Loft Helm Repo for UffizziCluster already deleted", "NamespacedName", req.NamespacedName)
		}
	}

	if lifecycleOpType == LIFECYCLE_OP_TYPE_UPDATE {
		if currentSpec != lastAppliedSpec {
			if _, err := r.upsertVClusterHelmRelease(true, ctx, uCluster); err != nil {
				logger.Error(err, "Failed to update HelmRelease")
				return ctrl.Result{}, err
			}
		}
	}

	// Requeue the request to check the status
	return ctrl.Result{Requeue: true}, nil
}

func (r *UffizziClusterReconciler) createLoftHelmRepo(ctx context.Context, req ctrl.Request) error {
	return r.createHelmRepo(ctx, LOFT_HELM_REPO, req.Namespace, LOFT_CHART_REPO_URL)
}

func (r *UffizziClusterReconciler) deleteLoftHelmRepo(ctx context.Context, req ctrl.Request) error {
	return r.deleteHelmRepo(ctx, LOFT_HELM_REPO, req.Namespace)
}

func (r *UffizziClusterReconciler) createVClusterIngress(ctx context.Context, uCluster *uclusteruffizzicomv1alpha1.UffizziCluster) (*string, error) {
	helmReleaseName := BuildVClusterHelmReleaseName(uCluster)
	clusterIngress := BuildVClusterIngress(helmReleaseName, uCluster)
	ingressHost := BuildVClusterIngressHost(uCluster)
	if err := controllerutil.SetControllerReference(uCluster, clusterIngress, r.Scheme); err != nil {
		return nil, errors.Wrap(err, "Failed to create Cluster Ingress")
	}
	// create the nginx ingress
	//if err := r.Create(ctx, clusterIngress); err != nil {
	//	if strings.Contains(err.Error(), "already exists") {
	//		// If the Ingress already exists, update it
	//		if err := r.Update(ctx, clusterIngress); err != nil {
	//			return nil, errors.Wrap(err, "Failed to update Cluster Ingress")
	//		}
	//	} else {
	//		return nil, errors.Wrap(err, "Failed to create Cluster Ingress")
	//	}
	//}
	return &ingressHost, nil
}

func (r *UffizziClusterReconciler) createVClusterInternalServiceIngress(uCluster *uclusteruffizzicomv1alpha1.UffizziCluster, service uclusteruffizzicomv1alpha1.ExposedVClusterService, ctx context.Context) (*uclusteruffizzicomv1alpha1.ExposedVClusterServiceStatus, error) {
	helmReleaseName := BuildVClusterHelmReleaseName(uCluster)

	vclusterInternalServiceIngress := BuildVClusterInternalServiceIngress(service, uCluster, helmReleaseName)

	// Let the ingresses be owned by UffizziCluster
	if err := controllerutil.SetControllerReference(uCluster, vclusterInternalServiceIngress, r.Scheme); err != nil {
		return nil, errors.Wrap(err, "Failed to set ownerReference for Ingress for internal service "+service.Name+" in namespace "+service.Namespace)
	}

	if err := r.Create(ctx, vclusterInternalServiceIngress); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			// If the Ingress already exists, update it
			if err := r.Update(ctx, vclusterInternalServiceIngress); err != nil {
				return nil, errors.Wrap(err, "Failed to update Ingress for internal service "+service.Name+" in namespace "+service.Namespace)
			}
		} else {
			return nil, errors.Wrap(err, "Failed to create Ingress for internal service "+service.Name+" in namespace "+service.Namespace)
		}
	}

	// Create the ingress status
	vclusterInternalServiceHost := BuildVClusterInternalServiceIngressHost(uCluster)
	status := &uclusteruffizzicomv1alpha1.ExposedVClusterServiceStatus{
		Name:      service.Name,
		Namespace: service.Namespace,
		Host:      vclusterInternalServiceHost,
	}

	return status, nil
}

func (r *UffizziClusterReconciler) upsertVClusterHelmRelease(update bool, ctx context.Context, uCluster *uclusteruffizzicomv1alpha1.UffizziCluster) (*fluxhelmv2beta1.HelmRelease, error) {
	helmReleaseName := BuildVClusterHelmReleaseName(uCluster)
	var (
		VClusterIngressHostname     = BuildVClusterIngressHost(uCluster)
		OutKubeConfigServerArgValue = "https://" + VClusterIngressHostname
	)

	uClusterHelmValues := VCluster{
		Init:    VClusterInit{},
		FsGroup: 12345,
		Ingress: VClusterIngress{
			Enabled: true,
			Host:    VClusterIngressHostname,
			Annotations: map[string]string{
				"app.uffizzi.com/ingress-sync": "true",
			},
		},
		Isolation: VClusterIsolation{
			Enabled:             true,
			PodSecurityStandard: "baseline",
			ResourceQuota: VClusterResourceQuota{
				Enabled: true,
				Quota: VClusterResourceQuotaDefiniton{
					RequestsCpu:                 "2.5",
					RequestsMemory:              "10Gi",
					RequestsEphemeralStorage:    "60Gi",
					RequestsStorage:             "100Gi",
					LimitsCpu:                   "2.5",
					LimitsMemory:                "10Gi",
					LimitsEphemeralStorage:      "160Gi",
					ServicesLoadbalancers:       1,
					ServicesNodePorts:           0,
					CountEndpoints:              40,
					CountConfigmaps:             100,
					CountPersistentVolumeClaims: 20,
					CountPods:                   20,
					CountSecrets:                100,
					CountServices:               20,
				},
			},
			LimitRange: VClusterLimitRange{
				Enabled: true,
				Default: LimitRangeResources{
					Cpu:              "1",
					Memory:           "512Mi",
					EphemeralStorage: "8Gi",
				},
				DefaultRequest: LimitRangeResources{
					Cpu:              "100m",
					Memory:           "128Mi",
					EphemeralStorage: "3Gi",
				},
			},
		},
		NodeSelector: VClusterNodeSelector{
			SandboxGKEIORuntime: "gvisor",
		},
		SecurityContext: VClusterSecurityContext{
			Capabilities: VClusterSecurityContextCapabilities{
				Drop: []string{"all"},
			},
		},
		Tolerations: []VClusterToleration{
			{
				Key:      "sandbox.gke.io/runtime",
				Effect:   "NoSchedule",
				Operator: "Exists",
			},
		},
		Plugin: VClusterPlugins{
			VClusterPlugin{
				Image:           "uffizzi/ucluster-sync-plugin:v0.2.3",
				ImagePullPolicy: "IfNotPresent",
				Rbac: VClusterRbac{
					Role: VClusterRbacRole{
						ExtraRules: []VClusterRbacRule{
							{
								ApiGroups: []string{"networking.k8s.io"},
								Resources: []string{"ingresses"},
								Verbs:     []string{"create", "delete", "patch", "update", "get", "list", "watch"},
							},
						},
					},
					ClusterRole: VClusterRbacClusterRole{
						ExtraRules: []VClusterRbacRule{
							{
								ApiGroups: []string{"apiextensions.k8s.io"},
								Resources: []string{"customresourcedefinitions"},
								Verbs:     []string{"patch", "update", "get", "list", "watch"},
							},
						},
					},
				},
			},
		},
		Syncer: VClusterSyncer{
			KubeConfigContextName: helmReleaseName,
			ExtraArgs: []string{
				"--enforce-toleration=sandbox.gke.io/runtime:NoSchedule",
				"--node-selector=sandbox.gke.io/runtime=gvisor",
				"--enforce-node-selector",
			},
		},
		Sync: VClusterSync{
			Ingresses: EnabledBool{
				Enabled: false,
			},
		},
	}

	if uCluster.Spec.Ingress.Host != "" {
		uClusterHelmValues.Plugin.UffizziClusterSyncPlugin.Env = []VClusterPluginEnv{
			{
				Name:  "VCLUSTER_INGRESS_HOST",
				Value: VClusterIngressHostname,
			},
		}
	}

	if uCluster.Spec.ResourceQuota != nil {
		// map uCluster.Spec.ResourceQuota to uClusterHelmValues.Isolation.ResourceQuota
		q := *uCluster.Spec.ResourceQuota
		qHelmValues := uClusterHelmValues.Isolation.ResourceQuota
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
		uClusterHelmValues.Isolation.ResourceQuota = qHelmValues
	}

	if uCluster.Spec.LimitRange != nil {
		// same for limit range
		lr := uCluster.Spec.LimitRange
		lrHelmValues := uClusterHelmValues.Isolation.LimitRange
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
		uClusterHelmValues.Isolation.LimitRange = lrHelmValues
	}

	//if uCluster.Spec.Ingress.SyncFromManifests != nil {
	//	uClusterHelmValues.Sync.Ingresses.Enabled = *uCluster.Spec.Ingress.SyncFromManifests
	//}

	uClusterHelmValues.Syncer.ExtraArgs = append(uClusterHelmValues.Syncer.ExtraArgs,
		"--tls-san="+VClusterIngressHostname,
		"--out-kube-config-server="+OutKubeConfigServerArgValue,
	)

	if uCluster.Spec.Ingress.Services != nil {
		for _, service := range uCluster.Spec.Ingress.Services {
			uClusterHelmValues.MapServices.FromVirtual = append(uClusterHelmValues.MapServices.FromVirtual, VClusterMapServicesFromVirtual{
				From: service.Namespace + "/" + service.Name,
				To:   helmReleaseName + "-" + service.Name,
			})
		}
	}

	if len(uCluster.Spec.Helm) > 0 {
		uClusterHelmValues.Init.Helm = uCluster.Spec.Helm
	}

	if uCluster.Spec.Manifests != nil {
		uClusterHelmValues.Init.Manifests = *uCluster.Spec.Manifests
	}

	// marshal HelmValues struct to JSON
	helmValuesRaw, err := json.Marshal(uClusterHelmValues)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal HelmValues struct to JSON")
	}

	// Create the apiextensionsv1.JSON instance with the raw data
	helmValuesJSONObj := v1.JSON{Raw: helmValuesRaw}

	// Create a new HelmRelease
	newHelmRelease := &fluxhelmv2beta1.HelmRelease{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      helmReleaseName,
			Namespace: uCluster.Namespace,
		},
		Spec: fluxhelmv2beta1.HelmReleaseSpec{
			Upgrade: &fluxhelmv2beta1.Upgrade{
				Force: false,
			},
			Chart: fluxhelmv2beta1.HelmChartTemplate{
				Spec: fluxhelmv2beta1.HelmChartTemplateSpec{
					Chart:   VCLUSTER_CHART,
					Version: VCLUSTER_CHART_VERSION,
					SourceRef: fluxhelmv2beta1.CrossNamespaceObjectReference{
						Kind:      "HelmRepository",
						Name:      LOFT_HELM_REPO,
						Namespace: uCluster.Namespace,
					},
				},
			},
			ReleaseName: helmReleaseName,
			Values:      &helmValuesJSONObj,
		},
	}

	if err := controllerutil.SetControllerReference(uCluster, newHelmRelease, r.Scheme); err != nil {
		return nil, errors.Wrap(err, "failed to set controller reference")
	}
	// get the helm release spec in string
	newHelmReleaseSpecBytes, err := json.Marshal(newHelmRelease.Spec)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal current spec")
	}
	newHelmReleaseSpec := string(newHelmReleaseSpecBytes)
	// upsert
	if !update && uCluster.Status.LastAppliedHelmReleaseSpec == nil {
		if err := r.Create(ctx, newHelmRelease); err != nil {
			return nil, errors.Wrap(err, "failed to create HelmRelease")
		}
		uCluster.Status.LastAppliedHelmReleaseSpec = &newHelmReleaseSpec
		if err := r.Status().Update(ctx, uCluster); err != nil {
			return nil, errors.Wrap(err, "Failed to update the default UffizziCluster lastAppliedHelmReleaseSpec")
		}

	} else if uCluster.Status.LastAppliedHelmReleaseSpec != nil {
		// create helm release if there is no existing helm release to update
		if update && *uCluster.Status.LastAppliedHelmReleaseSpec != newHelmReleaseSpec {
			if err := r.updateHelmRelease(newHelmRelease, uCluster, ctx); err != nil {
				return nil, errors.Wrap(err, "failed to update HelmRelease")
			}
			return nil, errors.Wrap(err, "couldn't update HelmRelease as LastAppliedHelmReleaseSpec does not exist on resource")
		}
	}

	return newHelmRelease, nil
}

func (r *UffizziClusterReconciler) updateHelmRelease(newHelmRelease *fluxhelmv2beta1.HelmRelease, uCluster *uclusteruffizzicomv1alpha1.UffizziCluster, ctx context.Context) error {
	existingHelmRelease := &fluxhelmv2beta1.HelmRelease{}
	existingHelmReleaseNN := types.NamespacedName{
		Name:      newHelmRelease.Name,
		Namespace: newHelmRelease.Namespace,
	}
	if err := r.Get(ctx, existingHelmReleaseNN, existingHelmRelease); err != nil {
		return errors.Wrap(err, "failed to find HelmRelease")
	}
	// check if the helm release is already progressing, if so, do not update
	if existingHelmRelease.Status.Conditions != nil {
		for _, condition := range existingHelmRelease.Status.Conditions {
			if condition.Type == fluxhelmv2beta1.ReleasedCondition && condition.Status == "Unknown" && condition.Reason == "Progressing" {
				return nil
			}
		}
	}

	newHelmRelease.Spec.Upgrade = &fluxhelmv2beta1.Upgrade{
		Force: true,
	}
	existingHelmRelease.Spec = newHelmRelease.Spec
	if err := r.Update(ctx, existingHelmRelease); err != nil {
		return errors.Wrap(err, "error while updating helm release")
	}
	// update the lastAppliedConfig
	updatedSpecBytes, err := json.Marshal(uCluster.Spec)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal current spec")
	}
	updatedHelmReleaseSpecBytes, err := json.Marshal(existingHelmRelease.Spec)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal current spec")
	}
	updatedSpec := string(updatedSpecBytes)
	updatedHelmReleaseSpec := string(updatedHelmReleaseSpecBytes)
	uCluster.Status.LastAppliedConfiguration = &updatedSpec
	uCluster.Status.LastAppliedHelmReleaseSpec = &updatedHelmReleaseSpec
	if err := r.Status().Update(ctx, uCluster); err != nil {
		return errors.Wrap(err, "Failed to update the default UffizziCluster lastAppliedConfig")
	}
	return nil
}

func (r *UffizziClusterReconciler) createHelmRepo(ctx context.Context, name, namespace, url string) error {
	// Create HelmRepository in the same namespace as the HelmRelease
	helmRepo := &fluxsourcev1.HelmRepository{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: fluxsourcev1.HelmRepositorySpec{
			URL: url,
		},
	}

	err := r.Create(ctx, helmRepo)
	return err
}

func (r *UffizziClusterReconciler) deleteHelmRepo(ctx context.Context, name, namespace string) error {
	// Create HelmRepository in the same namespace as the HelmRelease
	helmRepo := &fluxsourcev1.HelmRepository{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	err := r.Delete(ctx, helmRepo)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *UffizziClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&uclusteruffizzicomv1alpha1.UffizziCluster{}).
		// Watch HelmRelease reconciled by the Helm Controller
		Watches(
			&controllerruntimesource.Kind{Type: &fluxhelmv2beta1.HelmRelease{}},
			&handler.EnqueueRequestForOwner{
				IsController: true,
				OwnerType:    &uclusteruffizzicomv1alpha1.UffizziCluster{},
			}).
		Complete(r)
}
