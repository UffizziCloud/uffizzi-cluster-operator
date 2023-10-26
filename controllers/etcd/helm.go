package etcd

import (
	"context"
	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/controllers/constants"
	"github.com/UffizziCloud/uffizzi-cluster-operator/controllers/helm/build"
	"github.com/UffizziCloud/uffizzi-cluster-operator/controllers/helm/build/etcd"
	fluxhelmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	fluxsourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			Type: constants.OCI_TYPE,
		},
	}

	err := r.Create(ctx, helmRepo)
	return err
}

func (r *UffizziClusterEtcdReconciler) createBitnamiHelmRepo(ctx context.Context, req ctrl.Request) error {
	return r.createHelmRepo(ctx, constants.BITNAMI_HELM_REPO, req.Namespace, constants.BITNAMI_CHART_REPO_URL)
}

func BuildEtcdHelmRelease(uCluster *v1alpha1.UffizziCluster, helmValuesJSONobj v1.JSON) *fluxhelmv2beta1.HelmRelease {
	return &fluxhelmv2beta1.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      BuildEtcdHelmReleaseName(uCluster),
			Namespace: uCluster.Namespace,
			Labels: map[string]string{
				constants.UFFIZZI_APP_COMPONENT_LABEL: constants.ETCD,
			},
		},
		Spec: fluxhelmv2beta1.HelmReleaseSpec{
			ReleaseName: BuildEtcdHelmReleaseName(uCluster),
			Chart: fluxhelmv2beta1.HelmChartTemplate{
				Spec: fluxhelmv2beta1.HelmChartTemplateSpec{
					Chart:   constants.ETCD_CHART,
					Version: constants.ETCD_CHART_VERSION,
					SourceRef: fluxhelmv2beta1.CrossNamespaceObjectReference{
						Kind: "HelmRepository",
						Name: constants.BITNAMI_HELM_REPO,
					},
				},
			},
			Values: &helmValuesJSONobj,
		},
	}
}

func (r *UffizziClusterEtcdReconciler) upsertETCDHelmRelease(ctx context.Context, uCluster *v1alpha1.UffizziCluster) (*fluxhelmv2beta1.HelmRelease, error) {
	helmValues := etcd.BuildETCDHelmValues()
	helmValuesJSONObj := build.HelmValuesToJSON(helmValues)
	hr := BuildEtcdHelmRelease(uCluster, helmValuesJSONObj)

	// add UffizziCluster as the owner of this HelmRelease
	if err := ctrl.SetControllerReference(uCluster, hr, r.Scheme); err != nil {
		return nil, err
	}
	if err := r.Create(ctx, hr); err != nil {
		return nil, err
	}
	return hr, nil
}
