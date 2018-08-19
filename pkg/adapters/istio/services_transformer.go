package istio

import (
	"github.com/bevyx/remesh/pkg/adapters/istio/models"
	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
)

//TransformLayout is
func TransformLayout(layoutMap map[string]api.LayoutSpec) []models.TransformedService {
	transformedServices := make([]models.TransformedService, 0)
	for layoutName, layout := range layoutMap {
		for _, service := range layout.Services {
			subsetHash := computeHash(service.Labels)
			transformedService := findTransformedService(service.Host, &transformedServices)
			if transformedService == nil {
				transformedServices = append(transformedServices, models.TransformedService{
					Host:              service.Host,
					ServiceSubsetList: []models.ServiceSubset{makeServiceSubset(service.Labels, subsetHash, layoutName)},
				})
			} else {
				serviceSubset := findServiceSubset(subsetHash, &transformedService.ServiceSubsetList)
				if serviceSubset == nil {
					transformedService.ServiceSubsetList = append(transformedService.ServiceSubsetList, makeServiceSubset(service.Labels, subsetHash, layoutName))
				} else {
					serviceSubset.Layouts = append(serviceSubset.Layouts, layoutName)
				}
			}
		}
	}
	return transformedServices
}

func makeServiceSubset(labels map[string]string, subsetHash string, layoutName string) models.ServiceSubset {
	return models.ServiceSubset{
		Labels:     labels,
		SubsetHash: subsetHash,
		Layouts:    []string{layoutName},
	}
}

func findTransformedService(host string, transformedServices *[]models.TransformedService) *models.TransformedService {
	for i := range *transformedServices {
		transformedService := &(*transformedServices)[i]
		if transformedService.Host == host {
			return transformedService
		}
	}
	return nil
}

func findServiceSubset(subsetHash string, serviceSubsets *[]models.ServiceSubset) *models.ServiceSubset {
	for i := range *serviceSubsets {
		serviceSubset := &(*serviceSubsets)[i]
		if serviceSubset.SubsetHash == subsetHash {
			return serviceSubset
		}
	}
	return nil
}
