package reconciliation

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	keycloakDeploymentName = "default-sso-in-namespace-deployment"
)

func (r *ReconcileManager) reconcileKeycloak(ctx context.Context, image string, req ctrl.Request) error {
	r.Log.Info("Starting keycloak reconciliation flow")

	//time.Sleep(2 * time.Second)

	deployment, err := r.getDeployment(ctx, req.Namespace, keycloakDeploymentName)
	if err != nil {
		// TODO when we will manage the installation => check if errors.IsNotFound(err) and create the deployment
		return err
	}

	deployment.Spec.Template.Spec.Containers[0].Image = image

	// TODO merge env vars

	err = r.Update(ctx, &deployment)

	r.Log.Info("Finished keycloak reconciliation flow")

	return err
}
