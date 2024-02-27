package uffizzicluster

import (
	context "context"
	"fmt"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/constants"
	"github.com/UffizziCloud/uffizzi-cluster-operator/src/pkg/helm/build/vcluster"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func (r *UffizziClusterReconciler) scaleWorkloads(ctx context.Context, scale int, workloads ...runtime.Object) error {
	sss := []*appsv1.StatefulSet{}
	ds := []*appsv1.Deployment{}
	for _, w := range workloads {
		if w != nil {
			switch w.(type) {
			case *appsv1.StatefulSet:
				ss := w.(*appsv1.StatefulSet)
				sss = append(sss, ss)
			case *appsv1.Deployment:
				d := w.(*appsv1.Deployment)
				ds = append(ds, d)
			}
		}
	}
	if len(sss) > 0 {
		if err := r.scaleStatefulSets(ctx, scale, sss...); err != nil {
			return fmt.Errorf("failed to scale stateful sets: %w", err)
		}
	}
	if len(ds) > 0 {
		if err := r.scaleDeployments(ctx, scale, ds...); err != nil {
			return fmt.Errorf("failed to scale deployments: %w", err)
		}
	}
	return nil
}

func (r *UffizziClusterReconciler) waitForWorkloadToScale(ctx context.Context, scale int, w runtime.Object) error {
	if w != nil {
		switch w.(type) {
		case *appsv1.StatefulSet:
			ss := w.(*appsv1.StatefulSet)
			if err := r.waitForStatefulSetToScale(ctx, scale, ss); err != nil {
				return fmt.Errorf("failed to wait for stateful sets to scale: %w", err)
			}
		case *appsv1.Deployment:
			d := w.(*appsv1.Deployment)
			if err := r.waitForDeploymentToScale(ctx, scale, d); err != nil {
				return fmt.Errorf("failed to wait for deployments to scale: %w", err)
			}
		}
	}
	return nil
}

func (r *UffizziClusterReconciler) getReplicasForWorkload(w runtime.Object) int {
	if w != nil {
		switch w.(type) {
		case *appsv1.StatefulSet:
			ss := w.(*appsv1.StatefulSet)
			return int(*ss.Spec.Replicas)
		case *appsv1.Deployment:
			d := w.(*appsv1.Deployment)
			return int(*d.Spec.Replicas)
		}
	}
	return 0
}
