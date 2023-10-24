package uffizzicluster

import (
	"context"
	"github.com/UffizziCloud/uffizzi-cluster-operator/controllers/constants"
	fluxsourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *UffizziClusterEtcdReconciler) createHelmRepo(ctx context.Context, name, namespace, url string) error {
	// Create HelmRepository in the same namespace as the HelmRelease
	helmRepo := &fluxsourcev1.HelmRepository{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: fluxsourcev1.HelmRepositorySpec{
			URL:  url,
			Type: "oci",
		},
	}

	err := r.Create(ctx, helmRepo)
	return err
}

func (r *UffizziClusterEtcdReconciler) createBitnamiHelmRepo(ctx context.Context, req ctrl.Request) error {
	return r.createHelmRepo(ctx, constants.BITNAMI_HELM_REPO, req.Namespace, constants.BITNAMI_CHART_REPO_URL)
}
