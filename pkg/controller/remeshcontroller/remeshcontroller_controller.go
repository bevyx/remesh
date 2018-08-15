/*
Copyright 2018 Bevyx.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package remeshcontroller

import (
	"context"
	"log"

	istioapi "github.com/bevyx/istio-api-go/pkg/apis/networking/v1alpha3"
	remeshv1alpha1 "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	"github.com/bevyx/remesh/pkg/istio"
	"github.com/bevyx/remesh/pkg/remesh"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new Remesh Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileRemesh{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("remeshcontroller-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &remeshv1alpha1.Entrypoint{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &remeshv1alpha1.Layout{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &remeshv1alpha1.Release{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &remeshv1alpha1.Segment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	//TODO: watch child resources
	// err = c.Watch(&source.Kind{Type: &istioapi.Gateway{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &remeshv1alpha1.Remesh{},
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

var _ reconcile.Reconciler = &ReconcileRemesh{}

// ReconcileRemesh reconciles a Remesh object
type ReconcileRemesh struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Remesh object and makes changes based on the state read
// and what is in the Remesh.Spec
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=remesh.bevyx.com,resources=entrypoints,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=remesh.bevyx.com,resources=targetings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=remesh.bevyx.com,resources=layouts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.istio.io,resources=gateways,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.istio.io,resources=virtualservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.istio.io,resources=destinationrules,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileRemesh) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	layoutList, releaseList, segmentList, entrypointList, err := r.fetchRemeshResources(request)
	if err != nil {
		return reconcile.Result{}, err
	}

	entrypointFlows := remesh.Combine(layoutList, releaseList, segmentList, entrypointList)
	applier := istio.NewIstioApplier(request.Namespace, r)
	istioapi.AddToScheme(r.scheme) //todo: this should be in the istio applier, but passing reconciler to istio would cause circular reference
	applier.Apply(entrypointFlows)

	return reconcile.Result{}, nil
}

func (r *ReconcileRemesh) fetchRemeshResources(request reconcile.Request) (layoutList remeshv1alpha1.LayoutList, releaseList remeshv1alpha1.ReleaseList, segmentList remeshv1alpha1.SegmentList, entrypointList remeshv1alpha1.EntrypointList, err error) {
	options := client.ListOptions{
		Namespace: request.Namespace,
	}
	entrypointList = remeshv1alpha1.EntrypointList{}
	layoutList = remeshv1alpha1.LayoutList{}
	releaseList = remeshv1alpha1.ReleaseList{}
	segmentList = remeshv1alpha1.SegmentList{}

	err = r.List(context.TODO(), &options, &entrypointList)
	if err != nil {
		log.Printf("missing Entrypoints %v", err)
		return
	}
	err = r.List(context.TODO(), &options, &layoutList)
	if err != nil {
		log.Printf("missing Layouts %v", err)
		return
	}
	err = r.List(context.TODO(), &options, &releaseList)
	if err != nil {
		log.Printf("missing Releases %v", err)
	}
	err = r.List(context.TODO(), &options, &segmentList)
	if err != nil {
		log.Printf("missing Segments %v", err)
		//it's ok to not have segments
	}
	log.Printf("fetched remesh resources: %d entrypoints, %d layouts, %d releases, %d segments", len(entrypointList.Items), len(layoutList.Items), len(releaseList.Items), len(segmentList.Items))
	return
}
