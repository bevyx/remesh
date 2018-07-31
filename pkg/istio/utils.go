package istio

import (
	"fmt"
	"hash/fnv"

	"k8s.io/apimachinery/pkg/util/rand"
	hashutil "k8s.io/kubernetes/pkg/util/hash"
)

func computeHash(service *map[string]string) string {
	serviceSubsetHasher := fnv.New32a()
	hashutil.DeepHashObject(serviceSubsetHasher, *service)

	return rand.SafeEncodeString(fmt.Sprint(serviceSubsetHasher.Sum32()))
}
