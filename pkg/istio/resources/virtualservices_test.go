package resources

// import (
// 	"encoding/json"
// 	"testing"
// )

// func TestMakeIstioVirtualServices(t *testing.T) {
// 	transformedServices := []TransformedService{TransformedService{Host: "productpage", ServiceSubsetList: []ServiceSubset{ServiceSubset{Labels: map[string]string{"version": "v1", "stam": "v1"}, SubsetHash: "599866b499", Layouts: []string{"bookinfo", "bookinfo-ratings"}}}}, TransformedService{Host: "reviews", ServiceSubsetList: []ServiceSubset{ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", Layouts: []string{"bookinfo"}}, ServiceSubset{Labels: map[string]string{"version": "v2"}, SubsetHash: "84457d7684", Layouts: []string{"bookinfo-ratings"}}}}, TransformedService{Host: "details", ServiceSubsetList: []ServiceSubset{ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", Layouts: []string{"bookinfo", "bookinfo-ratings"}}}}, TransformedService{Host: "ratings", ServiceSubsetList: []ServiceSubset{ServiceSubset{Labels: map[string]string{"version": "v1"}, SubsetHash: "8444f85f99", Layouts: []string{"bookinfo-ratings"}}}}}

// 	istioVirtualServices := MakeIstioVirtualServices(transformedServices, "default", "bookinfo-gateway")
// 	json, _ := json.MarshalIndent(istioVirtualServices, "", "  ")
// 	t.Logf("%s", string(json))
// 	//t.Logf("%#v\n", istioVirtualServices)

