package controllers

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
)

func BuildVClusterIngressHost(uCluster *v1alpha1.UffizziCluster) string {
	return uCluster.Name + "-" + uCluster.Spec.Ingress.Host
}
