package istio

import (
	"reflect"
	"testing"

	"github.com/bevyx/remesh/pkg/adapters/istio/models"
	api "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
)

func TestTransformLayout(t *testing.T) {
	layoutMap := map[string]api.LayoutSpec{
		"bookinfo": {
			Http: []api.HTTPRoute{{
				Match: []api.HTTPMatchRequest{
					{
						Uri: &api.StringMatch{Exact: "/productpage"},
					},
					{
						Uri: &api.StringMatch{Exact: "/login"},
					},
					{
						Uri: &api.StringMatch{Exact: "/logout"},
					},
					{
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
				{
					Host: "productpage",
					Labels: map[string]string{
						"version": "v1",
						"stam":    "v1",
					},
				},
				{
					Host: "reviews",
					Labels: map[string]string{
						"version": "v1",
					},
				},
				{
					Host: "details",
					Labels: map[string]string{
						"version": "v1",
					},
				}},
		},
		"bookinfo-ratings": {
			Http: []api.HTTPRoute{{
				Match: []api.HTTPMatchRequest{
					{
						Uri: &api.StringMatch{Exact: "/productpage"},
					},
					{
						Uri: &api.StringMatch{Exact: "/login"},
					},
					{
						Uri: &api.StringMatch{Exact: "/logout"},
					},
					{
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
				{
					Host: "productpage",
					Labels: map[string]string{
						"stam":    "v1",
						"version": "v1",
					},
				},
				{
					Host: "reviews",
					Labels: map[string]string{
						"version": "v2",
					},
				},
				{
					Host: "ratings",
					Labels: map[string]string{
						"version": "v1",
					},
				},
				{
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

	expected := []models.TransformedService{{Host: "productpage", ServiceSubsetList: []models.ServiceSubset{models.ServiceSubset{Labels: map[string]string{"version": "v1", "stam": "v1"}, SubsetHash: "599866b499", Layouts: []string{"bookinfo", "bookinfo-ratings"}}}}, models.TransformedService{Host: "reviews", ServiceSubsetList: []models.ServiceSubset{models.ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", Layouts: []string{"bookinfo"}}, models.ServiceSubset{Labels: map[string]string{"version": "v2"}, SubsetHash: "84457d7684", Layouts: []string{"bookinfo-ratings"}}}}, models.TransformedService{Host: "details", ServiceSubsetList: []models.ServiceSubset{models.ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", Layouts: []string{"bookinfo", "bookinfo-ratings"}}}}, models.TransformedService{Host: "ratings", ServiceSubsetList: []models.ServiceSubset{models.ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", Layouts: []string{"bookinfo-ratings"}}}}}

	if !reflect.DeepEqual(transformLayout, expected) {
		t.Fail()
	}
}
