package reconciliation

import (
	"context"
	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
)

const componentManagerKubeId = "cm"

func (r *ReconcileManager) reconcileComponentManager(ctx context.Context, image string, req ctrl.Request, cr *v1alpha1.EntandoAppV2) error {
	r.Log.Info("Starting ComponentManager reconciliation flow")

	deployment, err := r.mustGetDeployment(ctx, cr.Spec.EntandoAppName, componentManagerKubeId)
	if err != nil {
		return err
	}

	deployment = r.updateCommonDeploymentData(deployment,
		image,
		cr.Spec.CommonEnvironmentVariables,
		cr.Spec.ComponentManager.EnvironmentVariables)

	if err = r.Update(ctx, deployment); err != nil {
		return err
	}

	// FIXME remove
	time.Sleep(3 * time.Second)

	r.Log.Info("Finished ComponentManager reconciliation flow")

	return nil
}
