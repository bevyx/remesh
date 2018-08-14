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

type Segment struct {
	HttpMatch []HTTPMatchRequest `json:"httpMatch,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TargetingSpec defines the desired state of Targeting
type TargetingSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Entrypoint string  `json:"entrypoint,omitempty"`
	Priority   int32   `json:"priority,omitempty"`
	Segment    Segment `json:"segment,omitempty"`
	Layout     string  `json:"layout,omitempty"`
}

// TargetingStatus defines the observed state of Targeting
type TargetingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Targeting is the Schema for the targetings API
// +k8s:openapi-gen=true
type Targeting struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TargetingSpec   `json:"spec,omitempty"`
	Status TargetingStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TargetingList contains a list of Targeting
type TargetingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Targeting `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Targeting{}, &TargetingList{})
}
