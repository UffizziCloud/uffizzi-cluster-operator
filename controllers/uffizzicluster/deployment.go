package uffizzicluster

import (
	"context"
	"errors"
	v1alpha1 "github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	"github.com/UffizziCloud/uffizzi-cluster-operator/controllers/helm/build/vcluster"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

func (r *UffizziClusterReconciler) getUffizziClusterDeployment(ctx context.Context, uCluster *v1alpha1.UffizziCluster) (*appsv1.Deployment, error) {
	ucDeployment := &appsv1.Deployment{}
	if err := r.Get(ctx, types.NamespacedName{
		Name:      vcluster.BuildVClusterHelmReleaseName(uCluster),
		Namespace: uCluster.Namespace}, ucDeployment); err != nil {
		return nil, err
	}
	return ucDeployment, nil
}

// scaleStatefulSets scales the stateful set to the given scale
func (r *UffizziClusterReconciler) scaleDeployments(ctx context.Context, scale int, deployments ...*appsv1.Deployment) error {
	// if the current replicas is greater than 0, then scale down to 0
	replicas := int32(scale)
	for _, d := range deployments {
		d.Spec.Replicas = &replicas
		if err := r.Update(ctx, d); err != nil {
			return err
		}
	}
	return nil
}

// waitForStatefulSetToScale is a goroutine which waits for the stateful set to be ready
func (r *UffizziClusterReconciler) waitForDeploymentToScale(ctx context.Context, scale int, ucDeployment *appsv1.Deployment) error {
	if ucDeployment == nil {
		return errors.New("deployment is nil")
	}
	// wait for the Deployment to be ready
	return wait.PollImmediate(time.Second*5, time.Minute*1, func() (bool, error) {
		if err := r.Get(ctx, types.NamespacedName{
			Name:      ucDeployment.Name,
			Namespace: ucDeployment.Namespace}, ucDeployment); err != nil {
			return false, err
		}
		if ucDeployment.Status.AvailableReplicas == int32(scale) {
			return true, nil
		}
		return false, nil
	})
}
