package istio

import (
	"log"

	istioapi "github.com/bevyx/istio-api-go/pkg/apis/networking/v1alpha3"
	"github.com/bevyx/remesh/pkg/istio/resources"
	"github.com/bevyx/remesh/pkg/models"
)

func Apply(entrypointFlows []models.EntrypointFlow, namespace string) {
	log.Printf("%#v", entrypointFlows)
	istioGateways := make([]istioapi.Gateway, 0)
	istioVirtualServices := make([]istioapi.VirtualService, 0)
	istioDestinationRules := make([]istioapi.DestinationRule, 0)

	for _, entrypointFlow := range entrypointFlows {
		istioGateway, gatewayName := resources.MakeIstioGateway(entrypointFlow.Entrypoint, namespace)
		istioGateways = append(istioGateways, istioGateway)

		httpRoutes := MakeRouteForEntrypoint(entrypointFlow)
		gatewayVirtualService := resources.MakeIstioVirtualServiceForGateway(httpRoutes, namespace, gatewayName)
		istioVirtualServices = append(istioVirtualServices, gatewayVirtualService)

		transformedServices := TransformVirtualEnvironment(entrypointFlow.VirtualEnvironments)
		virtualServices := resources.MakeIstioVirtualServices(transformedServices, namespace, gatewayName)
		destinationRules := resources.MakeIstioDestinationRules(transformedServices, namespace, gatewayName)

		istioVirtualServices = append(istioVirtualServices, virtualServices...)
		istioDestinationRules = append(istioDestinationRules, destinationRules...)
	}

}
