package istio

import (
	"context"
	"log"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	istioapi "github.com/bevyx/istio-api-go/pkg/apis/networking/v1alpha3"
	"github.com/bevyx/remesh/pkg/istio/resources"
	"github.com/bevyx/remesh/pkg/remesh"
	"k8s.io/apimachinery/pkg/api/meta"

	remeshv1alpha1 "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

//IstioApplier is an implementation of the remesh.Applier interface that can apply a virtualApp configuration to Istio service mesh
type IstioApplier struct {
	istioClient client.Client
}

//NewIstioApplier creates new IstioApplier
func NewIstioApplier() remesh.Applier {
	clientConfig := config.GetConfigOrDie()
	istioScheme := runtime.NewScheme()
	istioapi.AddToScheme(istioScheme)
	istioClient, err := client.New(clientConfig, client.Options{Scheme: istioScheme})
	if err != nil {
		log.Printf("%v", err)
	}
	return &IstioApplier{istioClient: istioClient}
}

//Apply is implementing the remesh.Applier interface
func (a *IstioApplier) Apply(virtualApps []remeshv1alpha1.VirtualApp, namespace string) {
	log.Printf("applying %d virtualApps", len(virtualApps))
	desiredGateways, desiredVirtualServices, desiredDestinationRules := a.getDesiredResources(virtualApps, namespace)
	log.Printf("desired: gateways %d, virtualservices: %d, destinationrules %d", len(desiredGateways), len(desiredVirtualServices), len(desiredDestinationRules))
	actualGateways, actualVirtualServices, actualDestinationRules := a.getActualResources(namespace)
	log.Printf("actual: gateways %d, virtualservices: %d, destinationrules %d", len(actualGateways), len(actualVirtualServices), len(actualDestinationRules))

	//existing && !desired -> delete
	//desired && !existing -> create
	//existing && desired -> update
	//!existing && !desired -> noop

	//TODO: figure out how to make this more generic (runtime.Object), less code dup

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
	a.deleteGateways(gatewaysToDelete)

	gatewaysToCreateObj := triageCreate(actualGatewaysObj, desiredGatewaysObj)
	gatewaysToCreate := make([]istioapi.Gateway, len(gatewaysToCreateObj))
	for i, x := range gatewaysToCreateObj {
		gatewaysToCreate[i] = *x.(*istioapi.Gateway)
	}
	log.Printf("creating %d gateways", len(gatewaysToCreate))
	a.createGateways(gatewaysToCreate)

	// gatewaysToUpdateObj := triageUpdate(actualGatewaysObj, desiredGatewaysObj)
	// gatewaysToUpdate := make([]istioapi.Gateway, len(gatewaysToUpdateObj))
	// for i, x := range gatewaysToUpdateObj {
	// 	gatewaysToUpdate[i] = *x.(*istioapi.Gateway)
	// }
	// log.Printf("updating %d destinationrules", len(gatewaysToUpdate))
	// a.updateGateways(gatewaysToUpdate)

	gatewaysToUpdateObj := triageUpdate(actualGatewaysObj, desiredGatewaysObj)
	gatewaysToUpdate := make([]struct {
		actual  istioapi.Gateway
		desired istioapi.Gateway
	}, len(gatewaysToUpdateObj))
	for i, x := range gatewaysToUpdateObj {
		tmp := x.(struct {
			actual  runtime.Object
			desired runtime.Object
		})
		gatewaysToUpdate[i] = struct {
			actual  istioapi.Gateway
			desired istioapi.Gateway
		}{
			*tmp.actual.(*istioapi.Gateway),
			*tmp.desired.(*istioapi.Gateway),
		}
	}
	log.Printf("updating %d gateways", len(gatewaysToUpdate))
	a.updateGateways(gatewaysToUpdate)

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
	a.deleteVirtualServices(virtualServicesToDelete)

	virtualServicesToCreateObj := triageCreate(actualVirtualServicesObj, desiredVirtualServicesObj)
	virtualServicesToCreate := make([]istioapi.VirtualService, len(virtualServicesToCreateObj))
	for i, x := range virtualServicesToCreateObj {
		virtualServicesToCreate[i] = *x.(*istioapi.VirtualService)
	}
	log.Printf("creating %d virtualservices", len(virtualServicesToCreate))
	a.createVirtualServices(virtualServicesToCreate)

	virtualServicesToUpdateObj := triageUpdate(actualVirtualServicesObj, desiredVirtualServicesObj)
	virtualServicesToUpdate := make([]struct {
		actual  istioapi.VirtualService
		desired istioapi.VirtualService
	}, len(virtualServicesToUpdateObj))
	for i, x := range virtualServicesToUpdateObj {
		tmp := x.(struct {
			actual  runtime.Object
			desired runtime.Object
		})
		virtualServicesToUpdate[i] = struct {
			actual  istioapi.VirtualService
			desired istioapi.VirtualService
		}{
			*tmp.actual.(*istioapi.VirtualService),
			*tmp.desired.(*istioapi.VirtualService),
		}
	}
	log.Printf("updating %d virtualServices", len(virtualServicesToUpdate))
	a.updateVirtualServices(virtualServicesToUpdate)

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
	a.deleteDestinationRules(destinationRulesToDelete)

	destinationRulesToCreateObj := triageCreate(actualDestinationRulesObj, desiredDestinationRulesObj)
	destinationRulesToCreate := make([]istioapi.DestinationRule, len(destinationRulesToCreateObj))
	for i, x := range destinationRulesToCreateObj {
		destinationRulesToCreate[i] = *x.(*istioapi.DestinationRule)
	}
	log.Printf("creating %d destinationrules", len(destinationRulesToCreate))
	a.createDestinationRules(destinationRulesToCreate)

	destinationRulesToUpdateObj := triageUpdate(actualDestinationRulesObj, desiredDestinationRulesObj)
	destinationRulesToUpdate := make([]struct {
		actual  istioapi.DestinationRule
		desired istioapi.DestinationRule
	}, len(destinationRulesToUpdateObj))
	for i, x := range destinationRulesToUpdateObj {
		tmp := x.(struct {
			actual  runtime.Object
			desired runtime.Object
		})
		destinationRulesToUpdate[i] = struct {
			actual  istioapi.DestinationRule
			desired istioapi.DestinationRule
		}{
			*tmp.actual.(*istioapi.DestinationRule),
			*tmp.desired.(*istioapi.DestinationRule),
		}
	}
	log.Printf("updating %d destinationRules", len(destinationRulesToUpdate))
	a.updateDestinationRules(destinationRulesToUpdate)
}

func (a *IstioApplier) deleteGateways(gateways []istioapi.Gateway) {
	// cfg := config.GetConfigOrDie() //TODO: inject?
	// istioclientset := istioapiclient.NewForConfigOrDie(cfg)
	// for _, x := range gateways {
	// 	err := istioclientset.NetworkingV1alpha3().Gateways(x.ObjectMeta.GetNamespace()).Delete(x.ObjectMeta.GetName(), nil)
	// 	if err != nil {
	// 		log.Printf("error delete %s/%s :  %v", x.ObjectMeta.GetNamespace(), x.ObjectMeta.GetName(), err)
	// 	}
	// }

	for _, x := range gateways {
		err := a.istioClient.Delete(context.TODO(), &x)
		if err != nil {
			log.Printf("error delete %s/%s :  %v", x.GetNamespace(), x.GetName(), err)
		}
	}
}

func (a *IstioApplier) createGateways(gateways []istioapi.Gateway) {
	// cfg := config.GetConfigOrDie() //TODO: inject?
	// istioclientset := istioapiclient.NewForConfigOrDie(cfg)
	// for _, x := range gateways {
	// 	_, err := istioclientset.NetworkingV1alpha3().Gateways(x.ObjectMeta.GetNamespace()).Create(&x)
	// 	if err != nil {
	// 		log.Printf("error create %s/%s :  %v", x.ObjectMeta.GetNamespace(), x.ObjectMeta.GetName(), err)
	// 	}
	// }

	for _, x := range gateways {
		err := a.istioClient.Create(context.TODO(), &x)
		if err != nil {
			log.Printf("error create %s/%s :  %v", x.GetNamespace(), x.GetName(), err)
		}
	}
}

func (a *IstioApplier) updateGateways(gateways []struct {
	actual  istioapi.Gateway
	desired istioapi.Gateway
}) {
	// cfg := config.GetConfigOrDie() //TODO: inject?
	// istioclientset := istioapiclient.NewForConfigOrDie(cfg)
	// for _, x := range gateways {
	// 	_, err := istioclientset.NetworkingV1alpha3().Gateways(x.ObjectMeta.GetNamespace()).Update(&x)
	// 	if err != nil {
	// 		log.Printf("error update %s/%s :  %v", x.ObjectMeta.GetNamespace(), x.ObjectMeta.GetName(), err)
	// 	}
	// }

	for _, x := range gateways {
		tmp := x.actual.DeepCopy()
		tmp.Spec = x.desired.Spec
		err := a.istioClient.Update(context.TODO(), tmp)
		if err != nil {
			log.Printf("error update %s/%s :  %v", tmp.GetNamespace(), tmp.GetName(), err)
		}
	}
}

func (a *IstioApplier) deleteVirtualServices(virtualServices []istioapi.VirtualService) {
	for _, x := range virtualServices {
		err := a.istioClient.Delete(context.TODO(), &x)
		if err != nil {
			log.Printf("error delete %s/%s :  %v", x.GetNamespace(), x.GetName(), err)
		}
	}
}

func (a *IstioApplier) createVirtualServices(virtualServices []istioapi.VirtualService) {
	for _, x := range virtualServices {
		err := a.istioClient.Create(context.TODO(), &x)
		if err != nil {
			log.Printf("error create %s/%s :  %v", x.GetNamespace(), x.GetName(), err)
		}
	}
}

func (a *IstioApplier) updateVirtualServices(virtualServices []struct {
	actual  istioapi.VirtualService
	desired istioapi.VirtualService
}) {
	for _, x := range virtualServices {
		tmp := x.actual.DeepCopy()
		tmp.Spec = x.desired.Spec
		err := a.istioClient.Update(context.TODO(), tmp)
		if err != nil {
			log.Printf("error update %s/%s :  %v", tmp.GetNamespace(), tmp.GetName(), err)
		}
	}
}

func (a *IstioApplier) deleteDestinationRules(destinationRules []istioapi.DestinationRule) {
	for _, x := range destinationRules {
		err := a.istioClient.Delete(context.TODO(), &x)
		if err != nil {
			log.Printf("error delete %s/%s :  %v", x.GetNamespace(), x.GetName(), err)
		}
	}
}

func (a *IstioApplier) createDestinationRules(destinationRules []istioapi.DestinationRule) {
	for _, x := range destinationRules {
		err := a.istioClient.Create(context.TODO(), &x)
		if err != nil {
			log.Printf("error create %s/%s :  %v", x.GetNamespace(), x.GetName(), err)
		}
	}
}

func (a *IstioApplier) updateDestinationRules(destinationRules []struct {
	actual  istioapi.DestinationRule
	desired istioapi.DestinationRule
}) {
	for _, x := range destinationRules {
		tmp := x.actual.DeepCopy()
		tmp.Spec = x.desired.Spec
		err := a.istioClient.Update(context.TODO(), tmp)
		if err != nil {
			log.Printf("error update %s/%s :  %v", tmp.GetNamespace(), tmp.GetName(), err)
		}
	}
}

func triageDelete(existing []runtime.Object, desired []runtime.Object) (res []runtime.Object) {
	for _, e := range existing {
		if _, found := findObject(e, desired); !found {
			res = append(res, e)
		}
	}
	return res
}

func triageCreate(existing []runtime.Object, desired []runtime.Object) (res []runtime.Object) {
	for _, d := range desired {
		if _, found := findObject(d, existing); !found {
			res = append(res, d)
		}
	}
	return res
}

func triageUpdate(existing []runtime.Object, desired []runtime.Object) (res []interface{}) {
	for _, d := range desired {
		if obj, found := findObject(d, existing); found {
			//if obj == d { //TODO: compare spec
			res = append(res, struct {
				actual  runtime.Object
				desired runtime.Object
			}{
				obj,
				d,
			})
			//}
		}
	}
	return res
}

func findObject(obj runtime.Object, list []runtime.Object) (runtime.Object, bool) {
	for _, i := range list {
		iAccessor, _ := meta.Accessor(i)
		objAccessor, _ := meta.Accessor(obj)
		if iAccessor.GetName() == objAccessor.GetName() {
			return i, true
		}
	}
	return nil, false
}

func (a *IstioApplier) getDesiredResources(virtualApps []remeshv1alpha1.VirtualApp, namespace string) (gateways []istioapi.Gateway, virtualServices []istioapi.VirtualService, destinationRules []istioapi.DestinationRule) {
	gateways = make([]istioapi.Gateway, 0)
	virtualServices = make([]istioapi.VirtualService, 0)
	destinationRules = make([]istioapi.DestinationRule, 0)

	for _, virtualApp := range virtualApps {
		gateway, gatewayName := resources.MakeIstioGateway(virtualApp, namespace)
		gateways = append(gateways, gateway)

		httpRoutes := MakeRouteForVirtualAppConfig(virtualApp)
		gatewayVirtualService, virtualServiceName := resources.MakeIstioVirtualServiceForGateway(httpRoutes, namespace, gatewayName)
		virtualServices = append(virtualServices, gatewayVirtualService)

		layoutMap := getLayoutMapFromReleaseFlows(virtualApp.Spec.ReleaseFlows)

		transformedServices := TransformLayout(layoutMap)
		transformedVirtualServices := resources.MakeIstioVirtualServices(transformedServices, namespace, virtualServiceName)
		transformedDestinationRules := resources.MakeIstioDestinationRules(transformedServices, namespace)
		virtualServices = append(virtualServices, transformedVirtualServices...)
		destinationRules = append(destinationRules, transformedDestinationRules...)
	}
	return
}

func (a *IstioApplier) getActualResources(namespace string) ([]istioapi.Gateway, []istioapi.VirtualService, []istioapi.DestinationRule) {
	options := client.ListOptions{
		Namespace: namespace,
	}
	gatewayList := istioapi.GatewayList{}
	virtualServiceList := istioapi.VirtualServiceList{}
	destinationRuleList := istioapi.DestinationRuleList{}

	err := a.istioClient.List(context.TODO(), &options, &gatewayList)
	if err != nil {
		log.Printf("%v", err)
	}
	err = a.istioClient.List(context.TODO(), &options, &virtualServiceList)
	if err != nil {
		log.Printf("%v", err)
	}
	err = a.istioClient.List(context.TODO(), &options, &destinationRuleList)
	if err != nil {
		log.Printf("%v", err)
	}

	// cfg := config.GetConfigOrDie() //TODO: inject?
	// istioclientset := istioapiclient.NewForConfigOrDie(cfg)
	// gatewayList, err := istioclientset.NetworkingV1alpha3().Gateways(namespace).List(metav1.ListOptions{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// virtualServiceList, err := istioclientset.NetworkingV1alpha3().VirtualServices(namespace).List(metav1.ListOptions{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// destinationRuleList, err := istioclientset.NetworkingV1alpha3().DestinationRules(namespace).List(metav1.ListOptions{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return gatewayList.Items, virtualServiceList.Items, destinationRuleList.Items
}
