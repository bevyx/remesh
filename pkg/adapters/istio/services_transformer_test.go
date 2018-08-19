package istio

import (
	"reflect"
	"testing"

	"github.com/bevyx/remesh/pkg/adapters/istio/models"
	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
)

func TestTransformLayout(t *testing.T) {
	layoutMap := map[string]api.LayoutSpec{
		"bookinfo": api.LayoutSpec{
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
		"bookinfo-ratings": api.LayoutSpec{
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
	}
	transformLayout := TransformLayout(layoutMap)
	//json, _ := json.MarshalIndent(transformLayout, "", "  ")
	//t.Logf("%s", string(json))
	//t.Logf("%#v\n", &transformLayout)

	expcted := []models.TransformedService{models.TransformedService{Host: "productpage", ServiceSubsetList: []models.ServiceSubset{models.ServiceSubset{Labels: map[string]string{"version": "v1", "stam": "v1"}, SubsetHash: "599866b499", Layouts: []string{"bookinfo", "bookinfo-ratings"}}}}, models.TransformedService{Host: "reviews", ServiceSubsetList: []models.ServiceSubset{models.ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", Layouts: []string{"bookinfo"}}, models.ServiceSubset{Labels: map[string]string{"version": "v2"}, SubsetHash: "84457d7684", Layouts: []string{"bookinfo-ratings"}}}}, models.TransformedService{Host: "details", ServiceSubsetList: []models.ServiceSubset{models.ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", Layouts: []string{"bookinfo", "bookinfo-ratings"}}}}, models.TransformedService{Host: "ratings", ServiceSubsetList: []models.ServiceSubset{models.ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", Layouts: []string{"bookinfo-ratings"}}}}}

	if !reflect.DeepEqual(transformLayout, expcted) {
		t.Fail()
	}
}
