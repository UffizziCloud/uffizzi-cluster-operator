package v1alpha1

import (
	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

func NewForConfig(c *rest.Config) (*uffizziClusterClient, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &uffizziClusterClient{restClient: client}, nil
}

func (c *uffizziClusterClient) UffizziClusterV1(namespace string) UffizziClusterInterface {
	return &uffizziClusterClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
