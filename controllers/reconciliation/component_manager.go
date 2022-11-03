package reconciliation

import (
	"context"
	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/utils"
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

	entandoAppV2 := v1alpha1.EntandoAppV2{}
	if err := r.Client.Get(ctx, req.NamespacedName, &entandoAppV2); err != nil {
		return err
	}

	envVars := utils.MergeEnvVars(entandoAppV2, deployment)
	deployment.Spec.Template.Spec.Containers[0].Env = envVars

	if err := r.Update(ctx, deployment); err != nil {
		return err
	}

	r.Log.Info("Finished ComponentManager reconciliation flow")

	return nil
}
