package etcd

import (
	"context"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/constants"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/helm/build"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/helm/build/etcd"
	fluxhelmv2beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	fluxsourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
)

func (r *UffizziClusterEtcdReconciler) createHelmRepo(ctx context.Context, name, namespace, url string) error {
	// Create HelmRepository in the same namespace as the HelmRelease
	helmRepo := &fluxsourcev1.HelmRepository{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: fluxsourcev1.HelmRepositorySpec{
			Interval: metav1.Duration{Duration: time.Minute * 5},
			Timeout:  &metav1.Duration{Duration: time.Second * 60},
			URL:      url,
			Type:     constants.OCI_TYPE,
		},
	}

	err := r.Create(ctx, helmRepo)
	return err
}

func (r *UffizziClusterEtcdReconciler) createBitnamiHelmRepo(ctx context.Context, req ctrl.Request) error {
	return r.createHelmRepo(ctx, constants.BITNAMI_HELM_REPO, req.Namespace, constants.BITNAMI_CHART_REPO_URL)
}

func (r *UffizziClusterEtcdReconciler) upsertETCDHelmRelease(ctx context.Context, uCluster *v1alpha1.UffizziCluster) (*fluxhelmv2beta1.HelmRelease, error) {
	etcdHelmValues := etcd.BuildETCDHelmValues()
	helmValuesJSONObj, err := build.HelmValuesToJSON(etcdHelmValues)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal helm values")
	}

	hr := &fluxhelmv2beta1.HelmRelease{
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
			Values: &helmValuesJSONObj,
		},
	}

	// add UffizziCluster as the owner of this HelmRelease
	if err := ctrl.SetControllerReference(uCluster, hr, r.Scheme); err != nil {
		return nil, errors.Wrap(err, "failed to set controller reference")
	}
	if err := r.Create(ctx, hr); err != nil {
		return nil, errors.Wrap(err, "failed to create HelmRelease")
	}
	return hr, nil
}
