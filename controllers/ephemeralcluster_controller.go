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
	"k8s.io/apimachinery/pkg/runtime"
	"math/rand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"

	eclusteruffizzicomv1alpha1 "github.com/UffizziCloud/ephemeral-cluster-operator/api/v1alpha1"
	helmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
)

// EphemeralClusterReconciler reconciles a EphemeralCluster object
type EphemeralClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	ecluster_name_suffix = "eclus-"
)

//+kubebuilder:rbac:groups=uffizzi.com,resources=ephemeralclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=uffizzi.com,resources=ephemeralclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=uffizzi.com,resources=ephemeralclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the EphemeralCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *EphemeralClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// For each EphemeralCluster custom resource, check if there is a HelmRelease
	// List all the HelmRelease custom resources and check if there are any with
	// a name that matches the following format:
	// eclus-<EpemeralCluster.Name>-<random-string>
	// if there aren't any, then create a new HelmRelease custom resource
	// with the name eclus-<EpemeralCluster.Name>-<random-string>

	eClusterList := &eclusteruffizzicomv1alpha1.EphemeralClusterList{}
	err := r.List(ctx, eClusterList)
	if err != nil {
		return ctrl.Result{}, err
	}

	helmReleaseList := &helmv2beta1.HelmReleaseList{}
	err = r.List(ctx, helmReleaseList)
	if err != nil {
		return ctrl.Result{}, err
	}

	for _, eCluster := range eClusterList.Items {
		correspondingHelmReleaseNameSuffix := ecluster_name_suffix + eCluster.Name

		// Check if there is a HelmRelease with the corresponding name
		// If there isn't, then create a new HelmRelease
		if !helmReleaseExists(correspondingHelmReleaseNameSuffix, helmReleaseList) {
			helmReleaseName := correspondingHelmReleaseNameSuffix + "-" + randString(5)
			// Create a new HelmRelease
			newHelmRelease := &helmv2beta1.HelmRelease{
				ObjectMeta: ctrl.ObjectMeta{
					Name:      helmReleaseName,
					Namespace: eCluster.Namespace,
				},
				Spec: helmv2beta1.HelmReleaseSpec{
					Chart: helmv2beta1.HelmChartTemplate{
						Spec: helmv2beta1.HelmChartTemplateSpec{
							Chart:   "stable/redis",
							Version: "10.6.0",
							SourceRef: helmv2beta1.CrossNamespaceObjectReference{
								Kind:      "HelmRepository",
								Name:      "stable",
								Namespace: "flux-system",
							},
						},
					},
					ReleaseName: helmReleaseName,
					Values:      nil,
				},
			}

			err = r.Create(ctx, newHelmRelease)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

	}

	return ctrl.Result{}, nil
}

func helmReleaseExists(name string, helmReleaseList *helmv2beta1.HelmReleaseList) bool {
	for _, helmRelease := range helmReleaseList.Items {
		if strings.Contains(helmRelease.Name, name) {
			return true
		}
	}
	return false
}

func randString(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// SetupWithManager sets up the controller with the Manager.
func (r *EphemeralClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eclusteruffizzicomv1alpha1.EphemeralCluster{}).
		// Watch HelmRelease resources
		Watches(&source.Kind{Type: &helmv2beta1.HelmRelease{}}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
