//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/fluxcd/pkg/apis/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HelmChart) DeepCopyInto(out *HelmChart) {
	*out = *in
	out.Chart = in.Chart
	out.Release = in.Release
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HelmChart.
func (in *HelmChart) DeepCopy() *HelmChart {
	if in == nil {
		return nil
	}
	out := new(HelmChart)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HelmChartInfo) DeepCopyInto(out *HelmChartInfo) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HelmChartInfo.
func (in *HelmChartInfo) DeepCopy() *HelmChartInfo {
	if in == nil {
		return nil
	}
	out := new(HelmChartInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HelmReleaseInfo) DeepCopyInto(out *HelmReleaseInfo) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HelmReleaseInfo.
func (in *HelmReleaseInfo) DeepCopy() *HelmReleaseInfo {
	if in == nil {
		return nil
	}
	out := new(HelmReleaseInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziCluster) DeepCopyInto(out *UffizziCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziCluster.
func (in *UffizziCluster) DeepCopy() *UffizziCluster {
	if in == nil {
		return nil
	}
	out := new(UffizziCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *UffizziCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziClusterAPIServer) DeepCopyInto(out *UffizziClusterAPIServer) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziClusterAPIServer.
func (in *UffizziClusterAPIServer) DeepCopy() *UffizziClusterAPIServer {
	if in == nil {
		return nil
	}
	out := new(UffizziClusterAPIServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziClusterDistro) DeepCopyInto(out *UffizziClusterDistro) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziClusterDistro.
func (in *UffizziClusterDistro) DeepCopy() *UffizziClusterDistro {
	if in == nil {
		return nil
	}
	out := new(UffizziClusterDistro)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziClusterIngress) DeepCopyInto(out *UffizziClusterIngress) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziClusterIngress.
func (in *UffizziClusterIngress) DeepCopy() *UffizziClusterIngress {
	if in == nil {
		return nil
	}
	out := new(UffizziClusterIngress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziClusterLimitRange) DeepCopyInto(out *UffizziClusterLimitRange) {
	*out = *in
	out.Default = in.Default
	out.DefaultRequest = in.DefaultRequest
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziClusterLimitRange.
func (in *UffizziClusterLimitRange) DeepCopy() *UffizziClusterLimitRange {
	if in == nil {
		return nil
	}
	out := new(UffizziClusterLimitRange)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziClusterLimitRangeDefault) DeepCopyInto(out *UffizziClusterLimitRangeDefault) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziClusterLimitRangeDefault.
func (in *UffizziClusterLimitRangeDefault) DeepCopy() *UffizziClusterLimitRangeDefault {
	if in == nil {
		return nil
	}
	out := new(UffizziClusterLimitRangeDefault)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziClusterLimitRangeDefaultRequest) DeepCopyInto(out *UffizziClusterLimitRangeDefaultRequest) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziClusterLimitRangeDefaultRequest.
func (in *UffizziClusterLimitRangeDefaultRequest) DeepCopy() *UffizziClusterLimitRangeDefaultRequest {
	if in == nil {
		return nil
	}
	out := new(UffizziClusterLimitRangeDefaultRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziClusterList) DeepCopyInto(out *UffizziClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]UffizziCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziClusterList.
func (in *UffizziClusterList) DeepCopy() *UffizziClusterList {
	if in == nil {
		return nil
	}
	out := new(UffizziClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *UffizziClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziClusterRequestsQuota) DeepCopyInto(out *UffizziClusterRequestsQuota) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziClusterRequestsQuota.
func (in *UffizziClusterRequestsQuota) DeepCopy() *UffizziClusterRequestsQuota {
	if in == nil {
		return nil
	}
	out := new(UffizziClusterRequestsQuota)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziClusterResourceCount) DeepCopyInto(out *UffizziClusterResourceCount) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziClusterResourceCount.
func (in *UffizziClusterResourceCount) DeepCopy() *UffizziClusterResourceCount {
	if in == nil {
		return nil
	}
	out := new(UffizziClusterResourceCount)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziClusterResourceQuota) DeepCopyInto(out *UffizziClusterResourceQuota) {
	*out = *in
	out.Requests = in.Requests
	out.Limits = in.Limits
	out.Services = in.Services
	out.Count = in.Count
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziClusterResourceQuota.
func (in *UffizziClusterResourceQuota) DeepCopy() *UffizziClusterResourceQuota {
	if in == nil {
		return nil
	}
	out := new(UffizziClusterResourceQuota)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziClusterResourceQuotaLimits) DeepCopyInto(out *UffizziClusterResourceQuotaLimits) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziClusterResourceQuotaLimits.
func (in *UffizziClusterResourceQuotaLimits) DeepCopy() *UffizziClusterResourceQuotaLimits {
	if in == nil {
		return nil
	}
	out := new(UffizziClusterResourceQuotaLimits)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziClusterServicesQuota) DeepCopyInto(out *UffizziClusterServicesQuota) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziClusterServicesQuota.
func (in *UffizziClusterServicesQuota) DeepCopy() *UffizziClusterServicesQuota {
	if in == nil {
		return nil
	}
	out := new(UffizziClusterServicesQuota)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziClusterSpec) DeepCopyInto(out *UffizziClusterSpec) {
	*out = *in
	out.APIServer = in.APIServer
	out.Ingress = in.Ingress
	if in.Helm != nil {
		in, out := &in.Helm, &out.Helm
		*out = make([]HelmChart, len(*in))
		copy(*out, *in)
	}
	if in.Manifests != nil {
		in, out := &in.Manifests, &out.Manifests
		*out = new(string)
		**out = **in
	}
	if in.ResourceQuota != nil {
		in, out := &in.ResourceQuota, &out.ResourceQuota
		*out = new(UffizziClusterResourceQuota)
		**out = **in
	}
	if in.LimitRange != nil {
		in, out := &in.LimitRange, &out.LimitRange
		*out = new(UffizziClusterLimitRange)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziClusterSpec.
func (in *UffizziClusterSpec) DeepCopy() *UffizziClusterSpec {
	if in == nil {
		return nil
	}
	out := new(UffizziClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UffizziClusterStatus) DeepCopyInto(out *UffizziClusterStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.HelmReleaseRef != nil {
		in, out := &in.HelmReleaseRef, &out.HelmReleaseRef
		*out = new(string)
		**out = **in
	}
	in.KubeConfig.DeepCopyInto(&out.KubeConfig)
	if in.Host != nil {
		in, out := &in.Host, &out.Host
		*out = new(string)
		**out = **in
	}
	if in.LastAppliedConfiguration != nil {
		in, out := &in.LastAppliedConfiguration, &out.LastAppliedConfiguration
		*out = new(string)
		**out = **in
	}
	if in.LastAppliedHelmReleaseSpec != nil {
		in, out := &in.LastAppliedHelmReleaseSpec, &out.LastAppliedHelmReleaseSpec
		*out = new(string)
		**out = **in
	}
	in.LastAwakeTime.DeepCopyInto(&out.LastAwakeTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UffizziClusterStatus.
func (in *UffizziClusterStatus) DeepCopy() *UffizziClusterStatus {
	if in == nil {
		return nil
	}
	out := new(UffizziClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VClusterIngressSpec) DeepCopyInto(out *VClusterIngressSpec) {
	*out = *in
	if in.IngressAnnotations != nil {
		in, out := &in.IngressAnnotations, &out.IngressAnnotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VClusterIngressSpec.
func (in *VClusterIngressSpec) DeepCopy() *VClusterIngressSpec {
	if in == nil {
		return nil
	}
	out := new(VClusterIngressSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VClusterKubeConfig) DeepCopyInto(out *VClusterKubeConfig) {
	*out = *in
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(meta.SecretKeyReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VClusterKubeConfig.
func (in *VClusterKubeConfig) DeepCopy() *VClusterKubeConfig {
	if in == nil {
		return nil
	}
	out := new(VClusterKubeConfig)
	in.DeepCopyInto(out)
	return out
}