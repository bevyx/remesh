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

type DestinationRoute struct {
	Host string        `json:"host,omitempty"`
	Port *PortSelector `json:"port,omitempty"`
}

type Service struct {
	Host   string            `json:"host,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

type HTTPRoute struct {
	Match            []*HTTPMatchRequest `json:"match,omitempty"`
	DestinationRoute DestinationRoute    `json:"destinationRoute,omitempty"`
}

type HTTPMatchRequest struct {
	Uri          map[string]string            `json:"uri,omitempty"`
	Scheme       map[string]string            `json:"scheme,omitempty"`
	Method       map[string]string            `json:"method,omitempty"`
	Authority    map[string]string            `json:"authority,omitempty"`
	Headers      map[string]map[string]string `json:"headers,omitempty"`
	Port         uint32                       `json:"port,omitempty"`
	SourceLabels map[string]string            `json:"source_labels,omitempty"`
	Gateways     []string                     `json:"gateways,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VirtualEnvironmentSpec defines the desired state of VirtualEnvironment
type VirtualEnvironmentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Http     []*HTTPRoute `json:"http,omitempty"`
	Services []*Service   `json:"services,omitempty"`
}

// VirtualEnvironmentStatus defines the observed state of VirtualEnvironment
type VirtualEnvironmentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualEnvironment is the Schema for the virtualenvironments API
// +k8s:openapi-gen=true
type VirtualEnvironment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualEnvironmentSpec   `json:"spec,omitempty"`
	Status VirtualEnvironmentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualEnvironmentList contains a list of VirtualEnvironment
type VirtualEnvironmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualEnvironment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualEnvironment{}, &VirtualEnvironmentList{})
}
