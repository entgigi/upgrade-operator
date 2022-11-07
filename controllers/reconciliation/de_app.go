package reconciliation

import (
	"context"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
)

func (r *ReconcileManager) reconcileDeApp(ctx context.Context, image string, req ctrl.Request, cr *v1alpha1.EntandoAppV2) error {
	r.Log.Info("Starting DeApp reconciliation flow")

	deployment, err := r.mustGetDeployment(ctx, cr.Spec.EntandoAppName, "")
	if err != nil {
		return err
	}

	deployment = r.updateCommonDeploymentData(deployment,
		image,
		r.envVarByVersion(cr, deAppManagerEnv),
		cr.Spec.CommonEnvironmentVariables,
		cr.Spec.DeApp.EnvironmentVariables)

	if err = r.Update(ctx, deployment); err != nil {
		return err
	}

	// FIXME remove
	time.Sleep(3 * time.Second)

	r.Log.Info("Finished DeApp reconciliation flow")

	return nil
}

var deAppManagerEnv = mapApplicationEnvVar{
	"7.1.1": applicationEnvVar{},
}
