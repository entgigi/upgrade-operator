package reconciliation

import (
	"context"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *ReconcileManager) reconcileDeApp(ctx context.Context, image string, req ctrl.Request) error {
	r.Log.Info("Starting DeApp reconciliation flow")
	time.Sleep(time.Second * 5)
	r.Log.Info("Finished DeApp reconciliation flow")

	return nil
}
