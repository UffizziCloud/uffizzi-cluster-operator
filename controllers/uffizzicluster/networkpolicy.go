package uffizzicluster

import (
	"fmt"
	uclusteruffizzicomv1alpha1 "github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *UffizziClusterReconciler) buildEgressPolicy(uCluster *uclusteruffizzicomv1alpha1.UffizziCluster) *networkingv1.NetworkPolicy {
	port443 := intstr.FromInt(443)
	port80 := intstr.FromInt(80)
	TCP := corev1.ProtocolTCP
	uClusterHelmReleaseName := BuildVClusterHelmReleaseName(uCluster)
	egressPolicy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-workloads-ingress", uClusterHelmReleaseName),
			Namespace: uCluster.Namespace,
		},
		Spec: networkingv1.NetworkPolicySpec{
			Egress: []networkingv1.NetworkPolicyEgressRule{
				{
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Port:     &port443,
							Protocol: &TCP,
						},
						{
							Port:     &port80,
							Protocol: &TCP,
						},
					},
					To: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app.kubernetes.io/component": "controller",
									"app.kubernetes.io/name":      "ingress-nginx",
								},
							},
						},
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"kubernetes.io/metadata.name": "uffizzi",
								},
							},
						},
					},
				},
			},
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"vcluster.loft.sh/managed-by": uClusterHelmReleaseName,
				},
			},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeEgress,
			},
		},
	}
	return egressPolicy
}
