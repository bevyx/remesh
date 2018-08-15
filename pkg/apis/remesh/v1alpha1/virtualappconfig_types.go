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

type Server_TLSOptions struct {
	HttpsRedirect     bool     `json:"https_redirect,omitempty"`
	Mode              string   `json:"mode,omitempty"`
	ServerCertificate string   `json:"server_certificate,omitempty"`
	PrivateKey        string   `json:"private_key,omitempty"`
	CaCertificates    string   `json:"ca_certificates,omitempty"`
	SubjectAltNames   []string `json:"subject_alt_names,omitempty"`
}

type Port struct {
	Number   uint32 `json:"number,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Name     string `json:"name,omitempty"`
}

type Server struct {
	Port  *Port              `json:"port,omitempty"`
	Hosts []string           `json:"hosts,omitempty"`
	Tls   *Server_TLSOptions `json:"tls,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VirtualAppConfigSpec defines the desired state of VirtualAppConfig
type VirtualAppConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Servers []*Server `json:"servers,omitempty"`
}

// VirtualAppConfigStatus defines the observed state of VirtualAppConfig
type VirtualAppConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualAppConfig is the Schema for the virtualappconfigs API
// +k8s:openapi-gen=true
type VirtualAppConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualAppConfigSpec   `json:"spec,omitempty"`
	Status VirtualAppConfigStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualAppConfigList contains a list of VirtualAppConfig
type VirtualAppConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualAppConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualAppConfig{}, &VirtualAppConfigList{})
}
