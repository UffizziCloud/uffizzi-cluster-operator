package controllers

import (
	"fmt"
	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	"k8s.io/api/networking/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func BuildVClusterIngress(helmReleaseName string, uCluster *v1alpha1.UffizziCluster) *v1.Ingress {
	uclusterIngressHost := BuildVClusterIngressHost(uCluster)
	ingress := &v1.Ingress{
		ObjectMeta: v12.ObjectMeta{
			Name:      helmReleaseName + "-ingress",
			Namespace: uCluster.Namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS",
				"nginx.ingress.kubernetes.io/ssl-redirect":     "true",
				"nginx.ingress.kubernetes.io/ssl-passthrough":  "true",
			},
		},
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					Host: uclusterIngressHost,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path: "/",
									PathType: func() *v1.PathType {
										pt := v1.PathTypeImplementationSpecific
										return &pt
									}(),
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: helmReleaseName,
											Port: v1.ServiceBackendPort{
												Number: 443,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Add annotations defined for the service. overrides existing annotations
	clusterIngressSpec := uCluster.Spec.Ingress.Cluster
	if len(clusterIngressSpec.IngressAnnotations) > 0 {
		for k, v := range clusterIngressSpec.IngressAnnotations {
			ingress.Annotations[k] = v
		}
	}

	if clusterIngressSpec.CertManagerTLSEnabled {
		ingress.Spec.TLS = []v1.IngressTLS{
			{
				Hosts:      []string{uclusterIngressHost},
				SecretName: fmt.Sprintf("%s-tls", helmReleaseName),
			},
		}
	}

	return ingress
}

func BuildVClusterIngressHost(uCluster *v1alpha1.UffizziCluster) string {
	return uCluster.Name + "-" + uCluster.Spec.Ingress.Host
}
