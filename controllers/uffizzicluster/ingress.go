package uffizzicluster

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
)

func BuildVClusterIngressHost(uCluster *v1alpha1.UffizziCluster) string {
	host := ""
	if uCluster.Spec.Ingress.Host != "" {
		host = uCluster.Name + "-" + uCluster.Spec.Ingress.Host
	}
	return host
}
