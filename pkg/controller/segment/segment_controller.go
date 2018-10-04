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

package segment

import (
	"context"

	remeshv1alpha1 "github.com/bevyx/remesh/pkg/apis/remesh/v1alpha1"
	"github.com/bevyx/remesh/pkg/controller/utils"
	"k8s.io/apimachinery/pkg/api/errors"
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

// Add creates a new Segment Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this remesh.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSegment{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("segment-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Segment
	err = c.Watch(&source.Kind{Type: &remeshv1alpha1.Segment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to VirtualApp
	vappMapper := utils.VirtualAppHandlerMapper(func(rf remeshv1alpha1.ReleaseFlow) []string {
		names := make([]string, 0)
		if rf.Segments != nil {
			for name := range *rf.Segments {
				names = append(names, name)
			}
		}
		return names
	})
	err = c.Watch(&source.Kind{Type: &remeshv1alpha1.VirtualApp{}}, &handler.EnqueueRequestsFromMapFunc{ToRequests: vappMapper})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileSegment{}

// ReconcileSegment reconciles a Segment object
type ReconcileSegment struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Segment object and makes changes based on the state read
// and what is in the Segment.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// +kubebuilder:rbac:groups=remesh.bevyx.com,resources=segments,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileSegment) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Segment instance
	instance := &remeshv1alpha1.Segment{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	deleted := false
	if err != nil {
		if errors.IsNotFound(err) {
			deleted = true
		} else {
			// Error reading the object - requeue the request.
			return reconcile.Result{}, err
		}
	}

	return utils.ReconcileAllVirtualAppsByFn(r.Client, request.Namespace, func(virtualApp *remeshv1alpha1.VirtualApp) {
		reconcileSegment(request.Name, instance, virtualApp, deleted)
	})
}

func reconcileSegment(segmentName string, segment *remeshv1alpha1.Segment, virtualApp *remeshv1alpha1.VirtualApp, deleted bool) {
	releaseFlows := virtualApp.Spec.ReleaseFlows
	for i := range releaseFlows {
		reconcileReleaseFlowSegments(releaseFlows[i].Segments, segmentName, deleted, segment)
	}
}

func reconcileReleaseFlowSegments(releaseFlowSegments *map[string]*remeshv1alpha1.SegmentSpec, segmentName string, deleted bool, segment *remeshv1alpha1.Segment) {
	if releaseFlowSegments != nil {
		if _, ok := (*releaseFlowSegments)[segmentName]; ok {
			if deleted {
				(*releaseFlowSegments)[segmentName] = nil
			} else {
				(*releaseFlowSegments)[segmentName] = segment.Spec.DeepCopy()
			}
		}
	}
}
