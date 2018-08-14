package istio

import (
	"reflect"
	"testing"

	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	"github.com/bevyx/remesh/pkg/istio/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestTransformLayout(t *testing.T) {
	layouts := []api.Layout{
		api.Layout{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bookinfo",
				Namespace: "default",
			},
			Spec: api.LayoutSpec{
				Http: []api.HTTPRoute{api.HTTPRoute{
					Match: []api.HTTPMatchRequest{
						api.HTTPMatchRequest{
							Uri: &api.StringMatch{Exact: "/productpage"},
						},
						api.HTTPMatchRequest{
							Uri: &api.StringMatch{Exact: "/login"},
						},
						api.HTTPMatchRequest{
							Uri: &api.StringMatch{Exact: "/logout"},
						},
						api.HTTPMatchRequest{
							Uri: &api.StringMatch{Prefix: "/api/v1/products"},
						},
					},
					Destination: api.Destination{
						Host: "productpage",
						Port: &api.PortSelector{
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
		api.Layout{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bookinfo-ratings",
				Namespace: "default",
			},
			Spec: api.LayoutSpec{
				Http: []api.HTTPRoute{api.HTTPRoute{
					Match: []api.HTTPMatchRequest{
						api.HTTPMatchRequest{
							Uri: &api.StringMatch{Exact: "/productpage"},
						},
						api.HTTPMatchRequest{
							Uri: &api.StringMatch{Exact: "/login"},
						},
						api.HTTPMatchRequest{
							Uri: &api.StringMatch{Exact: "/logout"},
						},
						api.HTTPMatchRequest{
							Uri: &api.StringMatch{Prefix: "/api/v1/products"},
						},
					},
					Destination: api.Destination{
						Host: "productpage",
						Port: &api.PortSelector{
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
	transformLayout := TransformLayout(layouts)
	//json, _ := json.MarshalIndent(transformLayout, "", "  ")
	//t.Logf("%s", string(json))
	//t.Logf("%#v\n", &transformLayout)

	expcted := []models.TransformedService{models.TransformedService{Host: "productpage", ServiceSubsetList: []models.ServiceSubset{models.ServiceSubset{Labels: map[string]string{"version": "v1", "stam": "v1"}, SubsetHash: "599866b499", Layouts: []string{"bookinfo", "bookinfo-ratings"}}}}, models.TransformedService{Host: "reviews", ServiceSubsetList: []models.ServiceSubset{models.ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", Layouts: []string{"bookinfo"}}, models.ServiceSubset{Labels: map[string]string{"version": "v2"}, SubsetHash: "84457d7684", Layouts: []string{"bookinfo-ratings"}}}}, models.TransformedService{Host: "details", ServiceSubsetList: []models.ServiceSubset{models.ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", Layouts: []string{"bookinfo", "bookinfo-ratings"}}}}, models.TransformedService{Host: "ratings", ServiceSubsetList: []models.ServiceSubset{models.ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", Layouts: []string{"bookinfo-ratings"}}}}}

	if !reflect.DeepEqual(transformLayout, expcted) {
		t.Fail()
	}
}
