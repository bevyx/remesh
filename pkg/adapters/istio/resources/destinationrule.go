package resources

import (
	istioapi "github.com/bevyx/istio-api-go/pkg/apis/networking/v1alpha3"
	istiomodels "github.com/bevyx/remesh/pkg/adapters/istio/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//MakeIstioDestinationRules is
func MakeIstioDestinationRules(transformedServices []istiomodels.TransformedService, namespace string) []istioapi.DestinationRule {
	destinationRules := make([]istioapi.DestinationRule, 0)
	for _, transformedService := range transformedServices {
		destinationRules = append(destinationRules, makeDestinationRule(transformedService, namespace))
	}
	return destinationRules
}

func makeDestinationRule(transformedService istiomodels.TransformedService, namespace string) istioapi.DestinationRule {
	return istioapi.DestinationRule{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DestinationRule",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Prefix + transformedService.Host + DrSuffix,
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
		Spec: makeDestinationRuleSpec(transformedService),
	}
}

func makeDestinationRuleSpec(transformedService istiomodels.TransformedService) istioapi.DestinationRuleSpec {
	subsets := make([]istioapi.Subset, 0)
	for _, subsetService := range transformedService.ServiceSubsetList {
		subsets = append(subsets, istioapi.Subset{
			Name:   GetSubsetName(transformedService.Host, subsetService.SubsetHash),
			Labels: subsetService.Labels,
		})
	}
	return istioapi.DestinationRuleSpec{
		Host:    transformedService.Host,
		Subsets: subsets,
	}
}