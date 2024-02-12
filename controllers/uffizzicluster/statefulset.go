package uffizzicluster

import (
	"context"
	"errors"
	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	uffizzicluster "github.com/UffizziCloud/uffizzi-cluster-operator/controllers/etcd"
	"github.com/UffizziCloud/uffizzi-cluster-operator/controllers/helm/build/vcluster"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

var (
	ErrStatefulSetNil = errors.New("statefulSet is nil")
)

func (r *UffizziClusterReconciler) getUffizziClusterStatefulSet(ctx context.Context, uCluster *v1alpha1.UffizziCluster) (*appsv1.StatefulSet, error) {
	ucStatefulSet := &appsv1.StatefulSet{}
	if err := r.Get(ctx, types.NamespacedName{
		Name:      vcluster.BuildVClusterHelmReleaseName(uCluster),
		Namespace: uCluster.Namespace}, ucStatefulSet); err != nil {
		return nil, err
	}
	return ucStatefulSet, nil
}

func (r *UffizziClusterReconciler) getEtcdStatefulSet(ctx context.Context, uCluster *v1alpha1.UffizziCluster) (*appsv1.StatefulSet, error) {
	etcdStatefulSet := &appsv1.StatefulSet{}
	if err := r.Get(ctx, types.NamespacedName{
		Name:      uffizzicluster.BuildEtcdHelmReleaseName(uCluster),
		Namespace: uCluster.Namespace}, etcdStatefulSet); err != nil {
		return nil, err
	}
	return etcdStatefulSet, nil
}

// scaleStatefulSets scales the stateful set to the given scale
func (r *UffizziClusterReconciler) scaleStatefulSets(ctx context.Context, scale int, statefulSets ...*appsv1.StatefulSet) error {
	// if the current replicas is greater than 0, then scale down to 0
	replicas := int32(scale)
	for _, ss := range statefulSets {
		if ss != nil {
			ss.Spec.Replicas = &replicas
			if err := r.Update(ctx, ss); err != nil {
				return err
			}
		}
	}
	return nil
}

// waitForStatefulSetToScale is a goroutine which waits for the stateful set to be ready
func (r *UffizziClusterReconciler) waitForStatefulSetToScale(ctx context.Context, scale int, ucStatefulSet *appsv1.StatefulSet) error {
	if ucStatefulSet == nil {
		return ErrStatefulSetNil
	}
	// wait for the StatefulSet to be ready
	return wait.PollImmediate(time.Second*5, time.Minute*1, func() (bool, error) {
		if err := r.Get(ctx, types.NamespacedName{
			Name:      ucStatefulSet.Name,
			Namespace: ucStatefulSet.Namespace}, ucStatefulSet); err != nil {
			return false, err
		}
		if ucStatefulSet.Status.AvailableReplicas == int32(scale) {
			return true, nil
		}
		return false, nil
	})
}
