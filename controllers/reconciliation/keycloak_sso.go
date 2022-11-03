package reconciliation

import (
	"context"
	"errors"
	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/utils"
	ctrl "sigs.k8s.io/controller-runtime"
)

const keycloakKubeId = "default-sso-in-namespace"

func (r *ReconcileManager) reconcileKeycloak(ctx context.Context, image string, req ctrl.Request, cr *v1alpha1.EntandoAppV2) error {
	r.Log.Info("Starting Keycloak reconciliation flow")

	labels := utils.BuildDeploymentLabelSelector(keycloakKubeId)
	deployments, err := utils.FindDeploymentsByLabels(ctx, r.Client, labels)
	if err != nil {
		return err
	}

	if len(deployments.Items) > 1 {
		return errors.New("more than 1 keycloak deployment found matching the received labels")
	}

	if len(deployments.Items) <= 0 {
		r.Log.Info("Keycloak deployment not found. Assuming external Keycloak deployment, skipping reconciliation")
		return nil
	}

	deployment := &deployments.Items[0]

	deployment = r.updateCommonDeploymentData(deployment,
		image,
		cr.Spec.CommonEnvironmentVariables,
		cr.Spec.Keycloak.EnvironmentVariables)

	if err := r.Update(ctx, deployment); err != nil {
		return err
	}

	r.Log.Info("Finished Keycloak reconciliation flow")

	return nil
}
