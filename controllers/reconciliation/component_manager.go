package reconciliation

import (
	"context"
	"fmt"
	"strings"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/utils"
<<<<<<< HEAD
=======
	appsv1 "k8s.io/api/apps/v1"
>>>>>>> 48b1b37 (revert disabled reconciler)
	ctrl "sigs.k8s.io/controller-runtime"
)

const componentManagerKubeId = "cm"

func (r *ReconcileManager) reconcileComponentManager(ctx context.Context, image string, req ctrl.Request, cr *v1alpha1.EntandoAppV2) error {
	r.Log.Info("Starting ComponentManager reconciliation flow")

	labels := utils.BuildEntandoDeploymentLabelSelector(cr.Spec.EntandoAppName, componentManagerKubeId)
	deployment, err := utils.MustGetFirstDeploymentByLabels(ctx, r.Client, labels)
	if err != nil {
		return err
	}

	deployment.Spec.Template.Spec.Containers[0].Image = image


	envVars := utils.MergeEnvVars(cr, deployment)
	deployment.Spec.Template.Spec.Containers[0].Env = envVars

	if err := r.Update(ctx, deployment); err != nil {
		return err
	}

	r.Log.Info("Finished ComponentManager reconciliation flow")

	return nil
}
