// +build !ignore_autogenerated

/*


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

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Resource) DeepCopyInto(out *Resource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Resource.
func (in *Resource) DeepCopy() *Resource {
	if in == nil {
		return nil
	}
	out := new(Resource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScaleSpec) DeepCopyInto(out *ScaleSpec) {
	*out = *in
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScaleSpec.
func (in *ScaleSpec) DeepCopy() *ScaleSpec {
	if in == nil {
		return nil
	}
	out := new(ScaleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduledPodAutoscaler) DeepCopyInto(out *ScheduledPodAutoscaler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduledPodAutoscaler.
func (in *ScheduledPodAutoscaler) DeepCopy() *ScheduledPodAutoscaler {
	if in == nil {
		return nil
	}
	out := new(ScheduledPodAutoscaler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ScheduledPodAutoscaler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduledPodAutoscalerList) DeepCopyInto(out *ScheduledPodAutoscalerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ScheduledPodAutoscaler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduledPodAutoscalerList.
func (in *ScheduledPodAutoscalerList) DeepCopy() *ScheduledPodAutoscalerList {
	if in == nil {
		return nil
	}
	out := new(ScheduledPodAutoscalerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ScheduledPodAutoscalerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduledPodAutoscalerSpec) DeepCopyInto(out *ScheduledPodAutoscalerSpec) {
	*out = *in
	out.Resource = in.Resource
	in.ScaleUp.DeepCopyInto(&out.ScaleUp)
	in.ScaleDown.DeepCopyInto(&out.ScaleDown)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduledPodAutoscalerSpec.
func (in *ScheduledPodAutoscalerSpec) DeepCopy() *ScheduledPodAutoscalerSpec {
	if in == nil {
		return nil
	}
	out := new(ScheduledPodAutoscalerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduledPodAutoscalerStatus) DeepCopyInto(out *ScheduledPodAutoscalerStatus) {
	*out = *in
	if in.LastScheduleTime != nil {
		in, out := &in.LastScheduleTime, &out.LastScheduleTime
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduledPodAutoscalerStatus.
func (in *ScheduledPodAutoscalerStatus) DeepCopy() *ScheduledPodAutoscalerStatus {
	if in == nil {
		return nil
	}
	out := new(ScheduledPodAutoscalerStatus)
	in.DeepCopyInto(out)
	return out
}