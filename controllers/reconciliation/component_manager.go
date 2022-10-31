package reconciliation

import (
	"context"
	"fmt"
	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/utils"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const componentManagerDeploymentEndName = "cm-deployment"

func (r *ReconcileManager) reconcileComponentManager(ctx context.Context, image string, req ctrl.Request) error {
	r.Log.Info("Starting ComponentManager reconciliation flow")

	deploymentList := &appsv1.DeploymentList{}

	if err := r.Client.List(ctx, deploymentList); err != nil {
		return err
	}

	var deployment *appsv1.Deployment
	for _, item := range deploymentList.Items {
		if strings.HasSuffix(item.Name, componentManagerDeploymentEndName) {
			deployment = &item
			break
		}
	}

	if deployment == nil {
		return fmt.Errorf("deployment ComponentManager not found")
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
