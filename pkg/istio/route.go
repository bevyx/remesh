package istio

// import "github.com/bevyx/remesh/pkg/models"

// func MakeRouteForEntrypoint(entrypointFlow models.EntrypointFlow) []TransformedService {
// 	transformedServices := make([]TransformedService, 0)
// 	for _, virtualEnvironment := range virtualEnvironments {
// 		for _, service := range virtualEnvironment.Spec.Services {
// 			subsetHash := computeHash(&service.Labels)
// 			transformedService := findTransformedService(service.Host, &transformedServices)
// 			if transformedService == nil {
// 				transformedServices = append(transformedServices, TransformedService{
// 					Host:              service.Host,
// 					ServiceSubsetList: []ServiceSubset{*makeServiceSubset(service.Labels, subsetHash, virtualEnvironment.Name)},
// 				})
// 			} else {
// 				serviceSubset := findServiceSubset(subsetHash, &transformedService.ServiceSubsetList)
// 				if serviceSubset == nil {
// 					transformedService.ServiceSubsetList = append(transformedService.ServiceSubsetList, *makeServiceSubset(service.Labels, subsetHash, virtualEnvironment.Name))
// 				} else {
// 					serviceSubset.VirtualEnvironments = append(serviceSubset.VirtualEnvironments, virtualEnvironment.Name)
// 				}
// 			}
// 		}
// 	}
// 	return transformedServices
// }

// func makeServiceSubset(labels map[string]string, subsetHash string, virtualEnvironmentName string) *ServiceSubset {
// 	return &ServiceSubset{
// 		Labels:              labels,
// 		SubsetHash:          subsetHash,
// 		VirtualEnvironments: []string{virtualEnvironmentName},
// 	}
// }

// func findTransformedService(host string, transformedServices *[]TransformedService) *TransformedService {
// 	for i := range *transformedServices {
// 		transformedService := &(*transformedServices)[i]
// 		if transformedService.Host == host {
// 			return transformedService
// 		}
// 	}
// 	return nil
// }

// func findServiceSubset(subsetHash string, serviceSubsets *[]ServiceSubset) *ServiceSubset {
// 	for i := range *serviceSubsets {
// 		serviceSubset := &(*serviceSubsets)[i]
// 		if serviceSubset.SubsetHash == subsetHash {
// 			return serviceSubset
// 		}
// 	}
// 	return nil
// }
