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

func Apply(entrypointFlows []models.EntrypointFlow, namespace string) {
	log.Printf("applying %d entrypointFlows", len(entrypointFlows))
	desiredGateways, desiredVirtualServices, desiredDestinationRules := getDesiredResources(entrypointFlows, namespace)
	log.Printf("desired: gateways %d, virtualservices: %d, destinationrules %d", len(desiredGateways), len(desiredVirtualServices), len(desiredDestinationRules))
	actualGateways, actualVirtualServices, actualDestinationRules := getActualResources(namespace)
	log.Printf("actual: gateways %d, virtualservices: %d, destinationrules %d", len(actualGateways), len(actualVirtualServices), len(actualDestinationRules))

	//existing && !desired -> delete
	//desired && !existing -> create
	//existing && desired -> update
	//!existing && !desired -> noop

	//TODO: figure out how to operate on runtime.Object s, at least for ops like delete. how to use dynamic client?

	//Gateway
	desiredGatewaysObj := make([]runtime.Object, len(desiredGateways))
	for i, x := range desiredGateways {
		desiredGatewaysObj[i] = x.DeepCopyObject()
	}
	actualGatewaysObj := make([]runtime.Object, len(actualGateways))
	for i, x := range actualGateways {
		actualGatewaysObj[i] = x.DeepCopyObject()
	}
	gatewaysToDeleteObj := triageDelete(actualGatewaysObj, desiredGatewaysObj)
	gatewaysToDelete := make([]istioapi.Gateway, len(gatewaysToDeleteObj))
	for i, x := range gatewaysToDeleteObj {
		gatewaysToDelete[i] = *x.(*istioapi.Gateway)
	}
	log.Printf("deleting %d gateways", len(gatewaysToDelete))
	deleteGateways(gatewaysToDelete)
	gatewaysToCreateObj := triageCreate(actualGatewaysObj, desiredGatewaysObj)
	gatewaysToCreate := make([]istioapi.Gateway, len(gatewaysToCreateObj))
	for i, x := range gatewaysToCreateObj {
		gatewaysToCreate[i] = *x.(*istioapi.Gateway)
	}
	log.Printf("creating %d gateways", len(gatewaysToCreate))
	createGateways(gatewaysToCreate)

	//VirtualService
	desiredVirtualServicesObj := make([]runtime.Object, len(desiredVirtualServices))
	for i, x := range desiredVirtualServices {
		desiredVirtualServicesObj[i] = x.DeepCopyObject()
	}
	actualVirtualServicesObj := make([]runtime.Object, len(actualVirtualServices))
	for i, x := range actualVirtualServices {
		actualVirtualServicesObj[i] = x.DeepCopyObject()
	}
	virtualServicesToDeleteObj := triageDelete(actualVirtualServicesObj, desiredVirtualServicesObj)
	virtualServicesToDelete := make([]istioapi.VirtualService, len(virtualServicesToDeleteObj))
	for i, x := range virtualServicesToDeleteObj {
		virtualServicesToDelete[i] = *x.(*istioapi.VirtualService)
	}
	log.Printf("deleting %d virtualservices", len(virtualServicesToDelete))
	deleteVirtualServices(virtualServicesToDelete)
	virtualServicesToCreateObj := triageCreate(actualVirtualServicesObj, desiredVirtualServicesObj)
	virtualServicesToCreate := make([]istioapi.VirtualService, len(virtualServicesToCreateObj))
	for i, x := range virtualServicesToCreateObj {
		virtualServicesToCreate[i] = *x.(*istioapi.VirtualService)
	}
	log.Printf("creating %d virtualservices", len(virtualServicesToCreate))
	createVirtualServices(virtualServicesToCreate)

	//DestinationRule
	desiredDestinationRulesObj := make([]runtime.Object, len(desiredDestinationRules))
	for i, x := range desiredDestinationRules {
		desiredDestinationRulesObj[i] = x.DeepCopyObject()
	}
	actualDestinationRulesObj := make([]runtime.Object, len(actualDestinationRules))
	for i, x := range actualDestinationRules {
		actualDestinationRulesObj[i] = x.DeepCopyObject()
	}
	destinationRulesToDeleteObj := triageDelete(actualDestinationRulesObj, desiredDestinationRulesObj)
	destinationRulesToDelete := make([]istioapi.DestinationRule, len(destinationRulesToDeleteObj))
	for i, x := range destinationRulesToDeleteObj {
		destinationRulesToDelete[i] = *x.(*istioapi.DestinationRule)
	}
	log.Printf("deleting %d destinationrules", len(destinationRulesToDelete))
	deleteDestinationRules(destinationRulesToDelete)
	destinationRulesToCreateObj := triageCreate(actualDestinationRulesObj, desiredDestinationRulesObj)
	destinationRulesToCreate := make([]istioapi.DestinationRule, len(destinationRulesToCreateObj))
	for i, x := range destinationRulesToCreateObj {
		destinationRulesToCreate[i] = *x.(*istioapi.DestinationRule)
	}
	log.Printf("creating %d destinationrules", len(destinationRulesToCreate))
	createDestinationRules(destinationRulesToCreate)
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
			log.Printf("error create %s/%s :  %v", x.ObjectMeta.GetNamespace(), x.ObjectMeta.GetName(), err)
		}
	}
}

