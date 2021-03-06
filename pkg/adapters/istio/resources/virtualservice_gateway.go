package resources

import (
	istioapi "github.com/bevyx/istio-api-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//MakeIstioVirtualServiceForGateway is
func MakeIstioVirtualServiceForGateway(httpRoutes []istioapi.HTTPRoute, namespace string, gateway string) (virtualService istioapi.VirtualService, virtualServiceName string) {
	virtualServiceName = gateway + VsSuffix
	virtualService = istioapi.VirtualService{
		TypeMeta: metav1.TypeMeta{
			Kind:       "VirtualService",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      virtualServiceName,
			Namespace: namespace,
			/*OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(&virtualappconfig, schema.GroupVersionKind{
					Group:   api.SchemeGroupVersion.Group,
					Version: api.SchemeGroupVersion.Version,
					Kind:    "Layout",
				}),
			},*/
			Labels: AutoGeneratedLabels,
		},
		Spec: istioapi.VirtualServiceSpec{
			Hosts:    []string{HostAll},
			Gateways: []string{gateway},
			Http:     httpRoutes,
		},
	}
	return
}
