package reconciliation

import (
	"context"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const componentManagerDeploymentEndName = "cm-deployment"

func (r *ReconcileManager) reconcileComponentManager(ctx context.Context, image string, req ctrl.Request) error {
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

	if deployment != nil {
		deployment.Spec.Template.Spec.Containers[0].Image = image
		if err := r.Update(ctx, deployment); err != nil {
			return err
		}
	}

	return nil
}
