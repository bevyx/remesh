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

package virtualappconfig

import (
	"context"

	remeshv1alpha1 "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	// istioapi "github.com/bevyx/istio-api-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new VirtualAppConfig Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileVirtualAppConfig{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("virtualappconfig-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to VirtualAppConfig
	err = c.Watch(&source.Kind{Type: &remeshv1alpha1.VirtualAppConfig{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to Gateway
	// err = c.Watch(&source.Kind{Type: &istioapi.Gateaway{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &remeshv1alpha1.VirtualAppConfig{},
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

var _ reconcile.Reconciler = &ReconcileVirtualAppConfig{}

// ReconcileVirtualAppConfig reconciles a VirtualAppConfig object
type ReconcileVirtualAppConfig struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a VirtualAppConfig object and makes changes based on the state read
// and what is in the VirtualAppConfig.Spec
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=remesh.bevyx.com,resources=virtualappconfigs,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileVirtualAppConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the VirtualAppConfig instance
	instance := &remeshv1alpha1.VirtualAppConfig{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	//Generate gateway from virtualappconfig and link them
	// gateway, gatewayName := remesh.GenerateIstioGateway(*instance, instance.ObjectMeta.Namespace)
	// if err := controllerutil.SetControllerReference(instance, gateway, r.scheme); err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Check if the Deployment already exists
	// found := &istioapi.Gateway{}
	// err = r.Get(context.TODO(), types.NamespacedName{Name: gatewayName, Namespace: gateway.Namespace}, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	log.Printf("Creating Gateaway %s/%s\n", gateway.Namespace, gatewayName)
	// 	err = r.Create(context.TODO(), gateway)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}
	// } else if err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Update the found object and write the result back if there are any changes
	// if !reflect.DeepEqual(gateway.Spec, found.Spec) {
	// 	found.Spec = deploy.Spec
	// 	log.Printf("Updating Gateway %s/%s\n", gateway.Namespace, gatewayName)
	// 	err = r.Update(context.TODO(), found)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}
	// }
	return reconcile.Result{}, nil
}
