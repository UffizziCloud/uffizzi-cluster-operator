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

type VClusterIngressSpec struct {
	IngressAnnotations    map[string]string `json:"ingressAnnotations,omitempty"`
	CertManagerTLSEnabled bool              `json:"certManagerTLSEnabled,omitempty"`
}

// UffiClusterIngress defines the ingress capabilities of the cluster,
// the basic host can be setup for all
type UffizziClusterIngress struct {
	Host  string `json:"host,omitempty"`
	Class string `json:"class,omitempty"`
	//+kubebuilder:default:=true
	SyncFromManifests *bool                    `json:"syncFromManifests,omitempty"`
	Cluster           VClusterIngressSpec      `json:"cluster,omitempty"`
	Services          []ExposedVClusterService `json:"services,omitempty"`
}

type UffizziClusterStorageSync struct {
	PersistentVolumes      bool `json:"persistentVolumes,omitempty"`
	StorageClasses         bool `json:"storageClasses,omitempty"`
	PersistentVolumeClaims bool `json:"persistentVolumeClaims,omitempty"`
}

// UffizziClusterSpec defines the desired state of UffizziCluster
type UffizziClusterSpec struct {
	Ingress   UffizziClusterIngress `json:"ingress,omitempty"`
	TTL       string                `json:"ttl,omitempty"`
	Helm      []HelmChart           `json:"helm,omitempty"`
	Manifests *string               `json:"manifests,omitempty"`
}

// UffizziClusterStatus defines the observed state of UffizziCluster
type UffizziClusterStatus struct {
	Conditions               []metav1.Condition             `json:"conditions,omitempty"`
	HelmReleaseRef           *string                        `json:"helmReleaseRef,omitempty"`
	KubeConfig               VClusterKubeConfig             `json:"kubeConfig,omitempty"`
	Host                     *string                        `json:"host,omitempty"`
	ExposedServices          []ExposedVClusterServiceStatus `json:"exposedServices,omitempty"`
	LastAppliedConfiguration *string                        `json:"lastAppliedConfiguration,omitempty"`
}

// VClusterKubeConfig is the KubeConfig SecretReference of the related VCluster
type VClusterKubeConfig struct {
	SecretRef *meta.SecretKeyReference `json:"secretRef,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=uc;ucluster
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=='Ready')].status`
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
