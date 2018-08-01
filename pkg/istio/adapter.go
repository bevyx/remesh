package istio

import (
	"log"

	"github.com/bevyx/remesh/pkg/models"
)

func Apply(entrypointFlows []models.EntrypointFlow, namespace string) {
	log.Printf("%#v", entrypointFlows)
}