// 	//expcted := []istioapi.VirtualService{istioapi.VirtualService{TypeMeta: metav1.TypeMeta{Kind: "VirtualService", APIVersion: "networking.istio.io/istioapi"}, ObjectMeta: metav1.ObjectMeta{Name: "productpage", GenerateName: "", Namespace: "default", SelfLink: "", UID: "", ResourceVersion: "", Generation: 0, CreationTimestamp: metav1.Time{Time: time.Time{wall: 0x0, ext: 0, loc: (*time.Location)(nil)}}, DeletionTimestamp: (*metav1.Time)(nil), DeletionGracePeriodSeconds: (*int64)(nil), Labels: map[string]string{}, Annotations: map[string]string(nil), OwnerReferences: []metav1.OwnerReference(nil), Initializers: (*metav1.Initializers)(nil), Finalizers: []string(nil), ClusterName: ""}, Spec: istioapi.VirtualServiceSpec{Hosts: []string{"productpage"}, Gateways: []string{"bookinfo-gateway", "mesh"}, Http: []istioapi.HTTPRoute{istioapi.HTTPRoute{Match: []istioapi.HTTPMatchRequest{istioapi.HTTPMatchRequest{Uri: (*istioapi.StringMatch)(nil), Scheme: (*istioapi.StringMatch)(nil), Method: (*istioapi.StringMatch)(nil), Authority: (*istioapi.StringMatch)(nil), Headers: map[string]istioapi.StringMatch{"ol-route": istioapi.StringMatch{Exact: "bookinfo", Prefix: "", Regex: ""}}}}, Route: []istioapi.DestinationWeight{istioapi.DestinationWeight{Destination: istioapi.Destination{Host: "productpage", Subset: "productpage-599866b499", Port: istioapi.PortSelector{Number: 0x0, Name: ""}}, Weight: 0}}, Redirect: (*istioapi.HTTPRedirect)(nil), Rewrite: (*istioapi.HTTPRewrite)(nil), WebsocketUpgrade: false, Timeout: "", Retries: (*istioapi.HTTPRetry)(nil), Fault: (*istioapi.HTTPFaultInjection)(nil), Mirror: (*istioapi.Destination)(nil), AppendHeaders: map[string]string(nil), RemoveResponseHeaders: map[string]string(nil)}, istioapi.HTTPRoute{Match: []istioapi.HTTPMatchRequest{istioapi.HTTPMatchRequest{Uri: (*istioapi.StringMatch)(nil), Scheme: (*istioapi.StringMatch)(nil), Method: (*istioapi.StringMatch)(nil), Authority: (*istioapi.StringMatch)(nil), Headers: map[string]istioapi.StringMatch{"ol-route": istioapi.StringMatch{Exact: "bookinfo-ratings", Prefix: "", Regex: ""}}}}, Route: []istioapi.DestinationWeight{istioapi.DestinationWeight{Destination: istioapi.Destination{Host: "productpage", Subset: "productpage-599866b499", Port: istioapi.PortSelector{Number: 0x0, Name: ""}}, Weight: 0}}, Redirect: (*istioapi.HTTPRedirect)(nil), Rewrite: (*istioapi.HTTPRewrite)(nil), WebsocketUpgrade: false, Timeout: "", Retries: (*istioapi.HTTPRetry)(nil), Fault: (*istioapi.HTTPFaultInjection)(nil), Mirror: (*istioapi.Destination)(nil), AppendHeaders: map[string]string(nil), RemoveResponseHeaders: map[string]string(nil)}}, Tcp: []istioapi.TCPRoute(nil)}}, istioapi.VirtualService{TypeMeta: metav1.TypeMeta{Kind: "VirtualService", APIVersion: "networking.istio.io/istioapi"}, ObjectMeta: metav1.ObjectMeta{Name: "reviews", GenerateName: "", Namespace: "default", SelfLink: "", UID: "", ResourceVersion: "", Generation: 0, CreationTimestamp: metav1.Time{Time: time.Time{wall: 0x0, ext: 0, loc: (*time.Location)(nil)}}, DeletionTimestamp: (*metav1.Time)(nil), DeletionGracePeriodSeconds: (*int64)(nil), Labels: map[string]string{}, Annotations: map[string]string(nil), OwnerReferences: []metav1.OwnerReference(nil), Initializers: (*metav1.Initializers)(nil), Finalizers: []string(nil), ClusterName: ""}, Spec: istioapi.VirtualServiceSpec{Hosts: []string{"reviews"}, Gateways: []string{"bookinfo-gateway", "mesh"}, Http: []istioapi.HTTPRoute{istioapi.HTTPRoute{Match: []istioapi.HTTPMatchRequest{istioapi.HTTPMatchRequest{Uri: (*istioapi.StringMatch)(nil), Scheme: (*istioapi.StringMatch)(nil), Method: (*istioapi.StringMatch)(nil), Authority: (*istioapi.StringMatch)(nil), Headers: map[string]istioapi.StringMatch{"ol-route": istioapi.StringMatch{Exact: "bookinfo", Prefix: "", Regex: ""}}}}, Route: []istioapi.DestinationWeight{istioapi.DestinationWeight{Destination: istioapi.Destination{Host: "reviews", Subset: "reviews-8444f85f99", Port: istioapi.PortSelector{Number: 0x0, Name: ""}}, Weight: 0}}, Redirect: (*istioapi.HTTPRedirect)(nil), Rewrite: (*istioapi.HTTPRewrite)(nil), WebsocketUpgrade: false, Timeout: "", Retries: (*istioapi.HTTPRetry)(nil), Fault: (*istioapi.HTTPFaultInjection)(nil), Mirror: (*istioapi.Destination)(nil), AppendHeaders: map[string]string(nil), RemoveResponseHeaders: map[string]string(nil)}, istioapi.HTTPRoute{Match: []istioapi.HTTPMatchRequest{istioapi.HTTPMatchRequest{Uri: (*istioapi.StringMatch)(nil), Scheme: (*istioapi.StringMatch)(nil), Method: (*istioapi.StringMatch)(nil), Authority: (*istioapi.StringMatch)(nil), Headers: map[string]istioapi.StringMatch{"ol-route": istioapi.StringMatch{Exact: "bookinfo-ratings", Prefix: "", Regex: ""}}}}, Route: []istioapi.DestinationWeight{istioapi.DestinationWeight{Destination: istioapi.Destination{Host: "reviews", Subset: "reviews-84457d7684", Port: istioapi.PortSelector{Number: 0x0, Name: ""}}, Weight: 0}}, Redirect: (*istioapi.HTTPRedirect)(nil), Rewrite: (*istioapi.HTTPRewrite)(nil), WebsocketUpgrade: false, Timeout: "", Retries: (*istioapi.HTTPRetry)(nil), Fault: (*istioapi.HTTPFaultInjection)(nil), Mirror: (*istioapi.Destination)(nil), AppendHeaders: map[string]string(nil), RemoveResponseHeaders: map[string]string(nil)}}, Tcp: []istioapi.TCPRoute(nil)}}, istioapi.VirtualService{TypeMeta: metav1.TypeMeta{Kind: "VirtualService", APIVersion: "networking.istio.io/istioapi"}, ObjectMeta: metav1.ObjectMeta{Name: "details", GenerateName: "", Namespace: "default", SelfLink: "", UID: "", ResourceVersion: "", Generation: 0, CreationTimestamp: metav1.Time{Time: time.Time{wall: 0x0, ext: 0, loc: (*time.Location)(nil)}}, DeletionTimestamp: (*metav1.Time)(nil), DeletionGracePeriodSeconds: (*int64)(nil), Labels: map[string]string{}, Annotations: map[string]string(nil), OwnerReferences: []metav1.OwnerReference(nil), Initializers: (*metav1.Initializers)(nil), Finalizers: []string(nil), ClusterName: ""}, Spec: istioapi.VirtualServiceSpec{Hosts: []string{"details"}, Gateways: []string{"bookinfo-gateway", "mesh"}, Http: []istioapi.HTTPRoute{istioapi.HTTPRoute{Match: []istioapi.HTTPMatchRequest{istioapi.HTTPMatchRequest{Uri: (*istioapi.StringMatch)(nil), Scheme: (*istioapi.StringMatch)(nil), Method: (*istioapi.StringMatch)(nil), Authority: (*istioapi.StringMatch)(nil), Headers: map[string]istioapi.StringMatch{"ol-route": istioapi.StringMatch{Exact: "bookinfo", Prefix: "", Regex: ""}}}}, Route: []istioapi.DestinationWeight{istioapi.DestinationWeight{Destination: istioapi.Destination{Host: "details", Subset: "details-8444f85f99", Port: istioapi.PortSelector{Number: 0x0, Name: ""}}, Weight: 0}}, Redirect: (*istioapi.HTTPRedirect)(nil), Rewrite: (*istioapi.HTTPRewrite)(nil), WebsocketUpgrade: false, Timeout: "", Retries: (*istioapi.HTTPRetry)(nil), Fault: (*istioapi.HTTPFaultInjection)(nil), Mirror: (*istioapi.Destination)(nil), AppendHeaders: map[string]string(nil), RemoveResponseHeaders: map[string]string(nil)}, istioapi.HTTPRoute{Match: []istioapi.HTTPMatchRequest{istioapi.HTTPMatchRequest{Uri: (*istioapi.StringMatch)(nil), Scheme: (*istioapi.StringMatch)(nil), Method: (*istioapi.StringMatch)(nil), Authority: (*istioapi.StringMatch)(nil), Headers: map[string]istioapi.StringMatch{"ol-route": istioapi.StringMatch{Exact: "bookinfo-ratings", Prefix: "", Regex: ""}}}}, Route: []istioapi.DestinationWeight{istioapi.DestinationWeight{Destination: istioapi.Destination{Host: "details", Subset: "details-8444f85f99", Port: istioapi.PortSelector{Number: 0x0, Name: ""}}, Weight: 0}}, Redirect: (*istioapi.HTTPRedirect)(nil), Rewrite: (*istioapi.HTTPRewrite)(nil), WebsocketUpgrade: false, Timeout: "", Retries: (*istioapi.HTTPRetry)(nil), Fault: (*istioapi.HTTPFaultInjection)(nil), Mirror: (*istioapi.Destination)(nil), AppendHeaders: map[string]string(nil), RemoveResponseHeaders: map[string]string(nil)}}, Tcp: []istioapi.TCPRoute(nil)}}, istioapi.VirtualService{TypeMeta: metav1.TypeMeta{Kind: "VirtualService", APIVersion: "networking.istio.io/istioapi"}, ObjectMeta: metav1.ObjectMeta{Name: "ratings", GenerateName: "", Namespace: "default", SelfLink: "", UID: "", ResourceVersion: "", Generation: 0, CreationTimestamp: metav1.Time{Time: time.Time{wall: 0x0, ext: 0, loc: (*time.Location)(nil)}}, DeletionTimestamp: (*metav1.Time)(nil), DeletionGracePeriodSeconds: (*int64)(nil), Labels: map[string]string{}, Annotations: map[string]string(nil), OwnerReferences: []metav1.OwnerReference(nil), Initializers: (*metav1.Initializers)(nil), Finalizers: []string(nil), ClusterName: ""}, Spec: istioapi.VirtualServiceSpec{Hosts: []string{"ratings"}, Gateways: []string{"bookinfo-gateway", "mesh"}, Http: []istioapi.HTTPRoute{istioapi.HTTPRoute{Match: []istioapi.HTTPMatchRequest{istioapi.HTTPMatchRequest{Uri: (*istioapi.StringMatch)(nil), Scheme: (*istioapi.StringMatch)(nil), Method: (*istioapi.StringMatch)(nil), Authority: (*istioapi.StringMatch)(nil), Headers: map[string]istioapi.StringMatch{"ol-route": istioapi.StringMatch{Exact: "bookinfo-ratings", Prefix: "", Regex: ""}}}}, Route: []istioapi.DestinationWeight{istioapi.DestinationWeight{Destination: istioapi.Destination{Host: "ratings", Subset: "ratings-8444f85f99", Port: istioapi.PortSelector{Number: 0x0, Name: ""}}, Weight: 0}}, Redirect: (*istioapi.HTTPRedirect)(nil), Rewrite: (*istioapi.HTTPRewrite)(nil), WebsocketUpgrade: false, Timeout: "", Retries: (*istioapi.HTTPRetry)(nil), Fault: (*istioapi.HTTPFaultInjection)(nil), Mirror: (*istioapi.Destination)(nil), AppendHeaders: map[string]string(nil), RemoveResponseHeaders: map[string]string(nil)}}, Tcp: []istioapi.TCPRoute(nil)}}}
// 	//if !reflect.DeepEqual(*istioVirtualServices, expcted) {
// 	//	t.Fail()
// 	//}
// }
