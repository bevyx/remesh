package istio

import (
	"sort"

	istioapi "github.com/bevyx/istio-api-go/pkg/apis/istio/v1alpha3"
	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	"github.com/bevyx/remesh/pkg/models"
)

func MakeRouteForEntrypoint(entrypointFlow models.EntrypointFlow) []istioapi.HTTPRoute {
	istioRouteList := make([]istioapi.HTTPRoute, 0)
	prioritizeTargetingFlows := make([]models.TargetingFlow, len(entrypointFlow.TargetingFlows))
	copy(prioritizeTargetingFlows, entrypointFlow.TargetingFlows)
	sort.Sort(models.ByPriority(prioritizeTargetingFlows))
	for _, targetingFlow := range prioritizeTargetingFlows {
		combainedIstioRouteList := makeCombainedIstioRouteList(targetingFlow.VirtualEnvironment, targetingFlow.Targeting)
		defaultIstioRouteList := makeDefaultIstioRouteList(entrypointFlow.DefaultVirtualEnvironment)

		istioRouteList = append(istioRouteList, combainedIstioRouteList...)
		istioRouteList = append(istioRouteList, defaultIstioRouteList...)
	}

	return istioRouteList
}

func makeCombainedIstioRouteList(virtualEnvironment api.VirtualEnvironment, targeting api.Targeting) []istioapi.HTTPRoute {
	istioRouteList := make([]istioapi.HTTPRoute, 0)
	for _, veRoute := range virtualEnvironment.Spec.Http {
		istioMatchList := make([]istioapi.HTTPMatchRequest, 0)
		for _, veMatch := range veRoute.Match {
			for _, targetingMatch := range targeting.Spec.Segment.HttpMatch {
				istioMatchList = append(istioMatchList, combaineMatchesToIstioMatch(veMatch, targetingMatch))
			}
		}
		istioRouteList = append(istioRouteList, makeIstioRoute(istioMatchList, veRoute.DestinationRoute.Host, veRoute.DestinationRoute.Port.Number, targeting.Spec.VirtualEnvironment))

	}
	return istioRouteList
}

func makeDefaultIstioRouteList(defaultVirtualEnvironment api.VirtualEnvironment) []istioapi.HTTPRoute {
	istioRouteList := make([]istioapi.HTTPRoute, 0)
	for _, veRoute := range defaultVirtualEnvironment.Spec.Http {
		istioMatchList := make([]istioapi.HTTPMatchRequest, 0)
		for _, veMatch := range veRoute.Match {
			istioMatchList = append(istioMatchList, combaineMatchesToIstioMatch(veMatch, api.HTTPMatchRequest{}))
		}
		istioRouteList = append(istioRouteList, makeIstioRoute(istioMatchList, veRoute.DestinationRoute.Host, veRoute.DestinationRoute.Port.Number, defaultVirtualEnvironment.Name))
	}
	return istioRouteList
}

func makeIstioRoute(istioMatchList []istioapi.HTTPMatchRequest, destinationHost string, destinationRoutePortNumber uint32, virtualEnvironmentName string) istioapi.HTTPRoute {
	return istioapi.HTTPRoute{
		Match: istioMatchList,
		Route: []istioapi.DestinationWeight{
			istioapi.DestinationWeight{
				Destination: istioapi.Destination{
					Host: destinationHost,
					Port: istioapi.PortSelector{
						Number: destinationRoutePortNumber,
					},
				},
				Weight: 100,
			},
		},
		AppendHeaders: map[string]string{
			HeaderRouteName: virtualEnvironmentName,
		},
	}
}

func combaineMatchesToIstioMatch(veMatchItem api.HTTPMatchRequest, targetingMatchItem api.HTTPMatchRequest) istioapi.HTTPMatchRequest {
	uri := selectStringMatchesToIstioStringMatch(veMatchItem.Uri, targetingMatchItem.Uri)
	scheme := selectStringMatchesToIstioStringMatch(veMatchItem.Scheme, targetingMatchItem.Scheme)
	method := selectStringMatchesToIstioStringMatch(veMatchItem.Method, targetingMatchItem.Method)
	authority := selectStringMatchesToIstioStringMatch(veMatchItem.Authority, targetingMatchItem.Authority)
	headers := comaineMapOfStringMatchesToIstio(veMatchItem.Headers, targetingMatchItem.Headers)

	port := uint32(0)
	if targetingMatchItem.Port > 0 {
		port = targetingMatchItem.Port
	} else if veMatchItem.Port > 0 {
		port = veMatchItem.Port
	}

	sourceLabels := comaineStringMaps(veMatchItem.SourceLabels, targetingMatchItem.SourceLabels)
	gateways := comaineStringSlicesUnique(veMatchItem.Gateways, targetingMatchItem.Gateways)

	return istioapi.HTTPMatchRequest{
		Uri:          uri,
		Scheme:       scheme,
		Method:       method,
		Authority:    authority,
		Headers:      headers,
		Port:         port,
		SourceLabels: sourceLabels,
		Gateways:     gateways,
	}
}

//TODO: move it to utils
func comaineStringMaps(map1 map[string]string, map2 map[string]string) map[string]string {
	newMap := make(map[string]string, 0)
	for key, value := range map1 {
		newMap[key] = value
	}
	for key, value := range map2 {
		newMap[key] = value
	}
	return newMap
}

//TODO: move it to utils
func comaineStringSlicesUnique(slice1 []string, slice2 []string) []string {
	stringMap := make(map[string]bool, 0)
	for _, v := range slice1 {
		stringMap[v] = true
	}
	for _, v := range slice2 {
		stringMap[v] = true
	}
	newSlice := make([]string, len(stringMap))
	for k := range stringMap {
		newSlice = append(newSlice, k)
	}
	return newSlice
}

func comaineMapOfStringMatchesToIstio(veMap map[string]api.StringMatch, targetingMap map[string]api.StringMatch) map[string]istioapi.StringMatch {
	istioMap := make(map[string]istioapi.StringMatch, 0)
	for keyVe, valueVe := range veMap {
		istioMap[keyVe] = stringMatchToIstioStringMatch(valueVe)
	}
	for keyTargeting, valueTargeting := range targetingMap {
		istioMap[keyTargeting] = stringMatchToIstioStringMatch(valueTargeting)
	}
	return istioMap
}

func selectStringMatchesToIstioStringMatch(veStringMatch api.StringMatch, targetingStringMatch api.StringMatch) istioapi.StringMatch {
	if !isStringMatchEmpty(targetingStringMatch) {
		return stringMatchToIstioStringMatch(targetingStringMatch)
	} else if !isStringMatchEmpty(veStringMatch) {
		return stringMatchToIstioStringMatch(veStringMatch)
	}
	return istioapi.StringMatch{}
}

func isStringMatchEmpty(stringMatch api.StringMatch) bool {
	return stringMatch.Exact == "" && stringMatch.Prefix == "" && stringMatch.Regex == ""
}

func stringMatchToIstioStringMatch(stringMatch api.StringMatch) istioapi.StringMatch {
	return istioapi.StringMatch{
		Exact:  stringMatch.Exact,
		Prefix: stringMatch.Prefix,
		Regex:  stringMatch.Regex,
	}
}
