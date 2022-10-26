/*
Copyright 2022.

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

	v1alpha1 "github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EntandoAppV2Reconciler reconciles a EntandoAppV2 object
type EntandoAppV2Reconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.entando.org,resources=entandoappv2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.entando.org,resources=entandoappv2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.entando.org,resources=entandoappv2s/finalizers,verbs=update

func (r *EntandoAppV2Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithName("Upgrade Controller")
	log.Info("Start reconciling EntandoAppV2 custom resources")

	EntandoAppV2 := v1alpha1.EntandoAppV2{}
	err := r.Client.Get(ctx, req.NamespacedName, &EntandoAppV2)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// FIXME! add start proregss in status EntandoAppV2 cr
	EntandoAppV2.Status.Progress = "starting updaye"
	r.updateProgressStatus(ctx, EntandoAppV2)
	// FIXME! add sleep 1 minutes

	// FIXME! add finished proregss in status EntandoAppV2 cr
	EntandoAppV2.Status.Progress = "starting updaye"
	r.updateProgressStatus(ctx, EntandoAppV2)

	log.Info("Reconciled EntandoAppV2 custom resources")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EntandoAppV2Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	//log := r.Log.WithName("Upgrade Controller")
	return ctrl.NewControllerManagedBy(mgr).
		// FIXME! add filter on create for EntandoAppV2 cr
		For(&v1alpha1.EntandoAppV2{}).
		Complete(r)
}

func (r *EntandoAppV2Reconciler) updateProgressStatus(ctx context.Context, cr v1alpha1.EntandoAppV2) {
	log := r.Log.WithName("Upgrade Controller")
	err := r.Status().Update(ctx, &cr)
	if err != nil {
		log.Error(err, "Unable to update EntandoAppV2's progress status", "progress", cr.Status.Progress)
	}

}
