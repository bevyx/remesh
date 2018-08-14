package remesh

import (
	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	"github.com/bevyx/remesh/pkg/models"
)

func Reconcile() (err error) {

	// entrypointFlows := combine(layoutList, targetingList, entrypointList)
	// istio.Apply(entrypointFlows, namespace)

	return nil
}

func Combine(layoutList api.LayoutList, targetingList api.TargetingList, entrypointList api.EntrypointList) []models.EntrypointFlow {
	entrypointFlows := make([]models.EntrypointFlow, 0)
	for _, entrypoint := range entrypointList.Items {
		defaultLayout, ok := findLayout(entrypoint.Spec.DefaultLayout, layoutList.Items)
		if ok {
			targetings := getAllTargetingsByEntrypoint(entrypoint.ObjectMeta.Name, targetingList.Items)
			targetingFlows, layoutSet := combineTargetingToLayouts(targetings, layoutList.Items)
			_, isDefaultInSet := findLayout(defaultLayout.Name, layoutSet)
			if !isDefaultInSet {
				layoutSet = append(layoutSet, defaultLayout)
			}
			entrypointFlows = append(entrypointFlows, models.EntrypointFlow{
				Entrypoint:     entrypoint,
				DefaultLayout:  defaultLayout,
				TargetingFlows: targetingFlows,
				Layouts:        layoutSet,
			})
		} else {
			// TODO: Notify that we are waiting for virtual env to be created
		}
	}
	return entrypointFlows
}

func combineTargetingToLayouts(targetings []api.Targeting, layouts []api.Layout) ([]models.TargetingFlow, []api.Layout) {
	targetingFlows := make([]models.TargetingFlow, 0)
	layoutMap := map[string]api.Layout{}
	for _, targeting := range targetings {
		layout, ok := findLayout(targeting.Spec.Layout, layouts)
		layoutMap[layout.Name] = layout
		if ok {
			targetingFlows = append(targetingFlows, models.TargetingFlow{
				Targeting: targeting,
				Layout:    layout})
		}
	}
	layoutSet := make([]api.Layout, 0)
	for _, value := range layoutMap {
		layoutSet = append(layoutSet, value)
	}
	return targetingFlows, layoutSet
}

func getAllTargetingsByEntrypoint(entrypointName string, targetings []api.Targeting) []api.Targeting {
	targetingsOfEntrypoint := make([]api.Targeting, 0)
	for _, targeting := range targetings {
		if targeting.Spec.Entrypoint == entrypointName {
			targetingsOfEntrypoint = append(targetingsOfEntrypoint, targeting)
		}
	}
	return targetingsOfEntrypoint
}

func findLayout(name string, layouts []api.Layout) (api.Layout, bool) {
	for _, layout := range layouts {
		if layout.ObjectMeta.Name == name {
			return layout, true
		}
	}
	return api.Layout{}, false
}
