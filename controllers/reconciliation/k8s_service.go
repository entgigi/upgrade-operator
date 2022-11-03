package reconciliation

import (
	"context"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	K8sServiceDeploymentName = "entando-k8s-service"
)

func (r *ReconcileManager) reconcileK8sService(ctx context.Context, req ctrl.Request, image string, entandoAppV2 v1alpha1.EntandoAppV2) error {
	r.Log.Info("Starting k8s-service reconciliation flow")

	deployment := &appsv1.Deployment{}

	identifier := types.NamespacedName{
		Name:      K8sServiceDeploymentName,
		Namespace: req.NamespacedName.Namespace}
	if err := r.Client.Get(ctx, identifier, deployment); err != nil {
		return err
	}

	deployment = r.updateCommonDeploymentData(deployment,
		image,
		entandoAppV2.Spec.CommonEnvironmentVariables,
		entandoAppV2.Spec.ComponentManager.EnvironmentVariables)

	if err := r.Update(ctx, deployment); err != nil {
		return err
	}

	r.Log.Info("Finished k8s-service reconciliation flow")

	return nil
}
