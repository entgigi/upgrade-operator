package reconciliation

import (
	"context"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
)

const appBuilderKubeId = "ab"

func (r *ReconcileManager) reconcileAppBuilder(ctx context.Context, image string, req ctrl.Request, cr *v1alpha1.EntandoAppV2) error {
	r.Log.Info("Starting AppBuilder reconciliation flow")

	deployment, err := r.mustGetDeployment(ctx, cr.Spec.EntandoAppName, appBuilderKubeId)
	if err != nil {
		return err
	}

	deployment = r.updateCommonDeploymentData(deployment,
		image,
		r.envVarByVersion(cr, appBuilderManagerEnv),
		cr.Spec.CommonEnvironmentVariables,
		cr.Spec.AppBuilder.EnvironmentVariables)

	if err = r.Update(ctx, deployment); err != nil {
		return err
	}

	// FIXME remove
	time.Sleep(3 * time.Second)

	r.Log.Info("Finished AppBuilder reconciliation flow")
	return nil
}

var appBuilderManagerEnv = listApplicationEnvVar{
	"7.1.1": applicationEnvVar{},
}
