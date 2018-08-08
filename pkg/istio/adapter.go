package istio

import (
	"log"

	istioapi "github.com/bevyx/istio-api-go/pkg/apis/networking/v1alpha3"
	istioapiclient "github.com/bevyx/istio-api-go/pkg/client/clientset/versioned"
	"github.com/bevyx/remesh/pkg/istio/resources"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/bevyx/remesh/pkg/models"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func Apply(entrypointFlow models.EntrypointFlow, namespace string) {
	log.Printf("applying %d entrypointFlows", len(entrypointFlow))
	desiredGateways, desiredVirtualServices := getDesiredResources(entrypointFlow)
	actualGateways, actualVirtualServices := getActualResources(namespace)

	//existing && !desired -> delete
	//desired && !existing -> create
	//existing && desired -> update
	//!existing && !desired -> noop

	//TODO: figure out how to operate on runtime.Object s, at least for ops like delete. how to use dynamic client?

	//Gateway
	desiredGatewaysObj := make([]runtime.Object, len(desiredGateways))
	for i, x := range desiredGateways {
		desiredGatewaysObj[i] = runtime.Object(&x)
	}
	actualGatewaysObj := make([]runtime.Object, len(actualGateways))
	for i, x := range actualGateways {
		actualGatewaysObj[i] = runtime.Object(&x)
	}
	gatewaysToDeleteObj := triageDelete(desiredGatewaysObj, actualGatewaysObj)
	gatewaysToDelete := make([]istioapi.Gateway, len(gatewaysToDeleteObj))
	for i, x := range gatewaysToDeleteObj {
		gatewaysToDelete[i] = *x.(*istioapi.Gateway)
	}
	deleteGateways(gatewaysToDelete)
	gatewaysToCreateObj := triageCreate(desiredGatewaysObj, actualGatewaysObj)
	gatewaysToCreate := make([]istioapi.Gateway, len(gatewaysToCreateObj))
	for i, x := range gatewaysToCreateObj {
		gatewaysToCreate[i] = *x.(*istioapi.Gateway)
	}
	createGateways(gatewaysToCreate)

	//VirtualService
	desiredVirtualServicesObj := make([]runtime.Object, len(desiredVirtualServices))
	for i, x := range desiredVirtualServices {
		desiredVirtualServicesObj[i] = runtime.Object(&x)
	}
	actualVirtualServicesObj := make([]runtime.Object, len(actualVirtualServices))
	for i, x := range actualVirtualServices {
		actualVirtualServicesObj[i] = runtime.Object(&x)
	}
	virtualServicesToDeleteObj := triageDelete(desiredVirtualServicesObj, actualVirtualServicesObj)
	virtualServicesToDelete := make([]istioapi.VirtualService, len(virtualServicesToDeleteObj))
	for i, x := range virtualServicesToDeleteObj {
		virtualServicesToDelete[i] = *x.(*istioapi.VirtualService)
	}
	deleteVirtualServices(virtualServicesToDelete)
	virtualServicesToCreateObj := triageCreate(desiredVirtualServicesObj, actualVirtualServicesObj)
	virtualServicesToCreate := make([]istioapi.VirtualService, len(virtualServicesToCreateObj))
	for i, x := range virtualServicesToCreateObj {
		virtualServicesToCreate[i] = *x.(*istioapi.VirtualService)
	}
	createVirtualServices(virtualServicesToCreate)
}

func deleteGateways(gateways []istioapi.Gateway) {
	cfg := config.GetConfigOrDie() //TODO: inject?
	istioclientset := istioapiclient.NewForConfigOrDie(cfg)
	for _, x := range gateways {
		err := istioclientset.NetworkingV1alpha3().Gateways(x.ObjectMeta.GetNamespace()).Delete(x.ObjectMeta.GetName(), nil)
		if err != nil {
			log.Printf("error deleting %s/%s :  %v", x.ObjectMeta.GetNamespace(), x.ObjectMeta.GetName(), err)
		}
	}
}

func createGateways(gateways []istioapi.Gateway) {
	cfg := config.GetConfigOrDie() //TODO: inject?
	istioclientset := istioapiclient.NewForConfigOrDie(cfg)
	for _, x := range gateways {
		_, err := istioclientset.NetworkingV1alpha3().Gateways(x.ObjectMeta.GetNamespace()).Create(&x)
		if err != nil {
			log.Printf("error creating %s/%s :  %v", x.ObjectMeta.GetNamespace(), x.ObjectMeta.GetName(), err)
		}
	}
}

