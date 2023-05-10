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
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"math/rand"
	"os/exec"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"

	eclusteruffizzicomv1alpha1 "github.com/UffizziCloud/ephemeral-cluster-operator/api/v1alpha1"
	fluxhelmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	fluxsourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
)

// EphemeralClusterReconciler reconciles a EphemeralCluster object
type EphemeralClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Helm values for the vcluster chart
type HelmValuesInit struct {
	Manifests string `json:"manifests"`
}
type HelmValues struct {
	Init HelmValuesInit `json:"init"`
}

const (
	ECLUSTER_NAME_SUFFIX   = "eclus-"
	LOFT_HELM_REPO         = "loft"
	VCLUSTER_CHART         = "vcluster"
	VCLUSTER_CHART_VERSION = "0.15.0"
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
	logger := log.FromContext(ctx)

	fluxYAML, err := getFluxInstallOutput()
	if err != nil {
		fmt.Printf("Error running flux install command: %v\n", err)
		return ctrl.Result{}, err
	}

	vclusterHelmValues := HelmValues{
		Init: HelmValuesInit{
			Manifests: fluxYAML,
		},
	}

	// marshal HelmValues struct to JSON
	helmValuesRaw, err := json.Marshal(vclusterHelmValues)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return ctrl.Result{}, err
	}

	// Create the apiextensionsv1.JSON instance with the raw data
	helmValuesJSONObj := v1.JSON{Raw: helmValuesRaw}

	// For each EphemeralCluster custom resource, check if there is a HelmRelease
	// List all the HelmRelease custom resources and check if there are any with
	// a name that matches the following format:
	// eclus-<EpemeralCluster.Name>-<random-string>
	// if there aren't any, then create a new HelmRelease custom resource
	// with the name eclus-<EpemeralCluster.Name>-<random-string>

	eClusterList := &eclusteruffizzicomv1alpha1.EphemeralClusterList{}
	err = r.List(ctx, eClusterList)
	if err != nil {
		return ctrl.Result{}, err
	}

	helmReleaseList := &fluxhelmv2beta1.HelmReleaseList{}
	err = r.List(ctx, helmReleaseList)
	if err != nil {
		return ctrl.Result{}, err
	}

	for _, eCluster := range eClusterList.Items {
		helmReleaseNameSuffix := ECLUSTER_NAME_SUFFIX + eCluster.Name

		// Check if there is a HelmRelease with the corresponding name
		// If there isn't, then create a new HelmRelease
		if !helmReleaseExists(helmReleaseNameSuffix, helmReleaseList) {
			// Set release name
			helmReleaseName := helmReleaseNameSuffix + "-" + randString(5)

			// Create HelmRepository in the same namespace as the HelmRelease
			loftHelmRepo := &fluxsourcev1.HelmRepository{
				ObjectMeta: ctrl.ObjectMeta{
					Name:      LOFT_HELM_REPO,
					Namespace: eCluster.Namespace,
				},
				Spec: fluxsourcev1.HelmRepositorySpec{
					URL: "https://charts.loft.sh",
				},
			}

			err = r.Create(ctx, loftHelmRepo)
			// check if error is because the HelmRepository already exists
			if err != nil && !strings.Contains(err.Error(), "already exists") {
				logger.Info("Error while creating HelmRepository", "", err)
			}

			// Create a new HelmRelease
			newHelmRelease := &fluxhelmv2beta1.HelmRelease{
				ObjectMeta: ctrl.ObjectMeta{
					Name:      helmReleaseName,
					Namespace: eCluster.Namespace,
				},
				Spec: fluxhelmv2beta1.HelmReleaseSpec{
					Chart: fluxhelmv2beta1.HelmChartTemplate{
						Spec: fluxhelmv2beta1.HelmChartTemplateSpec{
							Chart:   VCLUSTER_CHART,
							Version: VCLUSTER_CHART_VERSION,
							SourceRef: fluxhelmv2beta1.CrossNamespaceObjectReference{
								Kind:      "HelmRepository",
								Name:      LOFT_HELM_REPO,
								Namespace: eCluster.Namespace,
							},
						},
					},
					ReleaseName: helmReleaseName,
					Values:      &helmValuesJSONObj,
				},
			}

			// TODO: implement alternative to SetControllerReference where the owner reference is set
			// in the annotations of the HelmReleases.
			//
			// Set the owner reference for the HelmRelease object
			if err := controllerutil.SetControllerReference(&eCluster, newHelmRelease, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}

			err = r.Create(ctx, newHelmRelease)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

	}

	return ctrl.Result{}, nil
}

func helmReleaseExists(name string, helmReleaseList *fluxhelmv2beta1.HelmReleaseList) bool {
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
func (r *EphemeralClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&eclusteruffizzicomv1alpha1.EphemeralCluster{}).
		// Watch HelmRelease resources
		Watches(&source.Kind{Type: &fluxhelmv2beta1.HelmRelease{}}, &handler.EnqueueRequestForObject{}).
		Watches(&source.Kind{Type: &fluxsourcev1.HelmRepository{}}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
