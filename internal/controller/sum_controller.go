/*
Copyright 2024.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	calcv1 "github.com/sudipto-003/kcal/api/v1"
)

// SumReconciler reconciles a Sum object
type SumReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=calc.toad.d3s.io,resources=sums,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=calc.toad.d3s.io,resources=sums/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=calc.toad.d3s.io,resources=sums/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Sum object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *SumReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	sum := &calcv1.Sum{}
	err := r.Get(ctx, req.NamespacedName, sum)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "Unable to fetch resource. Name: ", req.NamespacedName.Name, "Namespace: ", req.NamespacedName.Namespace)
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	arg1, arg2 := sum.Spec.Argument1, sum.Spec.Argument2
	res1 := new(int64)
	*res1 = int64(*arg1) + int64(*arg2)
	if *res1 < 0 {
		err_msg := "The sum up value must be positive integer, greater than 0."
		sum.Status.Error = &err_msg
		*res1 = 1
	} else {
		sum.Status.Error = nil
	}

	sum.Status.Result = res1
	if err := r.Status().Update(ctx, sum); err != nil {
		logger.Error(err, "Failed to update the resource. Name: ", req.NamespacedName.Name, "Namespace: ", req.NamespacedName.Namespace)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SumReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&calcv1.Sum{}).
		Complete(r)
}
