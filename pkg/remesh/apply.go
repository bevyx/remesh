package remesh

import remeshv1alpha1 "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"

//Applier is an interface for a target mesh. It gets the desired state of the VirtualApp, and is responsible to
//do all the actions necessary to apply this desired state to the target service service mesh in the given namespace (including translating
//the VirtualApp resources to service mesh native resources)
//TODO: does namespace belong in the interface or is it istio implementation detail? vapp is a kube resource so it has a ns, maybe pass it along?
type Applier interface {
	Apply(virtualApps []remeshv1alpha1.VirtualApp, namespace string)
}
