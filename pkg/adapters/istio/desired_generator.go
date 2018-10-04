package istio

import (
	istioapi "github.com/bevyx/istio-api-go/pkg/apis/networking/v1alpha3"
	"github.com/bevyx/remesh/pkg/adapters/istio/resources"
	remeshv1alpha1 "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
)

func GetDesiredState(virtualApps []remeshv1alpha1.VirtualApp, namespace string) ([]istioapi.Gateway, []istioapi.VirtualService, []istioapi.DestinationRule) {
	gateways, gatewayVirtualServices, layoutMap, gatewayVirtualServiceNames := extractVirtualApps(virtualApps, namespace)
	transformedServices := TransformLayout(layoutMap)
	transformedVirtualServices := resources.MakeIstioVirtualServices(transformedServices, namespace, gatewayVirtualServiceNames)
	transformedDestinationRules := resources.MakeIstioDestinationRules(transformedServices, namespace)

	virtualServices := append(gatewayVirtualServices, transformedVirtualServices...)

	return gateways, virtualServices, transformedDestinationRules
}

func extractVirtualApps(virtualApps []remeshv1alpha1.VirtualApp, namespace string) ([]istioapi.Gateway, []istioapi.VirtualService, map[string]remeshv1alpha1.LayoutSpec, []string) {
	gateways := make([]istioapi.Gateway, 0)
	virtualServices := make([]istioapi.VirtualService, 0)
	layoutMap := map[string]remeshv1alpha1.LayoutSpec{}
	gatewayVirtualServiceNames := make([]string, 0)
	for _, virtualApp := range virtualApps {
		if doesVappHaveReleases(virtualApp) {
			gateway, gatewayVirtualService, gatewayVirtualServiceName, vappLayoutMap := extractVirtualApp(virtualApp, namespace)
			virtualServices = append(virtualServices, gatewayVirtualService)
			gateways = append(gateways, gateway)
			gatewayVirtualServiceNames = append(gatewayVirtualServiceNames, gatewayVirtualServiceName)
			for key, layout := range vappLayoutMap {
				layoutMap[key] = layout
			}
		}
	}
	return gateways, virtualServices, layoutMap, gatewayVirtualServiceNames
}

func extractVirtualApp(virtualApp remeshv1alpha1.VirtualApp, namespace string) (gateway istioapi.Gateway, gatewayVirtualService istioapi.VirtualService, gatewayVirtualServiceName string, layoutMap map[string]remeshv1alpha1.LayoutSpec) {

	gateway, gatewayName := resources.MakeIstioGateway(virtualApp, namespace)
	httpRoutes := TranslateVirtualAppConfig(virtualApp)
	gatewayVirtualService, gatewayVirtualServiceName = resources.MakeIstioVirtualServiceForGateway(httpRoutes, namespace, gatewayName)
	layoutMap = getLayoutMapFromReleaseFlows(virtualApp.Spec.ReleaseFlows)

	return

}

func doesVappHaveReleases(virtualApp remeshv1alpha1.VirtualApp) bool {
	return len(virtualApp.Spec.ReleaseFlows) > 0
}
