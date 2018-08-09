package resources

import (
	istioapi "github.com/bevyx/istio-api-go/pkg/apis/networking/v1alpha3"
	istiomodels "github.com/bevyx/remesh/pkg/istio/models"
	"github.com/davecgh/go-spew/spew"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MakeIstioVirtualServices(transformedServices []istiomodels.TransformedService, namespace string, gateway string) []istioapi.VirtualService {
	virtualServices := make([]istioapi.VirtualService, 0)
	for _, transformedService := range transformedServices {
		vs := makeVirtualService(transformedService, namespace, gateway)
		spew.Dump(vs)
		virtualServices = append(virtualServices, makeVirtualService(transformedService, namespace, gateway))
	}
	return virtualServices
}

func makeVirtualService(transformedService istiomodels.TransformedService, namespace string, gateway string) istioapi.VirtualService {
	return istioapi.VirtualService{
		TypeMeta: metav1.TypeMeta{
			Kind:       "VirtualService",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Prefix + transformedService.Host + VsSuffix,
			Namespace: namespace,
			/*OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(&entrypoint, schema.GroupVersionKind{
					Group:   api.SchemeGroupVersion.Group,
					Version: api.SchemeGroupVersion.Version,
					Kind:    "VirtualEnvironment",
				}),
			},*/
			Labels: AutoGeneratedLabels,
		},
		Spec: makeVirtualServiceSpec(transformedService, namespace, gateway),
	}
}

func makeVirtualServiceSpec(transformedService istiomodels.TransformedService, namespace string, gatewayVirtualServiceName string) istioapi.VirtualServiceSpec {
	https := make([]istioapi.HTTPRoute, 0)
	for _, subsetService := range transformedService.ServiceSubsetList {
		for _, virtualEnvironment := range subsetService.VirtualEnvironments {
			https = append(https, istioapi.HTTPRoute{
				Match: []istioapi.HTTPMatchRequest{
					istioapi.HTTPMatchRequest{
						Headers: map[string]istioapi.StringMatch{
							"ol-route": istioapi.StringMatch{
								Exact: virtualEnvironment,
							},
						},
					},
				},
				Route: []istioapi.DestinationWeight{
					istioapi.DestinationWeight{
						Destination: istioapi.Destination{
							Host:   transformedService.Host,
							Subset: GetSubsetName(transformedService.Host, subsetService.SubsetHash),
						},
					},
				},
			})
		}
	}
	return istioapi.VirtualServiceSpec{
		Hosts:    []string{transformedService.Host},
		Gateways: []string{gatewayVirtualServiceName, Mesh},
		Http:     https,
	}
}
