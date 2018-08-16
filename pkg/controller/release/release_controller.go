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

package release

import (
	"context"
	goerrors "errors"
	"reflect"
	"sort"

	"k8s.io/apimachinery/pkg/types"

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

// Add creates a new Release Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this remesh.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileRelease{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("release-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Release
	err = c.Watch(&source.Kind{Type: &remeshv1alpha1.Release{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to VirtualApp
	vappMapper := utils.VirtualAppHandlerMapper(func(rf remeshv1alpha1.ReleaseFlow) []string {
		return []string{rf.ReleaseName}
	})
	err = c.Watch(&source.Kind{Type: &remeshv1alpha1.VirtualApp{}}, &handler.EnqueueRequestsFromMapFunc{ToRequests: vappMapper})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileRelease{}

// ReconcileRelease reconciles a Release object
type ReconcileRelease struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Release object and makes changes based on the state read
// and what is in the Release.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// +kubebuilder:rbac:groups=remesh.bevyx.com,resources=releases,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileRelease) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Release instance
	instance := &remeshv1alpha1.Release{}
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

	namespacedName := types.NamespacedName{
		Namespace: request.Namespace,
		Name:      instance.Spec.VirtualAppConfig,
	}
	virtualApp := &remeshv1alpha1.VirtualApp{}
	err = r.Get(context.TODO(), namespacedName, virtualApp)

	if errors.IsNotFound(err) {
		tmp := instance.DeepCopy()
		tmp.Status.Phase = "Pending"
		tmp.Status.Reason = "missingVApp"
		r.Update(context.TODO(), tmp)
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	} else {
		//Ensure there is a single default release
		if instance.Spec.Targeting == nil || len(instance.Spec.Targeting.Segments) == 0 {
			releases := remeshv1alpha1.ReleaseList{}
			err := r.List(context.TODO(), &client.ListOptions{}, &releases)
			if err != nil {
				return reconcile.Result{}, err
			}
			//TODO: figure out how to move this to serverside by query
			for _, release := range releases.Items {
				if release.Name != instance.Name && (release.Spec.Targeting == nil || len(release.Spec.Targeting.Segments) == 0) {
					r.Delete(context.TODO(), instance) //ignore errors here. if couldnt delete, it will be tried again next time
					return reconcile.Result{}, goerrors.New("Default release already exists")

				}
			}
		}
		virtualAppCopy := virtualApp.DeepCopy()
		reconcileRelease(request.Name, instance, virtualApp, deleted)

		if reflect.DeepEqual(virtualApp.Spec, virtualAppCopy.Spec) {
			return reconcile.Result{}, nil
		}
		if err := r.Update(context.TODO(), virtualApp); err != nil {
			return reconcile.Result{}, err
		}
	}
	//shouldn't get here
	return reconcile.Result{}, nil
}

func reconcileRelease(releaseName string, release *remeshv1alpha1.Release, virtualApp *remeshv1alpha1.VirtualApp, deleted bool) {
	releaseFlows := virtualApp.Spec.ReleaseFlows
	rfKey := -1
	for i := range releaseFlows {
		if releaseFlows[i].ReleaseName == releaseName {
			rfKey = i
		}
	}

	if rfKey == -1 {
		var targeting *map[string]*remeshv1alpha1.SegmentSpec
		if release.Spec.Targeting != nil {
			targeting = &map[string]*remeshv1alpha1.SegmentSpec{}
			for _, seg := range release.Spec.Targeting.Segments {
				(*targeting)[seg] = nil
			}
		}

		virtualApp.Spec.ReleaseFlows = append(virtualApp.Spec.ReleaseFlows, remeshv1alpha1.ReleaseFlow{
			ReleaseName: releaseName,
			Release:     *release.Spec.DeepCopy(),
			Targeting:   targeting,
			LayoutName:  release.Spec.Layout,
			Layout:      nil,
		})
	} else {
		if deleted {
			virtualApp.Spec.ReleaseFlows = append(virtualApp.Spec.ReleaseFlows[:rfKey], virtualApp.Spec.ReleaseFlows[rfKey+1:]...)
		} else {
			releaseFlow := releaseFlows[rfKey]
			releaseFlow.Release = *release.Spec.DeepCopy()
			if release.Spec.Layout != releaseFlow.LayoutName {
				releaseFlow.LayoutName = release.Spec.Layout
				releaseFlow.Layout = nil
			}
			if release.Spec.Targeting == nil {
				releaseFlow.Targeting = nil
			} else {
				newTargeting := &map[string]*remeshv1alpha1.SegmentSpec{}
				for _, seg := range release.Spec.Targeting.Segments {
					if releaseFlow.Targeting != nil {
						if val, ok := (*releaseFlow.Targeting)[seg]; ok {
							(*newTargeting)[seg] = val
						} else {
							(*newTargeting)[seg] = nil
						}
					} else {
						(*newTargeting)[seg] = nil
					}
				}
				releaseFlow.Targeting = newTargeting
			}

		}
	}
	sort.Sort(remeshv1alpha1.ByPriority(virtualApp.Spec.ReleaseFlows))
}
