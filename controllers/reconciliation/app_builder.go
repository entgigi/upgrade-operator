package reconciliation

import (
	"context"
	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *ReconcileManager) reconcileAppBuilder(ctx context.Context, image string, req ctrl.Request, cr *v1alpha1.EntandoAppV2) error {
	r.Log.Info("Starting AppBuilder reconciliation flow")
	time.Sleep(time.Second * 5)
	r.Log.Info("Finished AppBuilder reconciliation flow")
	return nil
}
