package uffizzicluster

import (
	context "context"
	"github.com/UffizziCloud/uffizzi-cluster-operator/controllers/constants"
	"github.com/UffizziCloud/uffizzi-cluster-operator/controllers/helm/build/vcluster"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1alpha1 "github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
)

func (r *UffizziClusterReconciler) getUffizziClusterWorkload(ctx context.Context, uCluster *v1alpha1.UffizziCluster) (runtime.Object, error) {
	if uCluster.Spec.Distro == constants.VCLUSTER_K8S_DISTRO {
		return r.getUffizziClusterDeployment(ctx, uCluster)
	}
	return r.getUffizziClusterStatefulSet(ctx, uCluster)
}

// deleteWorkloads deletes all the workloads created by the vcluster
func (r *UffizziClusterReconciler) deleteWorkloads(ctx context.Context, uc *v1alpha1.UffizziCluster) error {
	// delete pods with labels
	podList := &corev1.PodList{}
	if err := r.List(ctx, podList, client.InNamespace(uc.Namespace), client.MatchingLabels(map[string]string{
		constants.VCLUSTER_MANAGED_BY_KEY: vcluster.BuildVClusterHelmReleaseName(uc),
	})); err != nil {
		return err
	}
	for _, pod := range podList.Items {
		if err := r.Delete(ctx, &pod); err != nil {
			return err
		}
	}
	return nil
}
