package istio

import (
	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	knativeistio "github.com/knative/serving/pkg/apis/istio/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GenerateIstioGateway(entrypoint api.Entrypoint, namespace string) (knativeistio.Gateway, string) {
	istioServers := make([]knativeistio.Server, 0)

	for _, server := range entrypoint.Spec.Servers {
		var tls *knativeistio.TLSOptions
		if server.Tls != nil {
			tls = &knativeistio.TLSOptions{
				HttpsRedirect:     server.Tls.HttpsRedirect,
				Mode:              knativeistio.TLSMode(server.Tls.Mode),
				ServerCertificate: server.Tls.ServerCertificate,
				PrivateKey:        server.Tls.PrivateKey,
				CaCertificates:    server.Tls.CaCertificates,
				SubjectAltNames:   server.Tls.SubjectAltNames,
			}
		}

		istioServers = append(istioServers, knativeistio.Server{
			Port: knativeistio.Port{
				Name:     server.Port.Name,
				Number:   int(server.Port.Number),
				Protocol: knativeistio.PortProtocol(server.Port.Protocol)},
			Hosts: server.Hosts,
			TLS:   tls,
		})
	}
	gatewaySpec := knativeistio.GatewaySpec{
		Selector: map[string]string{"istio": "ingressgateway"},
		Servers:  istioServers,
	}

	gatewayName := "opsler-" + entrypoint.ObjectMeta.Name

	gateway := createGateway(gatewayName, gatewaySpec, entrypoint)

	return *gateway, gatewayName
}

func createGateway(name string, gatewaySpec knativeistio.GatewaySpec, entrypoint api.Entrypoint) *knativeistio.Gateway {
	labels := map[string]string{}
	return &knativeistio.Gateway{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Gateway",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: entrypoint.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(&entrypoint, schema.GroupVersionKind{
					Group:   api.SchemeGroupVersion.Group,
					Version: api.SchemeGroupVersion.Version,
					Kind:    "Entrypoint",
				}),
			},
			Labels: labels,
		},
		Spec: gatewaySpec,
	}
}
