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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ScheduledPodAutoscalerSpec defines the desired state of ScheduledPodAutoscaler
type ScheduledPodAutoscalerSpec struct {
	// +kubebuilder:validation:MinLength=0

	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Resource field for ScheduledPodAutoscaler - the resource to scale:
	// Requires two fields - name and type:
	Resource Resource `json:"resource"`

	// Setup for ScaleUp filed
	// Includes two fields - time and value:
	ScaleUp ScaleSpec `json:"scaleUp"`

	// Setup for ScaleDown filed
	// Includes two fields - time and value:
	ScaleDown ScaleSpec `json:"scaleDown"`
}

type Resource struct {
	// name of resource to manage - deployment or HPA name
	Name string `json:"name"`

	// type of resource to manage - options are: deployment,
	// HPA or annotatedDeployment (for HPA-operator managed HPAs),
	// Note (this should default to deployment) :
	// +optional
	Type string `json:"type,omitempty"`
}

type ScaleSpec struct {
	// time of when scaling action to take place:
	Time string `json:"time"`

	// value to scale to:
	Value *int32 `json:"value"`
}

// ScheduledPodAutoscalerStatus defines the observed state of ScheduledPodAutoscaler
type ScheduledPodAutoscalerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Information when was the last time a scaling action was successfully scheduled.
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=spa

// ScheduledPodAutoscaler is the Schema for the scheduledpodautoscalers API
type ScheduledPodAutoscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScheduledPodAutoscalerSpec   `json:"spec,omitempty"`
	Status ScheduledPodAutoscalerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ScheduledPodAutoscalerList contains a list of ScheduledPodAutoscaler
type ScheduledPodAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScheduledPodAutoscaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ScheduledPodAutoscaler{}, &ScheduledPodAutoscalerList{})
}
