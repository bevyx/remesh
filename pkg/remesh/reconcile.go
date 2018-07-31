package remesh

import (
	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	"github.com/bevyx/remesh/pkg/models"
)

func Reconcile() (err error) {

	// entrypointFlows := combine(virtualEnvironmentList, targetingList, entrypointList)
	// istio.Apply(entrypointFlows, namespace)

	return nil
}

func combine(virtualEnvironmentList api.VirtualEnvironmentList, targetingList api.TargetingList, entrypointList api.EntrypointList) []models.EntrypointFlow {
	entrypointFlows := make([]models.EntrypointFlow, 0)
	for _, entrypoint := range entrypointList.Items {
		defaultVirtualEnvironment, ok := findVirtualEnvironment(entrypoint.Spec.DefaultVirtualEnvironment, virtualEnvironmentList.Items)
		if ok {
			targetings := getAllTargetingsByEntrypoint(entrypoint.ObjectMeta.Name, targetingList.Items)
			entrypointFlows = append(entrypointFlows, models.EntrypointFlow{
				Entrypoint:                entrypoint,
				DefaultVirtualEnvironment: defaultVirtualEnvironment,
				TargetingFlows:            combineTargetingToVirtualEnvironments(targetings, virtualEnvironmentList.Items),
				VirtualEnvironments:       virtualEnvironmentList.Items,
			})
		} else {
			// TODO: Notify that we are waiting for virtual env to be created
		}
	}
	return entrypointFlows
}

func combineTargetingToVirtualEnvironments(targetings []api.Targeting, virtualEnvironments []api.VirtualEnvironment) []models.TargetingFlow {
	targetingFlows := make([]models.TargetingFlow, 0)
	for _, targeting := range targetings {
		virtualEnvironment, ok := findVirtualEnvironment(targeting.Spec.VirtualEnvironment, virtualEnvironments)
		if ok {
			targetingFlows = append(targetingFlows, models.TargetingFlow{
				Targeting:          targeting,
				VirtualEnvironment: virtualEnvironment})
		}
	}
	return targetingFlows
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

func findVirtualEnvironment(name string, virtualEnvironments []api.VirtualEnvironment) (api.VirtualEnvironment, bool) {
	for _, virtualEnvironment := range virtualEnvironments {
		if virtualEnvironment.ObjectMeta.Name == name {
			return virtualEnvironment, true
		}
	}
	return api.VirtualEnvironment{}, false
}
