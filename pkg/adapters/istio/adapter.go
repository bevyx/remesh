package istio

import (
	"context"
	"log"
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	goerrors "errors"

	istioapi "github.com/bevyx/istio-api-go/pkg/apis/networking/v1alpha3"
	"github.com/bevyx/remesh/pkg/adapters"
	remeshv1alpha1 "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

//IstioApplier is an implementation of the remesh.Applier interface that can apply a virtualApp configuration to Istio service mesh
type IstioApplier struct {
	istioClient client.Client
}

//NewIstioApplier creates new IstioApplier
func NewIstioApplier() adapters.Applier {
	clientConfig := config.GetConfigOrDie()
	istioScheme := runtime.NewScheme()
	istioapi.AddToScheme(istioScheme)
	istioClient, err := client.New(clientConfig, client.Options{Scheme: istioScheme})
	if err != nil {
		log.Printf("%v", err)
	}
	return &IstioApplier{istioClient: istioClient}
}

//Apply is implementing the adapters.Applier interface
func (a *IstioApplier) Apply(virtualApps []remeshv1alpha1.VirtualApp, namespace string) {
	log.Printf("applying %d virtualApps", len(virtualApps))
	desiredGateways, desiredVirtualServices, desiredDestinationRules := GetDesiredState(virtualApps, namespace)
	log.Printf("desired: gateways %d, virtualservices: %d, destinationrules %d", len(desiredGateways), len(desiredVirtualServices), len(desiredDestinationRules))
	actualGateways, actualVirtualServices, actualDestinationRules := a.getActualResources(namespace)
	log.Printf("actual: gateways %d, virtualservices: %d, destinationrules %d", len(actualGateways), len(actualVirtualServices), len(actualDestinationRules))

	a.applyGateways(desiredGateways, actualGateways)

	a.applyVirtualServices(desiredVirtualServices, actualVirtualServices)

	a.applyDestinationRules(desiredDestinationRules, actualDestinationRules)
}

func (a *IstioApplier) applyDestinationRules(desiredDestinationRules []istioapi.DestinationRule, actualDestinationRules []istioapi.DestinationRule) {
	desiredDestinationRulesObj := make([]runtime.Object, len(desiredDestinationRules))
	for i, x := range desiredDestinationRules {
		desiredDestinationRulesObj[i] = x.DeepCopyObject()
	}
	actualDestinationRulesObj := make([]runtime.Object, len(actualDestinationRules))
	for i, x := range actualDestinationRules {
		actualDestinationRulesObj[i] = x.DeepCopyObject()
	}
	destinationRulesToDeleteObj := triageDelete(actualDestinationRulesObj, desiredDestinationRulesObj)
	log.Printf("deleting %d destinationrules", len(destinationRulesToDeleteObj))
	a.deleteObjects(destinationRulesToDeleteObj)
	destinationRulesToCreateObj := triageCreate(actualDestinationRulesObj, desiredDestinationRulesObj)
	log.Printf("creating %d destinationrules", len(destinationRulesToCreateObj))
	a.createObjects(destinationRulesToCreateObj)
	destinationRulesToUpdateObj := triageUpdate(actualDestinationRulesObj, desiredDestinationRulesObj)
	log.Printf("updating %d destinationRules", len(destinationRulesToUpdateObj))
	a.updateObjects(destinationRulesToUpdateObj)
}

func (a *IstioApplier) applyVirtualServices(desiredVirtualServices []istioapi.VirtualService, actualVirtualServices []istioapi.VirtualService) {
	desiredVirtualServicesObj := make([]runtime.Object, len(desiredVirtualServices))
	for i, x := range desiredVirtualServices {
		desiredVirtualServicesObj[i] = x.DeepCopyObject()
	}
	actualVirtualServicesObj := make([]runtime.Object, len(actualVirtualServices))
	for i, x := range actualVirtualServices {
		actualVirtualServicesObj[i] = x.DeepCopyObject()
	}
	virtualServicesToDeleteObj := triageDelete(actualVirtualServicesObj, desiredVirtualServicesObj)
	log.Printf("deleting %d virtualservices", len(virtualServicesToDeleteObj))
	a.deleteObjects(virtualServicesToDeleteObj)
	virtualServicesToCreateObj := triageCreate(actualVirtualServicesObj, desiredVirtualServicesObj)
	log.Printf("creating %d virtualservices", len(virtualServicesToCreateObj))
	a.createObjects(virtualServicesToCreateObj)
	virtualServicesToUpdateObj := triageUpdate(actualVirtualServicesObj, desiredVirtualServicesObj)
	log.Printf("updating %d virtualServices", len(virtualServicesToUpdateObj))
	a.updateObjects(virtualServicesToUpdateObj)
}

func (a *IstioApplier) applyGateways(desiredGateways []istioapi.Gateway, actualGateways []istioapi.Gateway) {
	desiredGatewaysObj := make([]runtime.Object, len(desiredGateways))
	for i, x := range desiredGateways {
		desiredGatewaysObj[i] = x.DeepCopyObject()
	}
	actualGatewaysObj := make([]runtime.Object, len(actualGateways))
	for i, x := range actualGateways {
		actualGatewaysObj[i] = x.DeepCopyObject()
	}
	gatewaysToDeleteObj := triageDelete(actualGatewaysObj, desiredGatewaysObj)
	log.Printf("deleting %d gateways", len(gatewaysToDeleteObj))
	a.deleteObjects(gatewaysToDeleteObj)
	gatewaysToCreateObj := triageCreate(actualGatewaysObj, desiredGatewaysObj)
	log.Printf("creating %d gateways", len(gatewaysToCreateObj))
	a.createObjects(gatewaysToCreateObj)
	gatewaysToUpdateObj := triageUpdate(actualGatewaysObj, desiredGatewaysObj)
	log.Printf("updating %d gateways", len(gatewaysToUpdateObj))
	a.updateObjects(gatewaysToUpdateObj)
}

func (a *IstioApplier) deleteObjects(objects []runtime.Object) {
	for _, o := range objects {
		err := a.istioClient.Delete(context.TODO(), o)
		if err != nil {
			oMeta, _ := meta.Accessor(o)
			log.Printf("error delete %s/%s :  %v", oMeta.GetNamespace(), oMeta.GetName(), err)
		}
	}
}

func (a *IstioApplier) createObjects(objects []runtime.Object) {
	for _, o := range objects {
		err := a.istioClient.Create(context.TODO(), o)
		if err != nil {
			oMeta, _ := meta.Accessor(o)
			log.Printf("error create %s/%s :  %v", oMeta.GetNamespace(), oMeta.GetName(), err)
		}
	}
}

func (a *IstioApplier) updateObjects(objects []struct {
	actual  runtime.Object
	desired runtime.Object
}) {
	for _, o := range objects {
		actualCopy := o.actual.DeepCopyObject()
		aVal := reflect.ValueOf(actualCopy).Elem()
		actualSpec := aVal.FieldByName("Spec")

		if !actualSpec.IsValid() {
			aMeta, _ := meta.Accessor(actualCopy)
			err := goerrors.New("actual.Spec is invalid")
			log.Printf("error update %s/%s :  %v", aMeta.GetNamespace(), aMeta.GetName(), err)
			continue
		}
		dVal := reflect.ValueOf(o.desired).Elem()
		desiredSpec := dVal.FieldByName("Spec")
		if !(desiredSpec.IsValid() && desiredSpec.CanSet()) {
			dMeta, _ := meta.Accessor(o.desired)
			err := goerrors.New("desired.Spec is not set-able")
			log.Printf("error update %s/%s :  %v", dMeta.GetNamespace(), dMeta.GetName(), err)
			continue
		}

		actualSpec.Set(desiredSpec)
		err := a.istioClient.Update(context.TODO(), actualCopy)
		if err != nil {
			aMeta, _ := meta.Accessor(actualCopy)
			log.Printf("error update %s/%s :  %v", aMeta.GetNamespace(), aMeta.GetName(), err)
			continue
		}
	}
}

// existing && !desired -> delete
func triageDelete(existing []runtime.Object, desired []runtime.Object) (res []runtime.Object) {
	for _, e := range existing {
		if _, found := findObject(e, desired); !found {
			res = append(res, e)
		}
	}
	return res
}

// desired && !existing -> create
func triageCreate(existing []runtime.Object, desired []runtime.Object) (res []runtime.Object) {
	for _, d := range desired {
		if _, found := findObject(d, existing); !found {
			res = append(res, d)
		}
	}
	return res
}

// existing && desired -> update
func triageUpdate(existing []runtime.Object, desired []runtime.Object) (res []struct {
	actual  runtime.Object
	desired runtime.Object
}) {
	for _, d := range desired {
		dVal := reflect.ValueOf(d).Elem()
		desiredSpec := dVal.FieldByName("Spec")
		if obj, found := findObject(d, existing); found {
			aVal := reflect.ValueOf(obj).Elem()
			actualSpec := aVal.FieldByName("Spec")
			if !reflect.DeepEqual(actualSpec, desiredSpec) {
				res = append(res, struct {
					actual  runtime.Object
					desired runtime.Object
				}{
					obj,
					d,
				})
			}
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

	return gatewayList.Items, virtualServiceList.Items, destinationRuleList.Items
}
