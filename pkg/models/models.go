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
