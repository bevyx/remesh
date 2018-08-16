package utils

import (
	"context"
	"reflect"

	remeshv1alpha1 "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func ReconcileAllVirtualAppsByFn(r client.Client, namespace string, reconciler func(*remeshv1alpha1.VirtualApp)) (reconcile.Result, error) {
	optionsNamespace := &client.ListOptions{
		Namespace: namespace,
	}

	virtualAppList := &remeshv1alpha1.VirtualAppList{}
	err := r.List(context.TODO(), optionsNamespace, virtualAppList)

	if err == nil {
		virtualApps := virtualAppList.Items
		virtualAppsCopy := virtualAppList.DeepCopy().Items
		for i := range virtualApps {
			reconciler(&(virtualApps[i]))
		}
		if reflect.DeepEqual(virtualApps, virtualAppsCopy) {
			return reconcile.Result{}, nil
		}
		for i := range virtualApps {
			if err := r.Update(context.TODO(), &(virtualApps[i])); err != nil {
				return reconcile.Result{}, err
			}
		}

	} else {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func VirtualAppHandlerMapper(mapper func(remeshv1alpha1.ReleaseFlow) []string) handler.ToRequestsFunc {
	return func(o handler.MapObject) []reconcile.Request {
		vapp := o.Object.(*remeshv1alpha1.VirtualApp)
		rfs := vapp.Spec.ReleaseFlows
		requests := make([]reconcile.Request, len(rfs))
		for _, rf := range rfs {
			names := mapper(rf)
			for _, name := range names {
				requests = append(requests, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      name,
						Namespace: vapp.Namespace,
					},
				})
			}

		}
		return requests
	}

}
