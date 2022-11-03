package reconciliation

import (
	"context"
	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	keycloakDeploymentName = "default-sso-in-namespace-deployment"
)

func (r *ReconcileManager) reconcileKeycloak(ctx context.Context, image string, req ctrl.Request, cr *v1alpha1.EntandoAppV2) error {
	r.Log.Info("Starting Keycloak reconciliation flow")
	time.Sleep(time.Second * 5)
	r.Log.Info("Finished Keycloak reconciliation flow")

	return nil
}
