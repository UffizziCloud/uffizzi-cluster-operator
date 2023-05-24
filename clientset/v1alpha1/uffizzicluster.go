package v1alpha1

import (
	"context"

	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type UffizziClusterInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.UffizziClusterList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.UffizziCluster, error)
	Create(*v1alpha1.UffizziCluster) (*v1alpha1.UffizziCluster, error)
}

type uffizziClusterClient struct {
	restClient rest.Interface
	ns         string
}

func (c *uffizziClusterClient) List(opts metav1.ListOptions) (*v1alpha1.UffizziClusterList, error) {
	result := v1alpha1.UffizziClusterList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("UffizziClusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *uffizziClusterClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.UffizziCluster, error) {
	result := v1alpha1.UffizziCluster{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("UffizziClusters").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *uffizziClusterClient) Create(project *v1alpha1.UffizziCluster) (*v1alpha1.UffizziCluster, error) {
	result := v1alpha1.UffizziCluster{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("UffizziClusters").
		Body(project).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
