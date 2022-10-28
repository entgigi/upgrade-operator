package reconciliation

import (
	"context"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *ReconcileManager) reconcileAppBuilder(ctx context.Context, image string, req ctrl.Request) error {
	time.Sleep(2 * time.Second)
	return nil
}
