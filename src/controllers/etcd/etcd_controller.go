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

package etcd

import (
	"context"
	uclusteruffizzicomv1alpha1 "github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/controllers/constants"
	fluxhelmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	controllerruntimesource "sigs.k8s.io/controller-runtime/pkg/source"
)

// UffizziClusterReconciler reconciles a UffizziCluster object
type UffizziClusterEtcdReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *UffizziClusterEtcdReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// check if the cluster is a k3s cluster
	// if it is, then we need to create the etcd cluster
	// default lifecycle operation
	logger := log.FromContext(ctx)

	// ----------------------
	// UCLUSTER FIND AND CHECK
	// ----------------------
	// Fetch the UffizziCluster instance in question and then see which kind of event might have been triggered
	uCluster := &uclusteruffizzicomv1alpha1.UffizziCluster{}
	if err := r.Get(ctx, req.NamespacedName, uCluster); err != nil {
		// Handle error
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if uCluster.Spec.ExternalDatastore == constants.ETCD {
		err := r.createBitnamiHelmRepo(ctx, req)
		if err != nil && k8serrors.IsAlreadyExists(err) {
			// logger.Info("Loft Helm Repo for UffizziCluster already exists", "NamespacedName", req.NamespacedName)
		} else {
			logger.Info("Bitnami Helm Repo for UffizziCluster created", "NamespacedName", req.NamespacedName)
		}
		// create a helm release for the etcd cluster
		// check if the helm release exists
		helmRelease := &fluxhelmv2beta1.HelmRelease{}
		err = r.Get(ctx, client.ObjectKey{
			Namespace: uCluster.Namespace,
			Name:      BuildEtcdHelmReleaseName(uCluster),
		}, helmRelease)
		if err != nil {
			// if the helm release does not exist, create it
			_, err = r.upsertETCDHelmRelease(ctx, uCluster)
			if err != nil {
				return ctrl.Result{}, err
			}
		} else {
			return ctrl.Result{}, nil
		}
	}

	return ctrl.Result{}, nil
}

func BuildEtcdHelmReleaseName(uCluster *uclusteruffizzicomv1alpha1.UffizziCluster) string {
	return constants.UCLUSTER_NAME_PREFIX + constants.ETCD + "-" + uCluster.Name
}

// SetupWithManager sets up the controller with the Manager.
func (r *UffizziClusterEtcdReconciler) SetupWithManager(mgr ctrl.Manager) error {
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