func deleteDestinationRules(destinationRules []istioapi.DestinationRule) {
	cfg := config.GetConfigOrDie() //TODO: inject?
	istioclientset := istioapiclient.NewForConfigOrDie(cfg)
	for _, x := range destinationRules {
		err := istioclientset.NetworkingV1alpha3().DestinationRules(x.ObjectMeta.GetNamespace()).Delete(x.ObjectMeta.GetName(), nil)
		if err != nil {
			log.Printf("error deleting %s/%s :  %v", x.ObjectMeta.GetNamespace(), x.ObjectMeta.GetName(), err)
		}
	}
}

func createDestinationRules(destinationRules []istioapi.DestinationRule) {
	cfg := config.GetConfigOrDie() //TODO: inject?
	istioclientset := istioapiclient.NewForConfigOrDie(cfg)
	for _, x := range destinationRules {
		_, err := istioclientset.NetworkingV1alpha3().DestinationRules(x.ObjectMeta.GetNamespace()).Create(&x)
		if err != nil {
			log.Printf("error create %s/%s :  %v", x.ObjectMeta.GetNamespace(), x.ObjectMeta.GetName(), err)
		}
	}
}

func triageDelete(existing []runtime.Object, desired []runtime.Object) (res []runtime.Object) {
	for _, e := range existing {
		if !objectExists(e, desired) {
			res = append(res, e)
		}
	}
	return res
}

func triageCreate(existing []runtime.Object, desired []runtime.Object) (res []runtime.Object) {
	for _, d := range desired {
		if !objectExists(d, existing) {
			res = append(res, d)
		}
	}
	return res
}

func objectExists(obj runtime.Object, list []runtime.Object) bool {
	for _, i := range list {
		iAccessor, _ := meta.Accessor(i)
		objAccessor, _ := meta.Accessor(obj)
		if iAccessor.GetName() == objAccessor.GetName() {
			return true
		}
	}
	return false
}

func getDesiredResources(entrypointFlows []models.EntrypointFlow, namespace string) (gateways []istioapi.Gateway, virtualServices []istioapi.VirtualService, destinationRules []istioapi.DestinationRule) {
	gateways = make([]istioapi.Gateway, 0)
	virtualServices = make([]istioapi.VirtualService, 0)
	destinationRules = make([]istioapi.DestinationRule, 0)

	for _, entrypointFlow := range entrypointFlows {
		gateway, gatewayName := resources.MakeIstioGateway(entrypointFlow.Entrypoint, namespace)
		gateways = append(gateways, gateway)

		httpRoutes := MakeRouteForEntrypoint(entrypointFlow)
		gatewayVirtualService := resources.MakeIstioVirtualServiceForGateway(httpRoutes, namespace, gatewayName)
		virtualServices = append(virtualServices, gatewayVirtualService)

		transformedServices := TransformVirtualEnvironment(entrypointFlow.VirtualEnvironments)
		transformedVirtualServices := resources.MakeIstioVirtualServices(transformedServices, namespace, gatewayName)
		transformedDestinationRules := resources.MakeIstioDestinationRules(transformedServices, namespace, gatewayName)

		virtualServices = append(virtualServices, transformedVirtualServices...)
		destinationRules = append(destinationRules, transformedDestinationRules...)
	}
	return
}

func getActualResources(namespace string) ([]istioapi.Gateway, []istioapi.VirtualService, []istioapi.DestinationRule) {

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
	destinationRuleList, err := istioclientset.NetworkingV1alpha3().DestinationRules(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return gatewayList.Items, virtualServiceList.Items, destinationRuleList.Items
}
