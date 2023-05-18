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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fluxcd/pkg/apis/meta"
	"github.com/pkg/errors"
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/rand"
	"os/exec"
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
)

// UffizziClusterReconciler reconciles a UffizziCluster object
type UffizziClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type VCluster struct {
	Init            VClusterInit            `json:"init,omitempty"`
	Syncer          VClusterSyncer          `json:"syncer,omitempty"`
	Ingress         VClusterIngress         `json:"ingress,omitempty"`
	FsGroup         int64                   `json:"fsgroup,omitempty"`
	Isolation       VClusterIsolation       `json:"isolation,omitempty"`
	NodeSelector    VClusterNodeSelector    `json:"nodeSelector,omitempty"`
	SecurityContext VClusterSecurityContext `json:"securityContext,omitempty"`
	Tolerations     []VClusterToleration    `json:"tolerations,omitempty"`
	MapServices     VClusterMapServices     `json:"mapServices,omitempty"`
}

// VClusterInit - resources which are created during the init phase of the vcluster
type VClusterInit struct {
	Manifests string                                 `json:"manifests"`
	Helm      []uclusteruffizzicomv1alpha1.HelmChart `json:"helm"`
}

// VClusterSyncer - parameters to create the syncer with
// https://www.vcluster.com/docs/architecture/basics#vcluster-syncer
type VClusterSyncer struct {
	ExtraArgs []string `json:"extraArgs"`
}

