/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/fluxcd/pkg/apis/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type HelmReleaseInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type HelmChartInfo struct {
	Name    string `json:"name"`
	Repo    string `json:"repo"`
	Version string `json:"version,omitempty"`
}

type HelmChart struct {
	Chart   HelmChartInfo   `json:"chart"`
	Values  string          `json:"values,omitempty"`
	Release HelmReleaseInfo `json:"release"`
}

type ExposedVClusterService struct {
	Name                  string            `json:"name"`
	Namespace             string            `json:"namespace"`
	Port                  int32             `json:"port"`
	IngressAnnotations    map[string]string `json:"ingressAnnotations,omitempty"`
	CertManagerTLSEnabled bool              `json:"certManagerTLSEnabled,omitempty"`
}

type ExposedVClusterServiceStatus struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Host      string `json:"host"`
}

type UffizziClusterIngress struct {
	Host     string                   `json:"host,omitempty"`
	Class    string                   `json:"class,omitempty"`
	Services []ExposedVClusterService `json:"services,omitempty"`
}

// UffizziClusterSpec defines the desired state of UffizziCluster
type UffizziClusterSpec struct {
	Ingress    UffizziClusterIngress `json:"ingress,omitempty"`
	Components string                `json:"components,omitempty"`
	TTL        string                `json:"ttl,omitempty"`
	Helm       []HelmChart           `json:"helm,omitempty"`
	Upgrade    bool                  `json:"upgrade,omitempty"`
}

// UffizziClusterStatus defines the observed state of UffizziCluster
type UffizziClusterStatus struct {
	Ready           bool                           `json:"ready"`
	HelmReleaseRef  string                         `json:"helmReleaseRef"`
	KubeConfig      VClusterKubeConfig             `json:"kubeConfig"`
	Host            string                         `json:"host"`
	ExposedServices []ExposedVClusterServiceStatus `json:"exposedServices"`
}

// VClusterKubeConfig is the KubeConfig SecretReference of the related VCluster
type VClusterKubeConfig struct {
	SecretRef meta.SecretKeyReference `json:"secretRef"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Ready",type=boolean,JSONPath=`.status.ready`
//+kubebuilder:printcolumn:name="Host",type=string,JSONPath=`.status.host`

// UffizziCluster is the Schema for the UffizziClusters API
type UffizziCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UffizziClusterSpec   `json:"spec,omitempty"`
	Status UffizziClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// UffizziClusterList contains a list of UffizziCluster
type UffizziClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UffizziCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&UffizziCluster{}, &UffizziClusterList{})
}
