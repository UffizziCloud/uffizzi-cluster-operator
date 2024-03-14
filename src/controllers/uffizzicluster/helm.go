package uffizzicluster

import (
	"context"
	"encoding/json"
	uclusteruffizzicomv1alpha1 "github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/constants"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/helm/build"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/helm/build/vcluster"
	fluxhelmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	fluxsourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

func (r *UffizziClusterReconciler) createLoftHelmRepo(ctx context.Context, req ctrl.Request) error {
	return r.createHelmRepo(ctx, constants.LOFT_HELM_REPO, req.Namespace, constants.LOFT_CHART_REPO_URL)
}

func (r *UffizziClusterReconciler) deleteLoftHelmRepo(ctx context.Context, req ctrl.Request) error {
	return r.deleteHelmRepo(ctx, constants.LOFT_HELM_REPO, req.Namespace)
}

func (r *UffizziClusterReconciler) upsertVClusterK3SHelmRelease(update bool, ctx context.Context, uCluster *uclusteruffizzicomv1alpha1.UffizziCluster) (*fluxhelmv2beta1.HelmRelease, error) {
	patch := client.MergeFrom(uCluster.DeepCopy())

	vclusterK3sHelmValues, helmReleaseName := vcluster.BuildK3SHelmValues(uCluster)
	helmValuesJSONObj, err := build.HelmValuesToJSON(vclusterK3sHelmValues)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal helm values")
	}

	// Create a new HelmRelease
	newHelmRelease := &fluxhelmv2beta1.HelmRelease{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      helmReleaseName,
			Namespace: uCluster.Namespace,
			Labels: map[string]string{
				constants.UFFIZZI_APP_COMPONENT_LABEL: constants.VCLUSTER,
			},
		},
		Spec: fluxhelmv2beta1.HelmReleaseSpec{
			Upgrade: &fluxhelmv2beta1.Upgrade{
				Force: false,
			},
			Chart: fluxhelmv2beta1.HelmChartTemplate{
				Spec: fluxhelmv2beta1.HelmChartTemplateSpec{
					Chart:   constants.VCLUSTER_CHART_K3S,
					Version: constants.VCLUSTER_CHART_K3S_VERSION,
					SourceRef: fluxhelmv2beta1.CrossNamespaceObjectReference{
						Kind:      "HelmRepository",
						Name:      constants.LOFT_HELM_REPO,
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
		if err := r.Status().Patch(ctx, uCluster, patch); err != nil {
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

func (r *UffizziClusterReconciler) upsertVClusterK8sHelmRelease(update bool, ctx context.Context, uCluster *uclusteruffizzicomv1alpha1.UffizziCluster) (*fluxhelmv2beta1.HelmRelease, error) {
	vclusterK8sHelmValues, helmReleaseName := vcluster.BuildK8SHelmValues(uCluster)
	helmValuesJSONObj, err := build.HelmValuesToJSON(vclusterK8sHelmValues)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal helm values")
	}

	// Create a new HelmRelease
	newHelmRelease := &fluxhelmv2beta1.HelmRelease{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      helmReleaseName,
			Namespace: uCluster.Namespace,
			Labels: map[string]string{
				constants.UFFIZZI_APP_COMPONENT_LABEL: constants.VCLUSTER,
			},
		},
		Spec: fluxhelmv2beta1.HelmReleaseSpec{
			Upgrade: &fluxhelmv2beta1.Upgrade{
				Force: false,
			},
			Chart: fluxhelmv2beta1.HelmChartTemplate{
				Spec: fluxhelmv2beta1.HelmChartTemplateSpec{
					Chart:   constants.VCLUSTER_CHART_K8S,
					Version: constants.VCLUSTER_CHART_K8S_VERSION,
					SourceRef: fluxhelmv2beta1.CrossNamespaceObjectReference{
						Kind:      "HelmRepository",
						Name:      constants.LOFT_HELM_REPO,
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
		patch := client.MergeFrom(uCluster.DeepCopy())
		uCluster.Status.LastAppliedHelmReleaseSpec = &newHelmReleaseSpec
		if err := r.Status().Patch(ctx, uCluster, patch); err != nil {
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
	helmReleasePatch := client.MergeFrom(existingHelmRelease.DeepCopy())
	newHelmRelease.Spec.Upgrade = &fluxhelmv2beta1.Upgrade{
		Force: true,
	}
	existingHelmRelease.Spec = newHelmRelease.Spec
	if err := r.Client.Patch(ctx, existingHelmRelease, helmReleasePatch); err != nil {
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
	uClusterPatch := client.MergeFrom(uCluster.DeepCopy())
	uCluster.Status.LastAppliedConfiguration = &updatedSpec
	uCluster.Status.LastAppliedHelmReleaseSpec = &updatedHelmReleaseSpec
	if err := r.Status().Patch(ctx, uCluster, uClusterPatch); err != nil {
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
			Interval: metav1.Duration{
				Duration: time.Minute * 5,
			},
			Timeout: &metav1.Duration{
				Duration: time.Second * 60,
			},
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
