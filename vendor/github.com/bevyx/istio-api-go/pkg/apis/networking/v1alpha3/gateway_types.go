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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TLS modes enforced by the proxy
type Server_TLSOptions_TLSmode string

const (
	// Forward the connection to the upstream server selected based on
	// the SNI string presented by the client.
	Server_TLSOptions_PASSTHROUGH Server_TLSOptions_TLSmode = "PASSTHROUGH"
	// Secure connections with standard TLS semantics.
	Server_TLSOptions_SIMPLE Server_TLSOptions_TLSmode = "SIMPLE"
	// Secure connections to the upstream using mutual TLS by presenting
	// client certificates for authentication.
	Server_TLSOptions_MUTUAL Server_TLSOptions_TLSmode = "MUTUAL"
)

// `Gateway` describes a load balancer operating at the edge of the mesh
// receiving incoming or outgoing HTTP/TCP connections. The specification
// describes a set of ports that should be exposed, the type of protocol to
// use, SNI configuration for the load balancer, etc.
//
// For example, the following Gateway configuration sets up a proxy to act
// as a load balancer exposing port 80 and 9080 (http), 443 (https), and
// port 2379 (TCP) for ingress.  The gateway will be applied to the proxy
// running on a pod with labels `app: my-gateway-controller`. While Istio
// will configure the proxy to listen on these ports, it is the
// responsibility of the user to ensure that external traffic to these
// ports are allowed into the mesh.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: Gateway
// metadata:
//   name: my-gateway
// spec:
//   selector:
//     app: my-gatweway-controller
//   servers:
//   - port:
//       number: 80
//       name: http
//       protocol: HTTP
//     hosts:
//     - uk.bookinfo.com
//     - eu.bookinfo.com
//     tls:
//       httpsRedirect: true # sends 301 redirect for http requests
//   - port:
//       number: 443
//       name: https
//       protocol: HTTPS
//     hosts:
//     - uk.bookinfo.com
//     - eu.bookinfo.com
//     tls:
//       mode: SIMPLE #enables HTTPS on this port
//       serverCertificate: /etc/certs/servercert.pem
//       privateKey: /etc/certs/privatekey.pem
//   - port:
//       number: 9080
//       name: http-wildcard
//       protocol: HTTP
//     hosts:
//     - "*"
//   - port:
//       number: 2379 # to expose internal service via external port 2379
//       name: mongo
//       protocol: MONGO
//     hosts:
//     - "*"
// ```
// The Gateway specification above describes the L4-L6 properties of a load
// balancer. A `VirtualService` can then be bound to a gateway to control
// the forwarding of traffic arriving at a particular host or gateway port.
//
// For example, the following VirtualService splits traffic for
// "https://uk.bookinfo.com/reviews", "https://eu.bookinfo.com/reviews",
// "http://uk.bookinfo.com:9080/reviews",
// "http://eu.bookinfo.com:9080/reviews" into two versions (prod and qa) of
// an internal reviews service on port 9080. In addition, requests
// containing the cookie "user: dev-123" will be sent to special port 7777
// in the qa version. The same rule is also applicable inside the mesh for
// requests to the "reviews.prod.svc.cluster.local" service. This rule is
// applicable across ports 443, 9080. Note that "http://uk.bookinfo.com"
// gets redirected to "https://uk.bookinfo.com" (i.e. 80 redirects to 443).
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: VirtualService
// metadata:
//   name: bookinfo-rule
// spec:
//   hosts:
//   - reviews.prod.svc.cluster.local
//   - uk.bookinfo.com
//   - eu.bookinfo.com
//   gateways:
//   - my-gateway
//   - mesh # applies to all the sidecars in the mesh
//   http:
//   - match:
//     - headers:
//         cookie:
//           user: dev-123
//     route:
//     - destination:
//         port:
//           number: 7777
//         host: reviews.qa.svc.cluster.local
//   - match:
//       uri:
//         prefix: /reviews/
//     route:
//     - destination:
//         port:
//           number: 9080 # can be omitted if its the only port for reviews
//         host: reviews.prod.svc.cluster.local
//       weight: 80
//     - destination:
//         host: reviews.qa.svc.cluster.local
//       weight: 20
// ```
//
// The following VirtualService forwards traffic arriving at (external)
// port 27017 to internal Mongo server on port 5555. This rule is not
// applicable internally in the mesh as the gateway list omits the
// reserved name `mesh`.
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: VirtualService
// metadata:
//   name: bookinfo-Mongo
// spec:
//   hosts:
//   - mongosvr.prod.svc.cluster.local #name of internal Mongo service
//   gateways:
//   - my-gateway
//   tcp:
//   - match:
//     - port: 27017
//     route:
//     - destination:
//         host: mongo.prod.svc.cluster.local
//         port:
//           number: 5555
// ```
type GatewaySpec struct {
	// REQUIRED: A list of server specifications.
	Servers []Server `json:"servers,omitempty"`
	// REQUIRED: One or more labels that indicate a specific set of pods/VMs
	// on which this gateway configuration should be applied.
	// The scope of label search is platform dependent.
	// On Kubernetes, for example, the scope includes pods running in
	// all reachable namespaces.
	Selector map[string]string `json:"selector,omitempty"`
}

