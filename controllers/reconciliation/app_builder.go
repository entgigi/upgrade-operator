package reconciliation

import (
	"context"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *ReconcileManager) reconcileAppBuilder(ctx context.Context, image string, req ctrl.Request) error {
	return nil
}
