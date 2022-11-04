package reconciliation

import (
	"context"
	"errors"
	"fmt"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/utils"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
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

	// if external and deployment present internally => delete it
	if cr.Spec.Keycloak.ExternalService && len(deployments.Items) > 0 {
		r.Log.Info("external Keycloak requirement but deployment found. Dismissing...")
		return r.Delete(ctx, &deployments.Items[0])
	}

	// if internal
	if !cr.Spec.Keycloak.ExternalService {

		// and no deployment
		if len(deployments.Items) <= 0 || len(deployments.Items) > 1 {
			return fmt.Errorf("internal Keycloak requirement but %d found", len(deployments.Items))
			// TODO manage the case with 0 deployment by deploy a new keycloak in the future?
		}
		// otherwise reconcile
		deployment := &deployments.Items[0]

		deployment = r.updateCommonDeploymentData(deployment,
			image,
			r.envVarByVersion(cr, keycloakManagerEnv),
			cr.Spec.CommonEnvironmentVariables,
			cr.Spec.Keycloak.EnvironmentVariables)

		if err := r.Update(ctx, deployment); err != nil {
			return err
		}

		// FIXME remove
		time.Sleep(3 * time.Second)

		r.Log.Info("Finished Keycloak reconciliation flow")
	}

	return nil
}

var keycloakManagerEnv = listApplicationEnvVar{
	"7.1.1": applicationEnvVar{},
}
