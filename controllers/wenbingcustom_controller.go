/*
Copyright 2021.

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

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	mydomainv1alpha1 "github.com/wenbingz/wenbing-operator/api/v1alpha1"
)

// WenbingCustomReconciler reconciles a WenbingCustom object
type WenbingCustomReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=example.wenbing.org.wenbing,resources=wenbingcustoms,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.wenbing.org.wenbing,resources=wenbingcustoms/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.wenbing.org.wenbing,resources=wenbingcustoms/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WenbingCustom object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *WenbingCustomReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("WenbingCustom", req.NamespacedName)

	// your logic here
	instance := &mydomainv1alpha1.WenbingCustom{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, nil
	}

	// Check if this Deployment already exists
	found := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{
			Name: instance.GetName(),
			Namespace: instance.GetNamespace(),
		},
		found)
	var result *reconcile.Result
	result, err = r.ensureDeployment(req, instance, r.backendDeployment(instance))
	if result != nil {
		log.Error(err, "Deployment not ready!")
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.backendService(instance))
	if result != nil {
		log.Error(err, "Service not ready!")
		return *result, err
	}
	log.Info("Skip reconcile: Deployment and service already exists",
		"Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WenbingCustomReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mydomainv1alpha1.WenbingCustom{}).
		Complete(r)
}
