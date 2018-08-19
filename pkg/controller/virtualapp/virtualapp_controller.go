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

package virtualapp

import (
	"context"
	"log"

	"github.com/bevyx/remesh/pkg/adapters"
	"github.com/bevyx/remesh/pkg/adapters/istio"
	remeshv1alpha1 "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new VirtualApp Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this remesh.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileVirtualApp{Client: mgr.GetClient(), scheme: mgr.GetScheme(), applier: istio.NewIstioApplier()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("virtualapp-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to VirtualApp
	err = c.Watch(&source.Kind{Type: &remeshv1alpha1.VirtualApp{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by VirtualApp - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &remeshv1alpha1.VirtualApp{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileVirtualApp{}

// ReconcileVirtualApp reconciles a VirtualApp object
type ReconcileVirtualApp struct {
	client.Client
	scheme  *runtime.Scheme
	applier adapters.Applier
}

// Reconcile reads that state of the cluster for a VirtualApp object and makes changes based on the state read
// and what is in the VirtualApp.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=remesh.bevyx.com,resources=virtualapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.istio.io,resources=gateways,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.istio.io,resources=virtualservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.istio.io,resources=destinationrules,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileVirtualApp) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the VirtualApp instance
	virtualAppList, err := r.getVirtualAppList(request.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}

	r.applier.Apply(virtualAppList.Items, request.Namespace)
	return reconcile.Result{}, nil
}

func (r *ReconcileVirtualApp) getVirtualAppList(namespace string) (virtualAppList remeshv1alpha1.VirtualAppList, err error) {
	options := client.ListOptions{
		Namespace: namespace,
	}
	virtualAppList = remeshv1alpha1.VirtualAppList{}

	err = r.List(context.TODO(), &options, &virtualAppList)
	if err != nil {
		log.Printf("missing virtualApps %v", err)
		return
	}
	log.Printf("fetched virtualApps: %d", len(virtualAppList.Items))
	return
}
