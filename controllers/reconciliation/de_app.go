package reconciliation

import (
	"context"
	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *ReconcileManager) reconcileDeApp(ctx context.Context, image string, req ctrl.Request, cr *v1alpha1.EntandoAppV2) error {
	r.Log.Info("Starting DeApp reconciliation flow")
	time.Sleep(time.Second * 5)
	r.Log.Info("Finished DeApp reconciliation flow")

	return nil
}
