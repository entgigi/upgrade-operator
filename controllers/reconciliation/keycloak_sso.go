package reconciliation

import (
	"context"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *ReconcileManager) reconcileKeycloak(ctx context.Context, image string, req ctrl.Request) error {
	r.log.Info("Starting keycloak reconciliation flow")

	r.log.Info("Finished keycloak reconciliation flow")

	return nil
}
