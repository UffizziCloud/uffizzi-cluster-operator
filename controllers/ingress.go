package controllers

import (
	"fmt"
	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	"k8s.io/api/networking/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
)

func BuildVClusterIngress(helmReleaseName string, uCluster *v1alpha1.UffizziCluster) *v1.Ingress {
	nginxIngressClass := INGRESS_CLASS_NGINX
	return &v1.Ingress{
		ObjectMeta: v12.ObjectMeta{
			Name:      helmReleaseName + "-ingress",
			Namespace: uCluster.Namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS",
				"nginx.ingress.kubernetes.io/ssl-redirect":     "true",
				"nginx.ingress.kubernetes.io/ssl-passthrough":  "true",
				"cert-manager.io/cluster-issuer":               "my-uffizzi-letsencrypt",
			},
		},
		Spec: v1.IngressSpec{
			IngressClassName: &nginxIngressClass,
			Rules: []v1.IngressRule{
				{
					Host: BuildVClusterIngressHost(uCluster),
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
}

func BuildVClusterInternalServiceIngress(service v1alpha1.ExposedVClusterService, uCluster *v1alpha1.UffizziCluster, helmReleaseName string) *v1.Ingress {
	const PATH = "/"
	nginxIngressClass := INGRESS_CLASS_NGINX
	internalServiceHost := BuildVClusterInternalServiceIngressHost(uCluster)
	ingress := &v1.Ingress{
		ObjectMeta: v12.ObjectMeta{
			Name:        fmt.Sprintf("%s-x-%s-x-%s", service.Name, service.Namespace, helmReleaseName),
			Namespace:   uCluster.Namespace,
			Annotations: map[string]string{},
		},
		Spec: v1.IngressSpec{
			IngressClassName: &nginxIngressClass,
			Rules: []v1.IngressRule{
				{
					Host: internalServiceHost,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path: PATH,
									PathType: func() *v1.PathType {
										pt := v1.PathTypeImplementationSpecific
										return &pt
									}(),
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: helmReleaseName + "-" + service.Name,
											Port: v1.ServiceBackendPort{
												Number: service.Port,
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
	if len(service.IngressAnnotations) > 0 {
		for k, v := range service.IngressAnnotations {
			ingress.Annotations[k] = v
		}
	}

	if service.CertManagerTLSEnabled {
		ingress.Spec.TLS = []v1.IngressTLS{
			{
				Hosts:      []string{internalServiceHost},
				SecretName: fmt.Sprintf("%s-x-%s-x-%s-x-tls", service.Name, service.Namespace, helmReleaseName),
			},
		}
	}

	return ingress
}

func BuildVClusterIngressHost(uCluster *v1alpha1.UffizziCluster) string {
	return uCluster.Name + ".uc." + uCluster.Spec.Ingress.Host
}

func BuildVClusterInternalServiceIngressHost(uCluster *v1alpha1.UffizziCluster) string {
	randomString := rand.String(5)
	vclusterInternalServiceHostExtension := fmt.Sprintf("%s-svc.uc.%s", randomString, uCluster.Spec.Ingress.Host)
	return vclusterInternalServiceHostExtension
}
