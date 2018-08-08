package istio

import (
	"log"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/bevyx/remesh/pkg/models"
	knativeistio "github.com/knative/serving/pkg/apis/istio/v1alpha3"
	knativeistioclient "github.com/knative/serving/pkg/client/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func Apply(entrypointFlows []models.EntrypointFlow, namespace string) {
	log.Printf("applying %d entrypointFlows", len(entrypointFlows))
	desiredGateways, desiredVirtualServices := getDesiredResources(entrypointFlows)
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
	gatewaysToDelete := make([]knativeistio.Gateway, len(gatewaysToDeleteObj))
	for i, x := range gatewaysToDeleteObj {
		gatewaysToDelete[i] = *x.(*knativeistio.Gateway)
	}
	deleteGateways(gatewaysToDelete)
	gatewaysToCreateObj := triageCreate(desiredGatewaysObj, actualGatewaysObj)
	gatewaysToCreate := make([]knativeistio.Gateway, len(gatewaysToCreateObj))
	for i, x := range gatewaysToCreateObj {
		gatewaysToCreate[i] = *x.(*knativeistio.Gateway)
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
	virtualServicesToDelete := make([]knativeistio.VirtualService, len(virtualServicesToDeleteObj))
	for i, x := range virtualServicesToDeleteObj {
		virtualServicesToDelete[i] = *x.(*knativeistio.VirtualService)
	}
	deleteVirtualServices(virtualServicesToDelete)
	virtualServicesToCreateObj := triageCreate(desiredVirtualServicesObj, actualVirtualServicesObj)
	virtualServicesToCreate := make([]knativeistio.VirtualService, len(virtualServicesToCreateObj))
	for i, x := range virtualServicesToCreateObj {
		virtualServicesToCreate[i] = *x.(*knativeistio.VirtualService)
	}
	createVirtualServices(virtualServicesToCreate)
}

func deleteGateways(gateways []knativeistio.Gateway) {
	cfg := config.GetConfigOrDie() //TODO: inject?
	knativeistioclientset := knativeistioclient.NewForConfigOrDie(cfg)
	for _, x := range gateways {
		err := knativeistioclientset.NetworkingV1alpha3().Gateways(x.ObjectMeta.GetNamespace()).Delete(x.ObjectMeta.GetName(), nil)
		if err != nil {
			log.Printf("error deleting %s/%s :  %v", x.ObjectMeta.GetNamespace(), x.ObjectMeta.GetName(), err)
		}
	}
}

func createGateways(gateways []knativeistio.Gateway) {
	cfg := config.GetConfigOrDie() //TODO: inject?
	knativeistioclientset := knativeistioclient.NewForConfigOrDie(cfg)
	for _, x := range gateways {
		_, err := knativeistioclientset.NetworkingV1alpha3().Gateways(x.ObjectMeta.GetNamespace()).Create(&x)
		if err != nil {
			log.Printf("error creating %s/%s :  %v", x.ObjectMeta.GetNamespace(), x.ObjectMeta.GetName(), err)
		}
	}
}

func deleteVirtualServices(virtualServices []knativeistio.VirtualService) {
	cfg := config.GetConfigOrDie() //TODO: inject?
	knativeistioclientset := knativeistioclient.NewForConfigOrDie(cfg)
	for _, x := range virtualServices {
		err := knativeistioclientset.NetworkingV1alpha3().VirtualServices(x.ObjectMeta.GetNamespace()).Delete(x.ObjectMeta.GetName(), nil)
		if err != nil {
			log.Printf("error deleting %s/%s :  %v", x.ObjectMeta.GetNamespace(), x.ObjectMeta.GetName(), err)
		}
	}
}

func createVirtualServices(virtualServices []knativeistio.VirtualService) {
	cfg := config.GetConfigOrDie() //TODO: inject?
	knativeistioclientset := knativeistioclient.NewForConfigOrDie(cfg)
	for _, x := range virtualServices {
		_, err := knativeistioclientset.NetworkingV1alpha3().VirtualServices(x.ObjectMeta.GetNamespace()).Create(&x)
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

func getDesiredResources(entrypointFlows []models.EntrypointFlow) (gateways []knativeistio.Gateway, virtualServices []knativeistio.VirtualService) {
	gateways = make([]knativeistio.Gateway, 0)
	virtualServices = make([]knativeistio.VirtualService, 0)
	return
}

func getActualResources(namespace string) ([]knativeistio.Gateway, []knativeistio.VirtualService) {

	// TODO: check why kube-controller's (controller-runtime's) client didn't work

	// gatewayList := knativeistio.GatewayList{}
	// virtualServiceList := knativeistio.VirtualServiceList{}
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
	knativeistioclientset := knativeistioclient.NewForConfigOrDie(cfg)
	gatewayList, err := knativeistioclientset.NetworkingV1alpha3().Gateways(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	virtualServiceList, err := knativeistioclientset.NetworkingV1alpha3().VirtualServices(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return gatewayList.Items, virtualServiceList.Items
}
