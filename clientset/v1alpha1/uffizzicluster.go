package v1alpha1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/UffizziCloud/uffizzi-cluster-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type UffizziClusterInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.UffizziClusterList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.UffizziCluster, error)
	Create(UffizziClusterProps) (*v1alpha1.UffizziCluster, error)
	Update(name string, updateClusterProps UpdateUffizziClusterProps) error
	Delete(name string) error
}

type UffizziClusterClient struct {
	restClient rest.Interface
	ns         string
}

type UffizziClusterProps struct {
	Name        string
	Spec        v1alpha1.UffizziClusterSpec
	Annotations map[string]string
}

type UpdateUffizziClusterProps struct {
	Spec        v1alpha1.UffizziClusterSpec
	Annotations map[string]string
}

type JSONPatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func (c *UffizziClusterClient) List(opts metav1.ListOptions) (*v1alpha1.UffizziClusterList, error) {
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

func (c *UffizziClusterClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.UffizziCluster, error) {
	result := v1alpha1.UffizziCluster{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("UffizziClusters").
		Name(name).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *UffizziClusterClient) Create(clusterProps UffizziClusterProps) (*v1alpha1.UffizziCluster, error) {
	uffizziCluster := v1alpha1.UffizziCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "UffizziCluster",
			APIVersion: "uffizzi.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        clusterProps.Name,
			Annotations: clusterProps.Annotations,
		},
		Spec: clusterProps.Spec,
	}

	result := v1alpha1.UffizziCluster{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("UffizziClusters").
		Body(&uffizziCluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *UffizziClusterClient) Update(
	name string,
	updateClusterProps UpdateUffizziClusterProps,
) error {
	uffizziCluster := &v1alpha1.UffizziCluster{}

	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("UffizziClusters").
		Name(name).
		Do(context.TODO()).
		Into(uffizziCluster)

	if err != nil {
		return err
	}

	resourceVersion := uffizziCluster.ObjectMeta.ResourceVersion

	uffizziCluster.Spec = updateClusterProps.Spec
	uffizziCluster.ObjectMeta.Annotations = updateClusterProps.Annotations
	uffizziCluster.TypeMeta = metav1.TypeMeta{
		Kind:       "UffizziCluster",
		APIVersion: "uffizzi.com/v1alpha1",
	}

	fmt.Sprintf("Updated cluster: %v", uffizziCluster)

	updatedJSON, err := json.Marshal(uffizziCluster)
	if err != nil {
		return err
	}

	result := c.restClient.
		Put().
		Namespace(c.ns).
		Resource("UffizziClusters").
		Name(name).
		Body(updatedJSON).
		Param("resourceVersion", resourceVersion).
		SetHeader("Content-Type", "application/json").
		Do(context.TODO())

	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

func (c *UffizziClusterClient) Delete(name string) error {
	result := v1alpha1.UffizziCluster{}
	err := c.restClient.
		Delete().
		Namespace(c.ns).
		Resource("UffizziClusters").
		Name(name).
		Do(context.TODO()).
		Into(&result)

	return err
}
