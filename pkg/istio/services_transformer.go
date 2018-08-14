package istio

import (
	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	"github.com/bevyx/remesh/pkg/istio/models"
)

//TransformLayout is
func TransformLayout(layouts []api.Layout) []models.TransformedService {
	transformedServices := make([]models.TransformedService, 0)
	for _, layout := range layouts {
		for _, service := range layout.Spec.Services {
			subsetHash := computeHash(service.Labels)
			transformedService := findTransformedService(service.Host, &transformedServices)
			if transformedService == nil {
				transformedServices = append(transformedServices, models.TransformedService{
					Host:              service.Host,
					ServiceSubsetList: []models.ServiceSubset{makeServiceSubset(service.Labels, subsetHash, layout.Name)},
				})
			} else {
				serviceSubset := findServiceSubset(subsetHash, &transformedService.ServiceSubsetList)
				if serviceSubset == nil {
					transformedService.ServiceSubsetList = append(transformedService.ServiceSubsetList, makeServiceSubset(service.Labels, subsetHash, layout.Name))
				} else {
					serviceSubset.Layouts = append(serviceSubset.Layouts, layout.Name)
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
