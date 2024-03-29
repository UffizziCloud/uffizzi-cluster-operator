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

package uffizzicluster

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/constants"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/helm/build/vcluster"
	"github.com/fluxcd/pkg/apis/meta"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	controllerruntimesource "sigs.k8s.io/controller-runtime/pkg/source"
	"time"

	fluxhelmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// UffizziClusterReconciler reconciles a UffizziCluster object
type UffizziClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

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

// add statefulset rbac
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// add deployments rbac
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// add the ingress rbac
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete

// add networkpolicy rbac
//+kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete

// add pods rbac
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete

// add services rbac
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// add secret rbac
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *UffizziClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		lifecycleOpType constants.LIFECYCLE_OP_TYPE
		currentSpec     string
		lastAppliedSpec string
	)
	// default lifecycle operation
	lifecycleOpType = constants.LIFECYCLE_OP_TYPE_CREATE
	logger := log.FromContext(ctx)

	// ----------------------
	// UCLUSTER INIT and LIFECYCLE OP TYPE determination
	// ----------------------
	// Fetch the UffizziCluster instance in question and then see which kind of event might have been triggered
	uCluster := &v1alpha1.UffizziCluster{}
	if err := r.Get(ctx, req.NamespacedName, uCluster); err != nil {
		// possibly a delete event
		if k8serrors.IsNotFound(err) {
			lifecycleOpType = constants.LIFECYCLE_OP_TYPE_DELETE
			logger.Info("UffizziCluster deleted", "NamespacedName", req.NamespacedName)
		} else {
			logger.Info("Failed to get UffizziCluster", "NamespacedName", req.NamespacedName, "Error", err)
		}
		// Handle error
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// if a new ucluster has been created then set the status to have the current ingress spec
	if lifecycleOpType == constants.LIFECYCLE_OP_TYPE_CREATE {
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
				Initializing(),
				InitializingAPI(),
				InitializingDataStore(),
				DefaultSleepState(),
			}
			helmReleaseRef = ""
			host           = ""
			kubeConfig     = v1alpha1.VClusterKubeConfig{
				SecretRef: &meta.SecretKeyReference{},
			}
			lastAwakeTime = metav1.Now().Rfc3339Copy()
		)
		patch := client.MergeFrom(uCluster.DeepCopy())
		uCluster.Status = v1alpha1.UffizziClusterStatus{
			Conditions:     intialConditions,
			HelmReleaseRef: &helmReleaseRef,
			Host:           &host,
			KubeConfig:     kubeConfig,
			LastAwakeTime:  lastAwakeTime,
		}
		if err := r.Status().Patch(ctx, uCluster, patch); err != nil {
			logger.Error(err, "Failed to update the default UffizziCluster status")
			return ctrl.Result{}, err
		}
	}

	// ----------------------
	// UCLUSTER HELM CHART and RELATED RESOURCES _CREATION_
	// ----------------------
	// Check if there is already exists a VClusterK3S HelmRelease for this UCluster, if not create one
	helmReleaseName := vcluster.BuildVClusterHelmReleaseName(uCluster)
	helmRelease := &fluxhelmv2beta1.HelmRelease{}
	// check if the helm release already exists
	// create one if helmRelease doesn't exist
	var newHelmRelease *fluxhelmv2beta1.HelmRelease
	helmReleaseNamespacedName := client.ObjectKey{
		Namespace: uCluster.Namespace,
		Name:      helmReleaseName,
	}
	err := r.Get(ctx, helmReleaseNamespacedName, helmRelease)
	// if the helm release does not exist, create it
	if err != nil && k8serrors.IsNotFound(err) {
		// create egress policy for vcluster which will allow the vcluster to talk to the outside world
		if result, createEgressError := r.createEgressPolicy(ctx, uCluster); createEgressError != nil {
			return result, createEgressError
		}
		// helm release does not exist so let's create one
		lifecycleOpType = constants.LIFECYCLE_OP_TYPE_CREATE
		// create either a k8s based vcluster or a k3s based vcluster
		newHelmRelease, result, err := r.createVClusterHelmChart(ctx, uCluster, newHelmRelease)
		if err != nil {
			return result, err
		}
		// if newHelmRelease is still nil, then the upsert vcluster helm release upsert wasn't concluded
		if newHelmRelease == nil {
			return ctrl.Result{}, nil
		}
		// get the ingress hostname for the vcluster
		vclusterIngressHost := vcluster.BuildVClusterIngressHost(uCluster) // r.createVClusterIngress(ctx, uCluster)
		patch := client.MergeFrom(uCluster.DeepCopy())
		uCluster.Status.Host = &vclusterIngressHost
		// reference the HelmRelease in the status
		uCluster.Status.HelmReleaseRef = &helmReleaseName
		uCluster.Status.KubeConfig.SecretRef = &meta.SecretKeyReference{
			Name: "vc-" + helmReleaseName,
		}
		if err := r.Status().Patch(ctx, uCluster, patch); err != nil {
			logger.Error(err, "Failed to update UffizziCluster status")
			return ctrl.Result{RequeueAfter: time.Second * 5}, err
		}
		logger.Info("Created HelmRelease", "HelmRelease", newHelmRelease.Name)
	} else if err != nil {
		logger.Error(err, "Failed to create HelmRelease, unknown error")
		return ctrl.Result{}, err
	} else {
		lifecycleOpType = constants.LIFECYCLE_OP_TYPE_UPDATE
		// if helm release already exists then replicate the status conditions onto the uffizzicluster object
		patch := client.MergeFrom(uCluster.DeepCopy())
		mirrorHelmStackConditions(helmRelease, uCluster)
		if err := r.Status().Patch(ctx, uCluster, patch); err != nil {
			//logger.Error(err, "Failed to update UffizziCluster status")
			return ctrl.Result{RequeueAfter: time.Second * 5}, err
		}
	}
	// create helm repo for loft if it doesn't already exist
	if lifecycleOpType == constants.LIFECYCLE_OP_TYPE_CREATE {
		err := r.createLoftHelmRepo(ctx, req)
		if err != nil && k8serrors.IsAlreadyExists(err) {
			// logger.Info("Loft Helm Repo for UffizziCluster already exists", "NamespacedName", req.NamespacedName)
		} else {
			logger.Info("Loft Helm Repo for UffizziCluster created", "NamespacedName", req.NamespacedName)
		}
		patch := client.MergeFrom(uCluster.DeepCopy())
		uCluster.Status.LastAppliedConfiguration = &currentSpec
		if err := r.Status().Patch(ctx, uCluster, patch); err != nil {
			logger.Error(err, "Failed to update the default UffizziCluster lastAppliedConfig")
			return ctrl.Result{}, err
		}
		logger.Info("UffizziCluster lastAppliedConfig has been set")
	} else if lifecycleOpType == constants.LIFECYCLE_OP_TYPE_DELETE {
		err := r.deleteLoftHelmRepo(ctx, req)
		if err != nil && k8serrors.IsNotFound(err) {
			logger.Info("Loft Helm Repo for UffizziCluster already deleted", "NamespacedName", req.NamespacedName)
		}
	}
	// ----------------------
	// UCLUSTER HELM CHART and RELATED RESOURCES _UPDATION_
	// ----------------------
	var updatedHelmRelease *fluxhelmv2beta1.HelmRelease
	if lifecycleOpType == constants.LIFECYCLE_OP_TYPE_UPDATE {
		if currentSpec != lastAppliedSpec {
			if uCluster.Spec.Distro == constants.VCLUSTER_K8S_DISTRO {
				if updatedHelmRelease, err = r.upsertVClusterK8sHelmRelease(true, ctx, uCluster); err != nil {
					logger.Error(err, "Failed to update HelmRelease")
					return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 5}, err
				}
			} else {
				if updatedHelmRelease, err = r.upsertVClusterK3SHelmRelease(true, ctx, uCluster); err != nil {
					logger.Info("Failed to update HelmRelease with error, reconciling", "Error", err.Error())
					return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 5}, err
				}
			}
		}
	}
	if updatedHelmRelease == nil {
		return ctrl.Result{}, nil
	}
	// ----------------------
	// UCLUSTER SLEEP
	// ----------------------
	if err := r.reconcileSleepState(ctx, uCluster); err != nil {
		if k8serrors.IsNotFound(err) {
			logger.Info("vcluster workload not found, will check again in the next round")
			return ctrl.Result{}, nil
		}
		// cluster did not sleep
		logger.Info("Failed to reconcile sleep state, reconciling again", "Error", err.Error())
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *UffizziClusterReconciler) createVClusterHelmChart(ctx context.Context, uCluster *v1alpha1.UffizziCluster, newHelmRelease *fluxhelmv2beta1.HelmRelease) (*fluxhelmv2beta1.HelmRelease, ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var err error = nil
	if uCluster.Spec.Distro == constants.VCLUSTER_K8S_DISTRO {
		newHelmRelease, err = r.upsertVClusterK8sHelmRelease(false, ctx, uCluster)
		if err != nil {
			logger.Error(err, "Failed to create HelmRelease")
			return nil, ctrl.Result{Requeue: true}, err
		}
	} else {
		// default to k3s
		newHelmRelease, err = r.upsertVClusterK3SHelmRelease(false, ctx, uCluster)
		if err != nil {
			logger.Error(err, "Failed to create HelmRelease")
			return nil, ctrl.Result{Requeue: true}, err
		}
	}
	return newHelmRelease, ctrl.Result{}, nil
}

