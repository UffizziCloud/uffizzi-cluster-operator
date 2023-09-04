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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	controllerruntimesource "sigs.k8s.io/controller-runtime/pkg/source"
	"time"

	uclusteruffizzicomv1alpha1 "github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	fluxhelmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// UffizziClusterReconciler reconciles a UffizziCluster object
type UffizziClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	UCLUSTER_NAME_PREFIX       = "uc-"
	LOFT_HELM_REPO             = "loft"
	VCLUSTER_CHART_K3S         = "vcluster"
	VCLUSTER_CHART_K8S         = "vcluster-k8s"
	VCLUSTER_CHART_K3S_VERSION = "0.15.5"
	VCLUSTER_CHART_K8S_VERSION = "0.15.5"
	LOFT_CHART_REPO_URL        = "https://charts.loft.sh"
	VCLUSTER_K3S_DISTRO        = "k3s"
	VCLUSTER_K8S_DISTRO        = "k8s"
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

// add networkpolicy rbac
//+kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete

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

	// ----------------------
	// UCLUSTER INIT and LIFECYCLE OP TYPE determination
	// ----------------------
	// Fetch the UffizziCluster instance in question and then see which kind of event might have been triggered
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
	// if a new ucluster has been created then set the status to have the current ingress spec
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
			helmReleaseRef = ""
			host           = ""
			kubeConfig     = uclusteruffizzicomv1alpha1.VClusterKubeConfig{
				SecretRef: &meta.SecretKeyReference{},
			}
		)
		uCluster.Status = uclusteruffizzicomv1alpha1.UffizziClusterStatus{
			Conditions:     intialConditions,
			HelmReleaseRef: &helmReleaseRef,
			Host:           &host,
			KubeConfig:     kubeConfig,
		}
		if err := r.Status().Update(ctx, uCluster); err != nil {
			logger.Error(err, "Failed to update the default UffizziCluster status")
			return ctrl.Result{}, err
		}
	}

	// ----------------------
	// UCLUSTER HELM CHART and RELATED RESOURCES _CREATION_
	// ----------------------
	// Check if there is already exists a VClusterK3S HelmRelease for this UCluster, if not create one
	helmReleaseName := BuildVClusterHelmReleaseName(uCluster)
	helmRelease := &fluxhelmv2beta1.HelmRelease{}
	var newHelmRelease *fluxhelmv2beta1.HelmRelease
	helmReleaseNamespacedName := client.ObjectKey{
		Namespace: uCluster.Namespace,
		Name:      helmReleaseName,
	}
	// check if the helm release already exists
	err := r.Get(ctx, helmReleaseNamespacedName, helmRelease)
	if err != nil && k8serrors.IsNotFound(err) {
		// create egress policy for vcluster which will allow the vcluster to talk to the outside world
		egressPolicy := r.buildEgressPolicy(uCluster)
		if err := r.Create(ctx, egressPolicy); err != nil {
			logger.Error(err, "Failed to create egress policy")
			return ctrl.Result{Requeue: true}, err
		}
		// helm release does not exist so let's create one
		lifecycleOpType = LIFECYCLE_OP_TYPE_CREATE
		// create either a k8s based vcluster or a k3s based vcluster
		if uCluster.Spec.Distro == VCLUSTER_K8S_DISTRO {
			newHelmRelease, err = r.upsertVClusterK8sHelmRelease(false, ctx, uCluster)
			if err != nil {
				logger.Error(err, "Failed to create HelmRelease")
				return ctrl.Result{Requeue: true}, err
			}
		} else {
			// default to k3s
			newHelmRelease, err = r.upsertVClusterK3sHelmRelease(false, ctx, uCluster)
			if err != nil {
				logger.Error(err, "Failed to create HelmRelease")
				return ctrl.Result{Requeue: true}, err
			}
		}
		// if newHelmRelease is still nil, then the upsert vcluster helm release upsert wasn't concluded
		if newHelmRelease == nil {
			return ctrl.Result{}, nil
		}
		// get the ingress hostname for the vcluster
		vclusterIngressHost := BuildVClusterIngressHost(uCluster) // r.createVClusterIngress(ctx, uCluster)
		uCluster.Status.Host = &vclusterIngressHost
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
			//logger.Error(err, "Failed to update UffizziCluster status")
			return ctrl.Result{RequeueAfter: time.Second * 5}, err
		}
	}
	// create helm repo for loft if it doesn't already exist
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
	// ----------------------
	// UCLUSTER HELM CHART and RELATED RESOURCES _UPDATION_
	// ----------------------
	var updatedHelmRelease *fluxhelmv2beta1.HelmRelease
	if lifecycleOpType == LIFECYCLE_OP_TYPE_UPDATE {
		if currentSpec != lastAppliedSpec {
			if uCluster.Spec.Distro == VCLUSTER_K8S_DISTRO {
				if updatedHelmRelease, err = r.upsertVClusterK8sHelmRelease(true, ctx, uCluster); err != nil {
					logger.Error(err, "Failed to update HelmRelease")
					return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 5}, err
				}
			} else {
				if updatedHelmRelease, err = r.upsertVClusterK3sHelmRelease(true, ctx, uCluster); err != nil {
					logger.Error(err, "Failed to update HelmRelease")
					return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 5}, err
				}
			}
		}
	}
	if updatedHelmRelease == nil {
		return ctrl.Result{}, nil
	}

	// Requeue the request to check the status
	return ctrl.Result{Requeue: true}, nil
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
