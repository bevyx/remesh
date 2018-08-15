package models

import (
	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
)

type ReleaseFlow struct {
	Release  api.Release
	Segments []api.Segment
	Layout   api.Layout
}

type EntrypointFlow struct {
	Entrypoint   api.Entrypoint
	ReleaseFlows []ReleaseFlow
	Layouts      []api.Layout
}

// ByPriority implements sort.Interface for []ReleaseFlow based on
// the Release Priority field.
type ByPriority []ReleaseFlow

func (a ByPriority) Len() int      { return len(a) }
func (a ByPriority) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool {
	if a[i].Release.Spec.Targeting == nil {
		return true
	}
	if a[j].Release.Spec.Targeting == nil {
		return false
	}
	return a[i].Release.Spec.Targeting.Priority > a[j].Release.Spec.Targeting.Priority
}