// VClusterIngress - parameters to create the ingress with
type VClusterIngress struct {
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

const (
	UCLUSTER_NAME_SUFFIX   = "uc-"
	LOFT_HELM_REPO         = "loft"
	VCLUSTER_CHART         = "vcluster"
	VCLUSTER_CHART_VERSION = "0.15.0"
	LOFT_CHART_REPO_URL    = "https://charts.loft.sh"
	INGRESS_CLASS_NGINX    = "nginx"
)

//+kubebuilder:rbac:groups=uffizzi.com,resources=UffizziClusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=uffizzi.com,resources=UffizziClusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=uffizzi.com,resources=UffizziClusters/finalizers,verbs=update

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
	logger := log.FromContext(ctx)

	// Fetch the UffizziCluster instance
	uCluster := &uclusteruffizzicomv1alpha1.UffizziCluster{}
	if err := r.Get(ctx, req.NamespacedName, uCluster); err != nil {
		// Handle error
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if there is a corresponding HelmRelease which already exists
	// if not create one
	helmReleaseList := &fluxhelmv2beta1.HelmReleaseList{}
	err := r.List(ctx, helmReleaseList)
	if err != nil {
		return ctrl.Result{}, err
	}

	helmReleaseName := UCLUSTER_NAME_SUFFIX + uCluster.Name

	// If there is a VCluster HelmRelease created for this UffizziCluster
	// then update the status of the UffizziCluster based on the status of the HelmRelease
	if helmRelease, exists := helmReleaseExists(helmReleaseName, helmReleaseList); exists {
		logger.Info("HelmRelease exists", "HelmRelease", helmRelease.Name)
		// Update the EphemeralCluster status based on the HelmRelease status
		for _, condition := range helmRelease.Status.Conditions {
			if condition.Type == "Ready" {
				if condition.Status == metav1.ConditionTrue {
					uCluster.Status.Ready = true
					if err := r.Status().Update(ctx, uCluster); err != nil {
						logger.Error(err, "Failed to update UffizziCluster status")
						return ctrl.Result{}, err
					}
				} else {
					logger.Info("UffizziCluster not ready yet", "HelmRelease", helmRelease.Name)
					// Requeue the request to check the status of the HelmRelease
					return ctrl.Result{RequeueAfter: time.Second * 2}, nil
				}
				logger.Info("UffizziCluster in Ready state", "UffizziCluster", uCluster.Name)
				break
			}
		}
		return ctrl.Result{}, err
	}

	err = r.createHelmRepo(ctx, LOFT_HELM_REPO, uCluster.Namespace, LOFT_CHART_REPO_URL)
	// check if error is because the HelmRepository already exists
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		logger.Error(err, "Failed to create HelmRepository")
	}

	newHelmRelease, err := r.createVClusterHelmRelease(ctx, uCluster, helmReleaseName)
	if err != nil {
		logger.Error(err, "Failed to create HelmRelease")
		return ctrl.Result{}, err
	}

	if uCluster.Spec.Ingress.Class == INGRESS_CLASS_NGINX {
		nginxIngressClass := INGRESS_CLASS_NGINX
		clusterIngress := &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      helmReleaseName + "-ingress",
				Namespace: uCluster.Namespace,
				Annotations: map[string]string{
					"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS",
					"nginx.ingress.kubernetes.io/ssl-redirect":     "true",
					"nginx.ingress.kubernetes.io/ssl-passthrough":  "true",
				},
			},
			Spec: networkingv1.IngressSpec{
				IngressClassName: &nginxIngressClass,
				Rules: []networkingv1.IngressRule{
					{
						Host: getUClusterIngressHost(uCluster),
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path: "/",
										PathType: func() *networkingv1.PathType {
											pt := networkingv1.PathTypeImplementationSpecific
											return &pt
										}(),
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: helmReleaseName,
												Port: networkingv1.ServiceBackendPort{
													Number: 443,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		if err := controllerutil.SetControllerReference(uCluster, clusterIngress, r.Scheme); err != nil {
			logger.Error(err, "Failed to create Cluster Ingress")
			return ctrl.Result{}, err
		}

		// create the nginx ingress
		if err := r.Create(ctx, clusterIngress); err != nil {
			if strings.Contains(err.Error(), "already exists") {
				// If the Ingress already exists, update it
				if err := r.Update(ctx, clusterIngress); err != nil {
					logger.Error(err, "Failed to update Cluster Ingress")
					return ctrl.Result{}, err
				}
			} else {
				logger.Error(err, "Failed to create Cluster Ingress")
				return ctrl.Result{}, err
			}
		}

		if uCluster.Spec.Ingress.Services != nil {
			// Create the ingress for each service that is specified in the UffizziCluster services config
			for _, service := range uCluster.Spec.Ingress.Services {
				internalServiceHost := getUClusterInternalServiceIngressHost(uCluster, &service)
				vclusterInternalServiceIngress := &networkingv1.Ingress{
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprintf("%s-%s-%s-ingress", service.Name, service.Namespace, helmReleaseName),
						Namespace: uCluster.Namespace,
						Annotations: map[string]string{
							"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS",
						},
					},
					Spec: networkingv1.IngressSpec{
						IngressClassName: &nginxIngressClass,
						Rules: []networkingv1.IngressRule{
							{
								Host: internalServiceHost,
								IngressRuleValue: networkingv1.IngressRuleValue{
									HTTP: &networkingv1.HTTPIngressRuleValue{
										Paths: []networkingv1.HTTPIngressPath{
											{
												Path: "/",
												PathType: func() *networkingv1.PathType {
													pt := networkingv1.PathTypeImplementationSpecific
													return &pt
												}(),
												Backend: networkingv1.IngressBackend{
													Service: &networkingv1.IngressServiceBackend{
														Name: helmReleaseName + "-" + service.Name,
														Port: networkingv1.ServiceBackendPort{
															Number: service.Port,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}

				// Let the ingresses be owned by UffizziCluster
				if err := controllerutil.SetControllerReference(uCluster, vclusterInternalServiceIngress, r.Scheme); err != nil {
					logger.Error(err, "Failed to set ownerReference for Ingress for internal service "+service.Name+" in namespace "+service.Namespace)
					return ctrl.Result{}, err
				}

				if err := r.Create(ctx, vclusterInternalServiceIngress); err != nil {
					if strings.Contains(err.Error(), "already exists") {
						// If the Ingress already exists, update it
						if err := r.Update(ctx, vclusterInternalServiceIngress); err != nil {
							logger.Error(err, "Failed to update Ingress for internal service "+service.Name+" in namespace "+service.Namespace)
							return ctrl.Result{}, err
						}
					} else {
						logger.Error(err, "Failed to create Ingress for internal service "+service.Name+" in namespace "+service.Namespace)
						return ctrl.Result{}, err
					}
				}

				// add the exposed service to the status
				uCluster.Status.ExposedServices = append(uCluster.Status.ExposedServices,
					uclusteruffizzicomv1alpha1.ExposedVClusterServiceStatus{
						Name:      service.Name,
						Namespace: service.Namespace,
						Host:      internalServiceHost,
					})
			}
		}
	}

	// reference the HelmRelease in the status
	uCluster.Status.HelmReleaseRef = helmReleaseName
	uCluster.Status.KubeConfig.SecretRef = meta.SecretKeyReference{
		Name: "vc-" + helmReleaseName,
	}
	uCluster.Status.Host = getUClusterIngressHost(uCluster)

	if err := r.Status().Update(ctx, uCluster); err != nil {
		logger.Error(err, "Failed to update UffizziCluster status")
		return ctrl.Result{}, err
	}

	logger.Info("Created HelmRelease", "HelmRelease", newHelmRelease.Name)
	// Requeue the request to check the status of the HelmRelease
	return ctrl.Result{Requeue: true}, nil
}

func (r *UffizziClusterReconciler) createVClusterHelmRelease(ctx context.Context, uCluster *uclusteruffizzicomv1alpha1.UffizziCluster, helmReleaseName string) (*fluxhelmv2beta1.HelmRelease, error) {

	fluxYAML, err := getFluxInstallOutput()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get flux install output")
	}

	var (
		TLSSanArgValue              = getUClusterIngressHost(uCluster)
		OutKubeConfigServerArgValue = "https://" + TLSSanArgValue
	)

	uClusterHelmValues := VCluster{
		Init: VClusterInit{
			Manifests: fluxYAML,
		},
		FsGroup: 12345,
		Isolation: VClusterIsolation{
			Enabled:             true,
			PodSecurityStandard: "baseline",
			ResourceQuota: VClusterResourceQuota{
				Quota: VClusterResourceQuotaDefiniton{
					ServicesNodePorts: 5,
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
		Syncer: VClusterSyncer{
			ExtraArgs: []string{
				"--enforce-toleration=sandbox.gke.io/runtime:NoSchedule",
				"--node-selector=sandbox.gke.io/runtime=gvisor",
				"--enforce-node-selector",
			},
		},
	}

	if uCluster.Spec.Ingress.Class == INGRESS_CLASS_NGINX {
		uClusterHelmValues.Ingress.Enabled = true
		uClusterHelmValues.Syncer.ExtraArgs = append(uClusterHelmValues.Syncer.ExtraArgs,
			"--tls-san="+TLSSanArgValue,
			"--out-kube-config-server="+OutKubeConfigServerArgValue,
			//"--out-kube-config-secret="+KubeConfigSecretName,
		)

		if uCluster.Spec.Ingress.Services != nil {
			for _, service := range uCluster.Spec.Ingress.Services {
				uClusterHelmValues.MapServices.FromVirtual = append(uClusterHelmValues.MapServices.FromVirtual, VClusterMapServicesFromVirtual{
					From: service.Namespace + "/" + service.Name,
					To:   helmReleaseName + "-" + service.Name,
				})
			}
		}
	}

	if len(uCluster.Spec.Helm) > 0 {
		uClusterHelmValues.Init.Helm = uCluster.Spec.Helm
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

	err = r.Create(ctx, newHelmRelease)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HelmRelease")
	}

	return newHelmRelease, nil
}

func getUClusterIngressHost(uCluster *uclusteruffizzicomv1alpha1.UffizziCluster) string {
	return uCluster.Name + "-ucluster." + uCluster.Spec.Ingress.Host
}

func getUClusterInternalServiceIngressHost(uCluster *uclusteruffizzicomv1alpha1.UffizziCluster, service *uclusteruffizzicomv1alpha1.ExposedVClusterService) string {
	randomString := rand.String(5)
	vclusterInternalServiceHostExtension := fmt.Sprintf("%s-uc-svc.%s", randomString, uCluster.Spec.Ingress.Host)
	return vclusterInternalServiceHostExtension
}

func (r *UffizziClusterReconciler) createHelmRepo(ctx context.Context, name, namespace, url string) error {
	// Create HelmRepository in the same namespace as the HelmRelease
	loftHelmRepo := &fluxsourcev1.HelmRepository{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: fluxsourcev1.HelmRepositorySpec{
			URL: url,
		},
	}

	err := r.Create(ctx, loftHelmRepo)
	return err
}

func helmReleaseExists(name string, helmReleaseList *fluxhelmv2beta1.HelmReleaseList) (*fluxhelmv2beta1.HelmRelease, bool) {
	for _, helmRelease := range helmReleaseList.Items {
		if strings.Contains(helmRelease.Name, name) {
			return &helmRelease, true
		}
	}
	return nil, false
}

func getFluxInstallOutput() (string, error) {
	cmd := exec.Command("flux", "install", "--namespace=flux-system", "--components=source-controller,helm-controller", "--toleration-keys=sandbox.gke.io/runtime", "--export")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UffizziClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&uclusteruffizzicomv1alpha1.UffizziCluster{}).
		// Watch HelmRelease reconciled by the Helm Controller
		Watches(&controllerruntimesource.Kind{Type: &fluxhelmv2beta1.HelmRelease{}}, &handler.EnqueueRequestForObject{}).
		// Watch HelmRepository reconciled by the Source Controller
		Watches(&controllerruntimesource.Kind{Type: &fluxsourcev1.HelmRepository{}}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
