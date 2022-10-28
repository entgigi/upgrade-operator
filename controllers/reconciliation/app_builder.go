package reconciliation

import (
	"context"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *ReconcileManager) reconcileAppBuilder(ctx context.Context, image string, req ctrl.Request) error {
	r.Log.Info("Starting AppBuilder reconciliation flow")
	time.Sleep(time.Second * 5)
	r.Log.Info("Finished AppBuilder reconciliation flow")
	return nil
}
