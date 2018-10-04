/*
Copyright 2018 Bevyx.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ReleaseFlow struct {
	ReleaseName string                   `json:"releaseName,omitempty"`
	Release     ReleaseSpec              `json:"release,omitempty"`
	Segments    *map[string]*SegmentSpec `json:"segments,omitempty"`
	LayoutName  string                   `json:"layoutName,omitempty"`
	Layout      *LayoutSpec              `json:"layout,omitempty"`
}

// ByPriority implements sort.Interface for []ReleaseFlow based on
// the Release Priority field.
type ByPriority []ReleaseFlow

func (a ByPriority) Len() int      { return len(a) }
func (a ByPriority) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool {
	if a[i].Release.Targeting == nil {
		return false
	}
	if a[j].Release.Targeting == nil {
		return true
	}
	return a[i].Release.Targeting.Priority > a[j].Release.Targeting.Priority
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VirtualAppSpec defines the desired state of VirtualApp
type VirtualAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	VirtualAppConfig VirtualAppConfigSpec `json:"virtualAppConfig,omitempty"`
	ReleaseFlows     []ReleaseFlow        `json:"releaseFlows,omitempty"`
}

// VirtualAppStatus defines the observed state of VirtualApp
type VirtualAppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualApp is the Schema for the virtualapps API
// +k8s:openapi-gen=true
type VirtualApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualAppSpec   `json:"spec,omitempty"`
	Status VirtualAppStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualAppList contains a list of VirtualApp
type VirtualAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualApp{}, &VirtualAppList{})
}
