package istio

import (
	"reflect"
	"testing"

	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestTransformVirtualEnvironment(t *testing.T) {
	virtualEnvironments := []api.VirtualEnvironment{
		api.VirtualEnvironment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bookinfo",
				Namespace: "default",
			},
			Spec: api.VirtualEnvironmentSpec{
				Http: []api.HTTPRoute{api.HTTPRoute{
					Match: []api.HTTPMatchRequest{
						api.HTTPMatchRequest{
							Uri: api.StringMatch{Exact: "/productpage"},
						},
						api.HTTPMatchRequest{
							Uri: api.StringMatch{Exact: "/login"},
						},
						api.HTTPMatchRequest{
							Uri: api.StringMatch{Exact: "/logout"},
						},
						api.HTTPMatchRequest{
							Uri: api.StringMatch{Prefix: "/api/v1/products"},
						},
					},
					DestinationRoute: api.DestinationRoute{
						Host: "productpage",
						Port: api.PortSelector{
							Number: 9080,
						},
					},
				}},
				Services: []api.Service{
					api.Service{
						Host: "productpage",
						Labels: map[string]string{
							"version": "v1",
							"stam":    "v1",
						},
					},
					api.Service{
						Host: "reviews",
						Labels: map[string]string{
							"version": "v1",
						},
					},
					api.Service{
						Host: "details",
						Labels: map[string]string{
							"version": "v1",
						},
					}},
			},
		},
		api.VirtualEnvironment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bookinfo-ratings",
				Namespace: "default",
			},
			Spec: api.VirtualEnvironmentSpec{
				Http: []api.HTTPRoute{api.HTTPRoute{
					Match: []api.HTTPMatchRequest{
						api.HTTPMatchRequest{
							Uri: api.StringMatch{Exact: "/productpage"},
						},
						api.HTTPMatchRequest{
							Uri: api.StringMatch{Exact: "/login"},
						},
						api.HTTPMatchRequest{
							Uri: api.StringMatch{Exact: "/logout"},
						},
						api.HTTPMatchRequest{
							Uri: api.StringMatch{Prefix: "/api/v1/products"},
						},
					},
					DestinationRoute: api.DestinationRoute{
						Host: "productpage",
						Port: api.PortSelector{
							Number: 9080,
						},
					},
				}},
				Services: []api.Service{
					api.Service{
						Host: "productpage",
						Labels: map[string]string{
							"stam":    "v1",
							"version": "v1",
						},
					},
					api.Service{
						Host: "reviews",
						Labels: map[string]string{
							"version": "v2",
						},
					},
					api.Service{
						Host: "ratings",
						Labels: map[string]string{
							"version": "v1",
						},
					},
					api.Service{
						Host: "details",
						Labels: map[string]string{
							"version": "v1",
						},
					}},
			},
		},
	}
	transformVirtualEnvironment := TransformVirtualEnvironment(virtualEnvironments)
	//json, _ := json.MarshalIndent(transformVirtualEnvironment, "", "  ")
	//t.Logf("%s", string(json))
	//t.Logf("%#v\n", &transformVirtualEnvironment)

	expcted := []TransformedService{TransformedService{Host: "productpage", ServiceSubsetList: []ServiceSubset{ServiceSubset{Labels: map[string]string{"version": "v1", "stam": "v1"}, SubsetHash: "599866b499", VirtualEnvironments: []string{"bookinfo", "bookinfo-ratings"}}}}, TransformedService{Host: "reviews", ServiceSubsetList: []ServiceSubset{ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", VirtualEnvironments: []string{"bookinfo"}}, ServiceSubset{Labels: map[string]string{"version": "v2"}, SubsetHash: "84457d7684", VirtualEnvironments: []string{"bookinfo-ratings"}}}}, TransformedService{Host: "details", ServiceSubsetList: []ServiceSubset{ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", VirtualEnvironments: []string{"bookinfo", "bookinfo-ratings"}}}}, TransformedService{Host: "ratings", ServiceSubsetList: []ServiceSubset{ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", VirtualEnvironments: []string{"bookinfo-ratings"}}}}}

	if !reflect.DeepEqual(transformVirtualEnvironment, expcted) {
		t.Fail()
	}
}