// `Server` describes the properties of the proxy on a given load balancer
// port. For example,
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: Gateway
// metadata:
//   name: my-ingress
// spec:
//   selector:
//     app: my-ingress-gateway
//   servers:
//   - port:
//       number: 80
//       name: http2
//       protocol: HTTP2
//     hosts:
//     - "*"
// ```
//
// Another example
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: Gateway
// metadata:
//   name: my-tcp-ingress
// spec:
//   selector:
//     app: my-tcp-ingress-gateway
//   servers:
//   - port:
//       number: 27018
//       name: mongo
//       protocol: MONGO
//     hosts:
//     - "*"
// ```
//
// The following is an example of TLS configuration for port 443
//
// ```yaml
// apiVersion: networking.istio.io/v1alpha3
// kind: Gateway
// metadata:
//   name: my-tls-ingress
// spec:
//   selector:
//     app: my-tls-ingress-gateway
//   servers:
//   - port:
//       number: 443
//       name: https
//       protocol: HTTPS
//     hosts:
//     - "*"
//     tls:
//       mode: SIMPLE
//       serverCertificate: /etc/certs/server.pem
//       privateKey: /etc/certs/privatekey.pem
// ```
type Server struct {
	// REQUIRED: The Port on which the proxy should listen for incoming
	// connections
	Port *Port `json:"port,omitempty"`
	// REQUIRED. A list of hosts exposed by this gateway. At least one
	// host is required. While typically applicable to
	// HTTP services, it can also be used for TCP services using TLS with
	// SNI. May contain a wildcard prefix for the bottom-level component of
	// a domain name. For example `*.foo.com` matches `bar.foo.com`
	// and `*.com` matches `bar.foo.com`, `example.com`, and so on.
	//
	// **Note**: A `VirtualService` that is bound to a gateway must have one
	// or more hosts that match the hosts specified in a server. The match
	// could be an exact match or a suffix match with the server's hosts. For
	// example, if the server's hosts specifies "*.example.com",
	// VirtualServices with hosts dev.example.com, prod.example.com will
	// match. However, VirtualServices with hosts example.com or
	// newexample.com will not match.
	Hosts []string `json:"hosts,omitempty"`
	// Set of TLS related options that govern the server's behavior. Use
	// these options to control if all http requests should be redirected to
	// https, and the TLS modes to use.
	Tls *Server_TLSOptions `json:"tls,omitempty"`
}

type Server_TLSOptions struct {
	// If set to true, the load balancer will send a 301 redirect for all
	// http connections, asking the clients to use HTTPS.
	HttpsRedirect bool `json:"https_redirect,omitempty"`
	// Optional: Indicates whether connections to this port should be
	// secured using TLS. The value of this field determines how TLS is
	// enforced.
	Mode Server_TLSOptions_TLSmode `json:"mode,omitempty"`
	// REQUIRED if mode is `SIMPLE` or `MUTUAL`. The path to the file
	// holding the server-side TLS certificate to use.
	ServerCertificate string `json:"server_certificate,omitempty"`
	// REQUIRED if mode is `SIMPLE` or `MUTUAL`. The path to the file
	// holding the server's private key.
	PrivateKey string `json:"private_key,omitempty"`
	// REQUIRED if mode is `MUTUAL`. The path to a file containing
	// certificate authority certificates to use in verifying a presented
	// client side certificate.
	CaCertificates string `json:"ca_certificates,omitempty"`
	// A list of alternate names to verify the subject identity in the
	// certificate presented by the client.
	SubjectAltNames []string `json:"subject_alt_names,omitempty"`
}

// Port describes the properties of a specific port of a service.
type Port struct {
	// REQUIRED: A valid non-negative integer port number.
	Number uint32 `json:"number,omitempty"`
	// REQUIRED: The protocol exposed on the port.
	// MUST BE one of HTTP|HTTPS|GRPC|HTTP2|MONGO|TCP|TLS.
	// TLS is used to indicate secure connections to non HTTP services.
	Protocol string `json:"protocol,omitempty"`
	// Label assigned to the port.
	Name string `json:"name,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Gateway is the Schema for the Gateways API
// +k8s:openapi-gen=true
type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GatewaySpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GatewayList contains a list of Gateway
type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gateway `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Gateway{}, &GatewayList{})
}