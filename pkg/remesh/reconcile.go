package remesh

import (
	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	"github.com/bevyx/remesh/pkg/models"
)

func Combine(layoutList api.LayoutList, releaseList api.ReleaseList, segmentList api.SegmentList, entrypointList api.EntrypointList) []models.EntrypointFlow {
	entrypointFlows := make([]models.EntrypointFlow, 0)
	for _, entrypoint := range entrypointList.Items {
		releases := findReleasesByEntrypoint(entrypoint.Name, releaseList.Items)
		if len(releases) > 0 {
			releaseFlows := combineReleasesToSegmentsAndLayouts(releases, segmentList.Items, layoutList.Items)
			if len(releaseFlows) > 0 {
				layoutSet := getLayoutSetOfEntrypointFlow(releaseFlows)
				entrypointFlows = append(entrypointFlows, models.EntrypointFlow{
					Entrypoint:   entrypoint,
					ReleaseFlows: releaseFlows,
					Layouts:      layoutSet,
				})
			} else {
				// TODO: notify entrypoint wans't created at all. Waiting for some releases to be ready
			}

		} else {
			// TODO: notify entrypoint wans't created at all. Waiting for some releases to be created
		}

	}
	return entrypointFlows
}

func combineReleasesToSegmentsAndLayouts(releases []api.Release, segments []api.Segment, layouts []api.Layout) (releaseFlows []models.ReleaseFlow) {
	releaseFlows = make([]models.ReleaseFlow, 0)
	for _, release := range releases {
		layout := findLayout(release.Spec.Layout, layouts)
		if layout != nil {
			releaseFlow := combineOneReleaseToSegmentsAndALayout(release, segments, *layout)
			if releaseFlow != nil {
				releaseFlows = append(releaseFlows, *releaseFlow)
			}
		} else {
			// TODO: notify release wans't created at all. Waiting for a layout to be created
		}
	}
	return
}

func combineOneReleaseToSegmentsAndALayout(release api.Release, segments []api.Segment, layout api.Layout) (releaseFlow *models.ReleaseFlow) {
	if release.Spec.Targeting == nil {
		releaseFlow = &models.ReleaseFlow{
			Release:  release,
			Segments: []api.Segment{},
			Layout:   layout}
	} else {
		existSegments, dontExistSegments := getSegmentsOfRelease(release, segments)
		if len(existSegments) > 0 {
			releaseFlow = &models.ReleaseFlow{
				Release:  release,
				Segments: existSegments,
				Layout:   layout}
		} else {
			// TODO: notify release wans't created at all. Waiting for some segments to be created
		}
		if len(dontExistSegments) > 0 {
			// TODO: notify some segments don't exist
		}
	}
	return
}

func getLayoutSetOfEntrypointFlow(releaseFlows []models.ReleaseFlow) (layoutSet []api.Layout) {
	layoutMap := map[string]api.Layout{}
	for _, releaseFlow := range releaseFlows {
		layoutMap[releaseFlow.Layout.Name] = releaseFlow.Layout
	}
	layoutSet = make([]api.Layout, 0)
	for _, value := range layoutMap {
		layoutSet = append(layoutSet, value)
	}
	return
}

func getSegmentsOfRelease(release api.Release, segments []api.Segment) (found []api.Segment, dontExistSegments []string) {
	found = make([]api.Segment, 0)
	dontExistSegments = make([]string, 0)
	if release.Spec.Targeting != nil {
		for _, segmentName := range release.Spec.Targeting.Segments {
			segment := findSegment(segmentName, segments)
			if segment != nil {
				found = append(found, *segment)
			} else {
				dontExistSegments = append(dontExistSegments, segmentName)
			}
		}
	}
	return
}

func findLayout(name string, layouts []api.Layout) *api.Layout {
	for _, layout := range layouts {
		if layout.Name == name {
			return &layout
		}
	}
	return nil
}

func findSegment(name string, segments []api.Segment) *api.Segment {
	for _, segment := range segments {
		if segment.Name == name {
			return &segment
		}
	}
	return nil
}

func findReleasesByEntrypoint(entrypoint string, releases []api.Release) []api.Release {
	found := make([]api.Release, 0)
	for _, release := range releases {
		if release.Spec.Entrypoint == entrypoint {
			found = append(found, release)
		}
	}
	return found
}

func findReleasesWithNoTargeting(releases []api.Release) []api.Release {
	found := make([]api.Release, 0)
	for _, release := range releases {
		if release.Spec.Targeting == nil {
			found = append(found, release)
		}
	}
	return found
}
