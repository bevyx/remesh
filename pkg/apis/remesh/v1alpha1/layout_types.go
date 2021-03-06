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

type PortSelector struct {
	Number uint32 `json:"number,omitempty"`
}

type Service struct {
	Host   string            `json:"host,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
	Http   []HTTPRoute       `json:"http,omitempty"`
}

type HTTPRoute struct {
	Match           []HTTPMatchRequest `json:"match,omitempty"`
	DestinationPort *PortSelector      `json:"port,omitempty"`
}

type HTTPMatchRequest struct {
	Uri          *StringMatch           `json:"uri,omitempty"`
	Scheme       *StringMatch           `json:"scheme,omitempty"`
	Method       *StringMatch           `json:"method,omitempty"`
	Authority    *StringMatch           `json:"authority,omitempty"`
	Headers      map[string]StringMatch `json:"headers,omitempty"`
	Port         uint32                 `json:"port,omitempty"`
	SourceLabels map[string]string      `json:"source_labels,omitempty"`
	Gateways     []string               `json:"gateways,omitempty"`
}

type StringMatch struct {
	Exact  string `json:"exact,omitempty"`
	Prefix string `json:"prefix,omitempty"`
	Regex  string `json:"regex,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LayoutSpec defines the desired state of Layout
type LayoutSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Services []Service `json:"services,omitempty"`
}

// LayoutStatus defines the observed state of Layout
type LayoutStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Layout is the Schema for the layouts API
// +k8s:openapi-gen=true
type Layout struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LayoutSpec   `json:"spec,omitempty"`
	Status LayoutStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LayoutList contains a list of Layout
type LayoutList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Layout `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Layout{}, &LayoutList{})
}