func (r *UffizziClusterReconciler) createEgressPolicy(ctx context.Context, uCluster *v1alpha1.UffizziCluster) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	egressPolicy := r.buildEgressPolicy(uCluster)
	if err := controllerutil.SetControllerReference(uCluster, egressPolicy, r.Scheme); err != nil {
		return ctrl.Result{Requeue: true}, errors.Wrap(err, "failed to set controller reference")
	}
	if err := r.Create(ctx, egressPolicy); err != nil {
		logger.Error(err, "Failed to create egress policy")
		return ctrl.Result{Requeue: true}, err
	}
	return ctrl.Result{}, nil
}

// reconcileSleepState reconciles the sleep state of the vcluster
// it also makes sure that the vcluster is up and running before setting the sleep state
func (r *UffizziClusterReconciler) reconcileSleepState(ctx context.Context, uCluster *v1alpha1.UffizziCluster) error {
	var (
		// get the patch copy of the uCluster so that we can have a good diff between the uCluster and patch object
		patch           = client.MergeFrom(uCluster.DeepCopy())
		ucWorkload      runtime.Object
		etcdStatefulSet *appsv1.StatefulSet
		err             error
	)
	// get the stateful set or deployment created by the helm chart
	// k3s is a statefulset, k8s is a deployment
	if ucWorkload, err = r.getUffizziClusterWorkload(ctx, uCluster); err != nil {
		return fmt.Errorf("failed to get uffizzicluster workload: %w", err)
	}
	// get the etcd stateful set created by the helm chart
	if uCluster.Spec.ExternalDatastore == constants.ETCD {
		if etcdStatefulSet, err = r.getEtcdStatefulSet(ctx, uCluster); err != nil {
			return fmt.Errorf("failed to get etcd statefulset: %w", err)
		}
	}
	// execute sleep reconciliation based on the type of workload
	currentReplicas := r.getReplicasForWorkload(ucWorkload)
	// scale the vcluster instance to 0 if the sleep flag is true
	if uCluster.Spec.Sleep && currentReplicas > 0 {
		if err = r.scaleWorkloads(ctx, 0, ucWorkload); err != nil {
			return err
		}
		if err = r.waitForWorkloadToScale(ctx, 0, ucWorkload); err == nil {
			setCondition(uCluster, APINotReady())
		}
		if uCluster.Spec.ExternalDatastore == constants.ETCD {
			if err = r.scaleWorkloads(ctx, 0, etcdStatefulSet); err != nil {
				return err
			}
			if err = r.waitForWorkloadToScale(ctx, 0, etcdStatefulSet); err == nil {
				setCondition(uCluster, DataStoreNotReady())
			}
		} else {
			setCondition(uCluster, DataStoreNotReady())
		}
		if err != nil {
			return err
		}
		err = r.deleteWorkloads(ctx, uCluster)
		if err != nil {
			return err
		}
		sleepingTime := metav1.Now().Rfc3339Copy()
		setCondition(uCluster, Sleeping(sleepingTime))
		// if the current replicas is 0, then do nothing
	} else if !uCluster.Spec.Sleep && currentReplicas == 0 {
		workloads := []runtime.Object{}
		workloads = append(workloads, ucWorkload)
		if uCluster.Spec.ExternalDatastore == constants.ETCD {
			workloads = append(workloads, etcdStatefulSet)
		}
		if err := r.scaleWorkloads(ctx, 1, workloads...); err != nil {
			return err
		}
	}
	// ensure that the statefulset is up if the cluster is not sleeping
	if !uCluster.Spec.Sleep {
		// set status for vcluster waking up
		lastAwakeTime := metav1.Now().Rfc3339Copy()
		uCluster.Status.LastAwakeTime = lastAwakeTime
		// if the above runs successfully, then set the status to awake
		setCondition(uCluster, Awoken(lastAwakeTime))
		var err error
		if uCluster.Spec.ExternalDatastore == constants.ETCD {
			if err = r.waitForWorkloadToScale(ctx, 1, etcdStatefulSet); err == nil {
				setCondition(uCluster, DataStoreReady())
			}
		} else {
			setCondition(uCluster, DataStoreReady())
		}
		if err = r.waitForWorkloadToScale(ctx, 1, ucWorkload); err == nil {
			setCondition(uCluster, APIReady())
		}
		if err != nil {
			return err
		}
	}

	if err := r.Status().Patch(ctx, uCluster, patch); err != nil {
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UffizziClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.UffizziCluster{}).
		// Watch HelmRelease reconciled by the Helm Controller
		Watches(
			&controllerruntimesource.Kind{Type: &fluxhelmv2beta1.HelmRelease{}},
			&handler.EnqueueRequestForOwner{
				IsController: true,
				OwnerType:    &v1alpha1.UffizziCluster{},
			}).
		Complete(r)
}
