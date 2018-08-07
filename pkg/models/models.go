package models

import (
	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
)

type TargetingFlow struct {
	Targeting          api.Targeting
	VirtualEnvironment api.VirtualEnvironment
}

type EntrypointFlow struct {
	Entrypoint                api.Entrypoint
	DefaultVirtualEnvironment api.VirtualEnvironment
	TargetingFlows            []TargetingFlow
	VirtualEnvironments       []api.VirtualEnvironment
}

// ByPriority implements sort.Interface for []TargetingFlow based on
// the Targeting Priority field.
type ByPriority []TargetingFlow

func (a ByPriority) Len() int      { return len(a) }
func (a ByPriority) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool {
	return a[i].Targeting.Spec.Priority < a[j].Targeting.Spec.Priority
}
