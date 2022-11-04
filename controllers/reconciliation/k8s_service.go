package reconciliation

import (
	"context"
	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
)

const (
	K8sServiceDeploymentName = "entando-k8s-service"
)

func (r *ReconcileManager) reconcileK8sService(ctx context.Context, req ctrl.Request, image string, cr *v1alpha1.EntandoAppV2) error {
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
		r.envVarByVersion(cr, k8sServiceManagerEnv),
		cr.Spec.CommonEnvironmentVariables,
		cr.Spec.K8sService.EnvironmentVariables)

	if err := r.Update(ctx, deployment); err != nil {
		return err
	}

	// FIXME remove
	time.Sleep(3 * time.Second)

	r.Log.Info("Finished k8s-service reconciliation flow")

	return nil
}

var k8sServiceManagerEnv = listApplicationEnvVar{
	"7.1.1": applicationEnvVar{},
}
