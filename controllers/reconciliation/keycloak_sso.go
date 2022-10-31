package reconciliation

import (
	"context"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	keycloakDeploymentName = "default-sso-in-namespace-deployment"
)

func (r *ReconcileManager) reconcileKeycloak(ctx context.Context, image string, req ctrl.Request) error {
	r.Log.Info("Starting keycloak reconciliation flow")
	time.Sleep(time.Second * 5)
	r.Log.Info("Finished keycloak reconciliation flow")

	return nil
}
