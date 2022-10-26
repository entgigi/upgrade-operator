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

	upgradev1alpha1 "github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EntandoUpgradeReconciler reconciles a EntandoUpgrade object
type EntandoUpgradeReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=upgrade.entando.org,resources=entandoupgrades,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=upgrade.entando.org,resources=entandoupgrades/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=upgrade.entando.org,resources=entandoupgrades/finalizers,verbs=update

func (r *EntandoUpgradeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithName("Upgrade Controller")
	log.Info("Start reconciling EntandoUpgrade custom resources")

	entandoUpgrade := upgradev1alpha1.EntandoUpgrade{}
	err := r.Client.Get(ctx, req.NamespacedName, &entandoUpgrade)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// FIXME! add start proregss in status EntandoUpgrade cr
	entandoUpgrade.Status.Progress = "starting updaye"
	r.updateProgressStatus(ctx, entandoUpgrade)
	// FIXME! add sleep 1 minutes

	// FIXME! add finished proregss in status EntandoUpgrade cr
	entandoUpgrade.Status.Progress = "starting updaye"
	r.updateProgressStatus(ctx, entandoUpgrade)

	log.Info("Reconciled EntandoUpgrade custom resources")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EntandoUpgradeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	//log := r.Log.WithName("Upgrade Controller")
	return ctrl.NewControllerManagedBy(mgr).
		// FIXME! add filter on create for EntandoUpgrade cr
		For(&upgradev1alpha1.EntandoUpgrade{}).
		Complete(r)
}

func (r *EntandoUpgradeReconciler) updateProgressStatus(ctx context.Context, cr upgradev1alpha1.EntandoUpgrade) {
	log := r.Log.WithName("Upgrade Controller")
	err := r.Status().Update(ctx, &cr)
	if err != nil {
		log.Error(err, "Unable to update EntandoUpgrade's progress status", "progress", cr.Status.Progress)
	}

}
