package reconciliation

import (
	"context"
	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/utils"
	v1 "k8s.io/api/core/v1"
	errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
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
		r.envVarByVersion(ctx, req, cr, componentManagerEnv),
		cr.Spec.CommonEnvironmentVariables,
		cr.Spec.ComponentManager.EnvironmentVariables)

	deployment = utils.ManageUpdateStrategy(deployment, cr)

	if err = r.Update(ctx, deployment); err != nil {
		return err
	}

	// FIXME remove
	time.Sleep(3 * time.Second)

	r.Log.Info("Finished ComponentManager reconciliation flow")

	return nil
}

var componentManagerEnv = mapApplicationEnvVar{
	"7.1.1": applicationEnvVar{
		"ENTANDO_APP_HOST_NAME": calculateAppHostName,
		"ENTANDO_APP_USE_TLS":   calculateAppTls,
	},
}

func calculateAppHostName(ctx context.Context, req ctrl.Request, r *ReconcileManager, cr *v1alpha1.EntandoAppV2) string {
	return cr.Spec.IngressHostName
}

func calculateAppTls(ctx context.Context, req ctrl.Request, r *ReconcileManager, cr *v1alpha1.EntandoAppV2) string {
	// retrieve the configmap entando-operator-config of the legacy operator
	// if configmap not found return false
	// if error log error and return false
	// if configmap contains key entando.tls.secret.name not empty then tls returns true

	configmap := &v1.ConfigMap{}

	err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      "entando-operator-config",
	}, configmap)

	if err != nil {
		if !errors.IsNotFound(err) {
			r.Log.Error(err, "error retrieving configmap entando-operator-config")
		}
		return "false"
	}

	if val, found := configmap.Data["entando.tls.secret.name"]; found && len(val) > 0 {
		return "true"
	}

	return "false"
}
