package istio

import (
	istioapi "github.com/bevyx/istio-api-go/pkg/apis/networking/v1alpha3"
	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
)

//MakeRouteForVirtualAppConfig is
func MakeRouteForVirtualAppConfig(virtualApp api.VirtualApp) []istioapi.HTTPRoute {
	istioRouteList := make([]istioapi.HTTPRoute, 0)
	for _, releaseFlow := range virtualApp.Spec.ReleaseFlows {
		if releaseFlow.Layout != nil {
			if releaseFlow.Segments == nil {
				defaultIstioRouteList := makeLayoutIstioRouteList(releaseFlow.LayoutName, *releaseFlow.Layout)
				istioRouteList = append(istioRouteList, defaultIstioRouteList...)
			} else {
				combainedIstioRouteList := makeCombainedIstioRouteList(releaseFlow.LayoutName, *releaseFlow.Layout, *releaseFlow.Segments)
				istioRouteList = append(istioRouteList, combainedIstioRouteList...)
			}
		}

	}

	return istioRouteList
}

func makeCombainedIstioRouteList(layoutName string, layout api.LayoutSpec, segments map[string]*api.SegmentSpec) []istioapi.HTTPRoute {
	istioRouteList := make([]istioapi.HTTPRoute, 0)
	for _, segment := range segments {
		if segment != nil {
			for _, veRoute := range layout.Http {
				istioMatchList := make([]istioapi.HTTPMatchRequest, 0)
				for _, veMatch := range veRoute.Match {
					for _, targetingMatch := range segment.HttpMatch {
						istioMatchList = append(istioMatchList, combineMatchesToIstioMatch(veMatch, targetingMatch))
					}
				}
				istioRouteList = append(istioRouteList, makeIstioRoute(istioMatchList, veRoute.Destination.Host, veRoute.Destination.Port, layoutName))
			}
		}
	}

	return istioRouteList
}

func makeLayoutIstioRouteList(layoutName string, layout api.LayoutSpec) []istioapi.HTTPRoute {
	istioRouteList := make([]istioapi.HTTPRoute, 0)
	for _, veRoute := range layout.Http {
		istioMatchList := make([]istioapi.HTTPMatchRequest, 0)
		for _, veMatch := range veRoute.Match {
			istioMatchList = append(istioMatchList, combineMatchesToIstioMatch(veMatch, api.HTTPMatchRequest{}))
		}
		istioRouteList = append(istioRouteList, makeIstioRoute(istioMatchList, veRoute.Destination.Host, veRoute.Destination.Port, layoutName))
	}
	return istioRouteList
}

func makeIstioRoute(istioMatchList []istioapi.HTTPMatchRequest, destinationHost string, destinationRoutePort *api.PortSelector, layoutName string) istioapi.HTTPRoute {
	var port *istioapi.PortSelector
	if destinationRoutePort != nil {
		port = &istioapi.PortSelector{
			Number: destinationRoutePort.Number,
		}
	}
	return istioapi.HTTPRoute{
		Match: istioMatchList,
		Route: []istioapi.DestinationWeight{
			{
				Destination: istioapi.Destination{
					Host: destinationHost,
					Port: port,
				},
				Weight: 100,
			},
		},
		AppendHeaders: map[string]string{
			HeaderRouteName: layoutName,
		},
	}
}

func combineMatchesToIstioMatch(veMatchItem api.HTTPMatchRequest, targetingMatchItem api.HTTPMatchRequest) istioapi.HTTPMatchRequest {
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

	sourceLabels := combineStringMaps(veMatchItem.SourceLabels, targetingMatchItem.SourceLabels)
	gateways := combineStringSlicesUnique(veMatchItem.Gateways, targetingMatchItem.Gateways)

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
func combineStringMaps(map1 map[string]string, map2 map[string]string) map[string]string {
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
func combineStringSlicesUnique(slice1 []string, slice2 []string) []string {
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
		istioMap[keyVe] = *stringMatchToIstioStringMatch(&valueVe)
	}
	for keyTargeting, valueTargeting := range targetingMap {
		istioMap[keyTargeting] = *stringMatchToIstioStringMatch(&valueTargeting)
	}
	return istioMap
}

func selectStringMatchesToIstioStringMatch(veStringMatch *api.StringMatch, targetingStringMatch *api.StringMatch) *istioapi.StringMatch {
	if !isStringMatchEmpty(targetingStringMatch) {
		return stringMatchToIstioStringMatch(targetingStringMatch)
	} else if !isStringMatchEmpty(veStringMatch) {
		return stringMatchToIstioStringMatch(veStringMatch)
	}
	return nil
}

func isStringMatchEmpty(stringMatch *api.StringMatch) bool {
	return stringMatch == nil || (stringMatch.Exact == "" && stringMatch.Prefix == "" && stringMatch.Regex == "")
}

func stringMatchToIstioStringMatch(stringMatch *api.StringMatch) *istioapi.StringMatch {
	return &istioapi.StringMatch{
		Exact:  stringMatch.Exact,
		Prefix: stringMatch.Prefix,
		Regex:  stringMatch.Regex,
	}
}
