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
	"github.com/entgigi/upgrade-operator.git/controllers/reconciliation"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	entandoAppFinalizer = "app.entando.org/finalizer"
	controllerLogName   = "EntandoAppV2 Controller"
)

// EntandoAppV2Reconciler reconciles a EntandoAppV2 object
type EntandoAppV2Reconciler struct {
	// TODO centralize log variable into one single struct to embed
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.entando.org,resources=entandoappv2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.entando.org,resources=entandoappv2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.entando.org,resources=entandoappv2s/finalizers,verbs=update
//+kubebuilder:rbac:groups="*",resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="operators.coreos.com",resources=clusterserviceversions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="*",resources=configmaps,verbs=get;list;watch

func (r *EntandoAppV2Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithName(controllerLogName)
	log.Info("Start reconciling EntandoAppV2 custom resources")

	entandoAppV2 := v1alpha1.EntandoAppV2{}
	err := r.Client.Get(ctx, req.NamespacedName, &entandoAppV2)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if the EntandoApp instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isEntandoAppV2MarkedToBeDeleted := entandoAppV2.GetDeletionTimestamp() != nil
	if isEntandoAppV2MarkedToBeDeleted {
		if err := r.removeFinalizer(ctx, &entandoAppV2, log); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer for this CR
	err = r.addFinalizer(ctx, &entandoAppV2)
	if err != nil {
		return ctrl.Result{}, err
	}

	manager := reconciliation.NewReconcileManager(r.Client, r.Log)
	if err = manager.MainReconcile(ctx, req); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Reconciled EntandoAppV2 custom resources")
	return ctrl.Result{}, nil

}

// =====================================================================
// SetupWithManager sets up the controller with the Manager.
// =====================================================================
func (r *EntandoAppV2Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	//log := r.Log.WithName("Upgrade Controller")
	return ctrl.NewControllerManagedBy(mgr).
		// FIXME! add filter on create for EntandoAppV2 cr
		For(&v1alpha1.EntandoAppV2{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}). //solo modifiche a spec
		Complete(r)
}

// =====================================================================
// Add the cleanup steps that the operator
// needs to do before the CR can be deleted. Examples
// of finalizers include performing backups and deleting
// resources that are not owned by this CR, like a PVC.
// =====================================================================
func (r *EntandoAppV2Reconciler) finalizeEntandoApp(log logr.Logger, m *v1alpha1.EntandoAppV2) error {
	log.Info("Successfully finalized entandoApp")
	return nil
}

func (r *EntandoAppV2Reconciler) addFinalizer(ctx context.Context, entandoAppV2 *v1alpha1.EntandoAppV2) error {
	if !controllerutil.ContainsFinalizer(entandoAppV2, entandoAppFinalizer) {
		controllerutil.AddFinalizer(entandoAppV2, entandoAppFinalizer)
		return r.Update(ctx, entandoAppV2)
	}
	return nil
}

func (r *EntandoAppV2Reconciler) removeFinalizer(ctx context.Context, entandoAppV2 *v1alpha1.EntandoAppV2, log logr.Logger) error {
	if controllerutil.ContainsFinalizer(entandoAppV2, entandoAppFinalizer) {
		// Run finalization logic for entandoAppFinalizer. If the
		// finalization logic fails, don't remove the finalizer so
		// that we can retry during the next reconciliation.
		if err := r.finalizeEntandoApp(log, entandoAppV2); err != nil {
			return err
		}

		// Remove entandoAppFinalizer. Once all finalizers have been
		// removed, the object will be deleted.
		controllerutil.RemoveFinalizer(entandoAppV2, entandoAppFinalizer)
		err := r.Update(ctx, entandoAppV2)
		if err != nil {
			return err
		}
	}
	return nil
}