func deleteVirtualServices(virtualServices []istioapi.VirtualService) {
	cfg := config.GetConfigOrDie() //TODO: inject?
	istioclientset := istioapiclient.NewForConfigOrDie(cfg)
	for _, x := range virtualServices {
		err := istioclientset.NetworkingV1alpha3().VirtualServices(x.ObjectMeta.GetNamespace()).Delete(x.ObjectMeta.GetName(), nil)
		if err != nil {
			log.Printf("error deleting %s/%s :  %v", x.ObjectMeta.GetNamespace(), x.ObjectMeta.GetName(), err)
		}
	}
}

func createVirtualServices(virtualServices []istioapi.VirtualService) {
	cfg := config.GetConfigOrDie() //TODO: inject?
	istioclientset := istioapiclient.NewForConfigOrDie(cfg)
	for _, x := range virtualServices {
		_, err := istioclientset.NetworkingV1alpha3().VirtualServices(x.ObjectMeta.GetNamespace()).Create(&x)
		if err != nil {
			log.Printf("error deleting %s/%s :  %v", x.ObjectMeta.GetNamespace(), x.ObjectMeta.GetName(), err)
		}
	}
}

func triageDelete(existing []runtime.Object, desired []runtime.Object) (res []runtime.Object) {
Outer:
	for _, e := range existing {
		for _, d := range desired {
			eAccessor, _ := meta.Accessor(e)
			dAccessor, _ := meta.Accessor(d)
			if eAccessor.GetName() == dAccessor.GetName() {
				break Outer
			}
		}
		res = append(res, e)
	}
	return res
}

func triageCreate(existing []runtime.Object, desired []runtime.Object) (res []runtime.Object) {
Outer:
	for _, d := range desired {
		for _, e := range existing {
			eAccessor, _ := meta.Accessor(e)
			dAccessor, _ := meta.Accessor(d)
			if eAccessor.GetName() == dAccessor.GetName() {
				break Outer
			}
		}
		res = append(res, d)
	}
	return res
}

func getDesiredResources(entrypointFlow models.EntrypointFlow) (gateways []istioapi.Gateway, virtualServices []istioapi.VirtualService) {
	gateways = make([]istioapi.Gateway, 0)
	virtualServices = make([]istioapi.VirtualService, 0)
	destinationRules = make([]istioapi.DestinationRule, 0)

	for _, entrypointFlow := range entrypointFlows {
		istioGateway, gatewayName := resources.MakeIstioGateway(entrypointFlow.Entrypoint, namespace)
		istioGateways = append(istioGateways, istioGateway)

		httpRoutes := MakeRouteForEntrypoint(entrypointFlow)
		gatewayVirtualService := resources.MakeIstioVirtualServiceForGateway(httpRoutes, namespace, gatewayName)
		istioVirtualServices = append(istioVirtualServices, gatewayVirtualService)

		transformedServices := TransformVirtualEnvironment(entrypointFlow.VirtualEnvironments)
		virtualServices := resources.MakeIstioVirtualServices(transformedServices, namespace, gatewayName)
		destinationRules := resources.MakeIstioDestinationRules(transformedServices, namespace, gatewayName)

		istioVirtualServices = append(istioVirtualServices, virtualServices...)
		istioDestinationRules = append(istioDestinationRules, destinationRules...)
	}

}

func getActualResources(namespace string) ([]istioapi.Gateway, []istioapi.VirtualService) {

	// TODO: check why kube-controller's (controller-runtime's) client didn't work

	// gatewayList := istioapi.GatewayList{}
	// virtualServiceList := istioapi.VirtualServiceList{}
	// destinationRuleList := istioapi.DestinationRuleList{}

	// err = client.List(context.TODO(), nil, &gatewayList)
	// if err != nil {
	// 	log.Printf("%v", err)
	// }
	// err = client.List(context.TODO(), nil, &virtualServiceList)
	// if err != nil {
	// 	log.Printf("%v", err)
	// }
	// err = client.List(context.TODO(), nil, &destinationRuleList)
	// if err != nil {
	// 	log.Printf("%v", err)
	// }

	cfg := config.GetConfigOrDie() //TODO: inject?
	istioclientset := istioapiclient.NewForConfigOrDie(cfg)
	gatewayList, err := istioclientset.NetworkingV1alpha3().Gateways(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	virtualServiceList, err := istioclientset.NetworkingV1alpha3().VirtualServices(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return gatewayList.Items, virtualServiceList.Items
}
