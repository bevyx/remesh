package istio

import (
	"fmt"
	"hash/fnv"

	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	"k8s.io/apimachinery/pkg/util/rand"
	hashutil "k8s.io/kubernetes/pkg/util/hash"
)

func computeHash(service map[string]string) string {
	serviceSubsetHasher := fnv.New32a()
	hashutil.DeepHashObject(serviceSubsetHasher, service)

	return rand.SafeEncodeString(fmt.Sprint(serviceSubsetHasher.Sum32()))
}

func getLayoutMapFromReleaseFlows(releaseFlows []api.ReleaseFlow) map[string]api.LayoutSpec {
	layoutMap := map[string]api.LayoutSpec{}
	for _, releaseFlow := range releaseFlows {
		if releaseFlow.Layout != nil {
			layoutMap[releaseFlow.LayoutName] = *releaseFlow.Layout
		}
	}
	return layoutMap
}
