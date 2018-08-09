package istio

// import (
// 	"reflect"
// 	"testing"

// 	"github.com/davecgh/go-spew/spew"

// 	istioapi "github.com/bevyx/istio-api-go/pkg/apis/networking/v1alpha3"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/runtime"
// )

// var (
// 	gw1 = istioapi.Gateway{
// 		TypeMeta: metav1.TypeMeta{
// 			Kind:       "Gateway",
// 			APIVersion: "networking.istio.io/v1alpha3",
// 		},
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "gw1",
// 			Namespace: "ns1",
// 		},
// 		Spec: istioapi.GatewaySpec{
// 			Selector: map[string]string{
// 				"istio": "ingressgateway",
// 			},
// 			Servers: []istioapi.Server{
// 				istioapi.Server{
// 					Port: &istioapi.Port{
// 						Number:   8080,
// 						Protocol: "HTTP",
// 						Name:     "http",
// 					},
// 					Hosts: []string{
// 						"*",
// 					},
// 				},
// 			},
// 		},
// 	}
// 	gw2 = istioapi.Gateway{
// 		TypeMeta: metav1.TypeMeta{
// 			Kind:       "Gateway",
// 			APIVersion: "networking.istio.io/v1alpha3",
// 		},
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "gw2",
// 			Namespace: "ns1",
// 		},
// 		Spec: istioapi.GatewaySpec{
// 			Selector: map[string]string{
// 				"istio": "ingressgateway",
// 			},
// 			Servers: []istioapi.Server{
// 				istioapi.Server{
// 					Port: &istioapi.Port{
// 						Number:   8080,
// 						Protocol: "HTTP",
// 						Name:     "http",
// 					},
// 					Hosts: []string{
// 						"*",
// 					},
// 				},
// 			},
// 		},
// 	}
// 	gw3 = istioapi.Gateway{
// 		TypeMeta: metav1.TypeMeta{
// 			Kind:       "Gateway",
// 			APIVersion: "networking.istio.io/v1alpha3",
// 		},
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "gw3",
// 			Namespace: "ns1",
// 		},
// 		Spec: istioapi.GatewaySpec{
// 			Selector: map[string]string{
// 				"istio": "anothergateway",
// 			},
// 			Servers: []istioapi.Server{
// 				istioapi.Server{
// 					Port: &istioapi.Port{
// 						Number:   8080,
// 						Protocol: "HTTP",
// 						Name:     "http",
// 					},
// 					Hosts: []string{
// 						"*",
// 					},
// 				},
// 			},
// 		},
// 	}
// )

// func TestToDelete(t *testing.T) {
// 	existing := []istioapi.Gateway{
// 		*gw1.DeepCopy(),
// 	}
// 	desired := []istioapi.Gateway{
// 		*gw1.DeepCopy(),
// 	}
// 	existingObj := make([]runtime.Object, len(existing))
// 	for i, e := range existing {
// 		existingObj[i] = runtime.Object(&e)
// 	}
// 	desiredObj := make([]runtime.Object, len(desired))
// 	for i, d := range desired {
// 		desiredObj[i] = runtime.Object(&d)
// 	}
// 	actualObj := triageDelete(existingObj, desiredObj)
// 	var expected []runtime.Object
// 	if !reflect.DeepEqual(actualObj, expected) {
// 		t.Fail()
// 	}

// 	// --------------------

// 	existing = []istioapi.Gateway{
// 		*gw1.DeepCopy(),
// 	}
// 	existing[0].Name = "another"
// 	desired = []istioapi.Gateway{
// 		*gw1.DeepCopy(),
// 	}
// 	existingObj = make([]runtime.Object, len(existing))
// 	for i, e := range existing {
// 		existingObj[i] = runtime.Object(&e)
// 	}
// 	desiredObj = make([]runtime.Object, len(desired))
// 	for i, d := range desired {
// 		desiredObj[i] = runtime.Object(&d)
// 	}
// 	actualObj = triageDelete(existingObj, desiredObj)
// 	expected = []runtime.Object{
// 		runtime.Object(&existing[0]),
// 	}
// 	if !reflect.DeepEqual(actualObj, expected) {
// 		t.Fail()
// 	}

// }

// func TestToCreate(t *testing.T) {
// 	existing := []istioapi.Gateway{
// 		*gw1.DeepCopy(),
// 	}
// 	desired := []istioapi.Gateway{
// 		*gw1.DeepCopy(),
// 	}
// 	existingObj := make([]runtime.Object, len(existing))
// 	for i, e := range existing {
// 		existingObj[i] = runtime.Object(&e)
// 	}
// 	desiredObj := make([]runtime.Object, len(desired))
// 	for i, d := range desired {
// 		desiredObj[i] = runtime.Object(&d)
// 	}
// 	actualObj := triageCreate(existingObj, desiredObj)
// 	var expected []runtime.Object
// 	if !reflect.DeepEqual(actualObj, expected) {
// 		t.Fail()
// 	}

// 	// --------------------

// 	existing = []istioapi.Gateway{
// 		*gw1.DeepCopy(),
// 	}
// 	existing[0].Name = "another"
// 	desired = []istioapi.Gateway{
// 		*gw1.DeepCopy(),
// 	}
// 	existingObj = make([]runtime.Object, len(existing))
// 	for i, e := range existing {
// 		existingObj[i] = runtime.Object(&e)
// 	}
// 	desiredObj = make([]runtime.Object, len(desired))
// 	for i, d := range desired {
// 		desiredObj[i] = runtime.Object(&d)
// 	}
// 	actualObj = triageCreate(existingObj, desiredObj)
// 	expected = []runtime.Object{
// 		runtime.Object(&desired[0]),
// 	}

// 	if !reflect.DeepEqual(actualObj, expected) {
// 		t.Fail()
// 	}
// }

// func TestGetActualResources(t *testing.T) {
// 	gateways, virtualServices, destinationRules := getActualResources("istio-system")
// 	spew.Dump(gateways, virtualServices, destinationRules)
// }
