package reconciliation

import (
	"context"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/common"
	"github.com/entgigi/upgrade-operator.git/legacy"
	"github.com/entgigi/upgrade-operator.git/service"
	"github.com/entgigi/upgrade-operator.git/utils"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	reconciliationLogName = "ReconcileManager"
	numberOfSteps         = 4
)

type ReconcileManager struct {
	common.BaseK8sStructure
	statusUpdater *StatusUpdater
}

type ReconcileComponentFunc func(ctx context.Context, image string, req ctrl.Request, cr *v1alpha1.EntandoAppV2) error

// NewReconcileManager initialize a ReconcileManager
func NewReconcileManager(client client.Client, log logr.Logger) *ReconcileManager {
	logger := log.WithName(reconciliationLogName)
	statusUpdater := NewStatusUpdater(client, logger)
	return &ReconcileManager{
		common.BaseK8sStructure{Client: client, Log: logger},
		statusUpdater,
	}
}

func (r *ReconcileManager) MainReconcile(ctx context.Context, req ctrl.Request) error {

	r.Log.Info("Starting main reconciliation flow")
	r.statusUpdater.SetReconcileStarted(ctx, req.NamespacedName, numberOfSteps)

	var err error
	crReadOnly := &v1alpha1.EntandoAppV2{}
	if err := r.Client.Get(ctx, req.NamespacedName, crReadOnly); err != nil {
		return err
	}

	imageManager := service.NewImageManager(r.Log)
	var images service.EntandoAppImages

	if images, err = imageManager.FetchAndComposeImagesMap(crReadOnly); err != nil {
		return err
	}
	//r.Log.Info(fmt.Sprintf("%+v\n", images))

	//TODO reconcile secrets for ca before EntandoApp components

	if _, err = r.reconcileComponent(ctx, req, "Keycloak", r.reconcileKeycloak, images.FetchKeycloak(), crReadOnly); err != nil {
		return err
	}

	if _, err = r.reconcileComponent(ctx, req, "DeApp", r.reconcileDeApp, images.FetchDeApp(), crReadOnly); err != nil {
		return err
	}

	if _, err = r.reconcileComponent(ctx, req, "AppBuilder", r.reconcileAppBuilder, images.FetchAppBuilder(), crReadOnly); err != nil {
		return err
	}

	// TODO before check entando-k8s-service app ready
	cr := &v1alpha1.EntandoAppV2{}
	if cr, err = r.reconcileComponent(ctx, req, "ComponentManager", r.reconcileComponentManager, images.FetchComponentManager(), crReadOnly); err != nil {
		return err
	}

	// =========================== start legacy section ===========================
	// progress step not added because is not a business step but jsut technical
	legacyReconcile := legacy.NewLegacyReconcileManager(r.Client, r.Log)
	if utils.IsOlmInstallation() {
		r.statusUpdater.SetReconcileProcessingComponent(ctx, req.NamespacedName, "Csv")

		if err = legacyReconcile.ReconcileClusterServiceVersion(ctx, req, images); err != nil {
			r.statusUpdater.SetReconcileFailed(ctx, req.NamespacedName, "CsvReconciliationFailed")
			return err
		}
	} else {
		r.statusUpdater.SetReconcileProcessingComponent(ctx, req.NamespacedName, "ImageInfo")
		if err = legacyReconcile.ReconcileImageInfo(ctx, req, images); err != nil {
			r.statusUpdater.SetReconcileFailed(ctx, req.NamespacedName, "ImageInfoReconciliationFailed")
			return err
		}
	}
	// =========================== end legacy section =============================

	// reconcile k8s-service component taking into account of the legacy Operator behavior,
	// in non-olm install we need to reconcile k8s-service deployment
	if !utils.IsOlmInstallation() {
		// K8sService
		r.statusUpdater.SetReconcileProcessingComponent(ctx, req.NamespacedName, "K8sService")

		// TODO decide if add the k8service in the progress count. in that case we could also consider to adapt the k8s-service reconciliation function to the standard format

		if err = r.reconcileK8sService(ctx, req, images.FetchK8sService(), crReadOnly); err != nil {
			r.statusUpdater.SetReconcileFailed(ctx, req.NamespacedName, "K8sServiceReconciliationFailed")
			return err
		}
		// legacy K8sCoordinator restart ? no needs

	}

	// Check for progress/total mismatch
	if cr.Status.Progress != numberOfSteps {
		r.Log.Info("WARNING: progress different from total at the end of reconciliation", "progress", cr.Status.Progress, "total", numberOfSteps)
	}

	r.statusUpdater.SetReconcileSuccessfullyCompleted(ctx, req.NamespacedName)

	return nil
}

// reconcileComponent pattern function to reconcile a single component
func (r *ReconcileManager) reconcileComponent(ctx context.Context,
	req ctrl.Request,
	componentName string,
	reconcile ReconcileComponentFunc,
	imageUrl string,
	cr *v1alpha1.EntandoAppV2) (*v1alpha1.EntandoAppV2, error) {

	r.statusUpdater.SetReconcileProcessingComponent(ctx, req.NamespacedName, componentName)
	if err := reconcile(ctx, imageUrl, req, cr); err != nil {
		r.statusUpdater.SetReconcileFailed(ctx, req.NamespacedName, componentName+"ReconciliationFailed")
		return nil, err
	}
	return r.statusUpdater.IncrementProgress(ctx, req.NamespacedName)
}

// mustGetDeployment try to get the first deployment starting by the entando label dedicated to identify the component
// return the found Deployment if only 1 Deployment corresponds to the search criteria. return err otherwise
func (r *ReconcileManager) mustGetDeployment(ctx context.Context, entandoAppName string, kubeComponentId string) (*v1.Deployment, error) {

	labels := utils.BuildDeploymentLabelSelectorWithAppName(entandoAppName, kubeComponentId)
	return utils.MustGetFirstDeploymentByLabels(ctx, r.Client, labels)
}

// updateCommonDeploymentData update the received deployment applying the changes common to every entando deployment
// return the updated deployment
// it does not update the deployment present in Kubernetes
func (r *ReconcileManager) updateCommonDeploymentData(deployment *v1.Deployment,
	image string,
	componentVersionEnvVars []corev1.EnvVar,
	genericEnvVars []corev1.EnvVar,
	specificEnvVars []corev1.EnvVar) *v1.Deployment {

	deployment = r.updateImage(deployment, image)
	deployment = r.applyVersionUpgradeEnvVars(deployment, componentVersionEnvVars)
	deployment = r.mergeEnvVars(deployment, genericEnvVars, specificEnvVars)

	return deployment
}

// updateImage update the Spec.Template.Spec.Containers[0].Image property of the received deployment with the one
// contained in the image param
// return the updated deployment
func (r *ReconcileManager) updateImage(deployment *v1.Deployment, image string) *v1.Deployment {

	if image == "" {
		r.Log.Info("no new image found for the deployment " + deployment.Name)
	} else {
		deployment.Spec.Template.Spec.Containers[0].Image = image
	}

	return deployment
}

// mergeEnvVars merge the received environment variables
// the order is:
//   - deployment env vars
//   - overridden by genericEnvVars
//   - overridden by specificEnvVars
//
// return the updated deployment
func (r *ReconcileManager) mergeEnvVars(deployment *v1.Deployment,
	genericEnvVars []corev1.EnvVar,
	specificEnvVars []corev1.EnvVar) *v1.Deployment {

	envVars := utils.MergeEnvVars(deployment, genericEnvVars, specificEnvVars)
	deployment.Spec.Template.Spec.Containers[0].Env = envVars
	return deployment
}

func (r *ReconcileManager) applyVersionUpgradeEnvVars(deployment *v1.Deployment, newEnvs []corev1.EnvVar) *v1.Deployment {

	actualEnvMap := utils.ConvertEnvVarSliceToMap(deployment.Spec.Template.Spec.Containers[0].Env)
	for _, env := range newEnvs {
		if _, found := actualEnvMap[env.Name]; !found {
			deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, env)
		}
	}

	return deployment
}

type applicationEnvVar map[string]func(r *ReconcileManager, cr *v1alpha1.EntandoAppV2) string

type mapApplicationEnvVar map[string]applicationEnvVar

func (r *ReconcileManager) envVarByVersion(cr *v1alpha1.EntandoAppV2, upgradeVersionEnvMap mapApplicationEnvVar) []corev1.EnvVar {
	newAppEnvs := make([]corev1.EnvVar, 0)
	if envs, ok := upgradeVersionEnvMap[cr.Spec.Version]; ok {
		for key, value := range envs {
			newAppEnvs = append(newAppEnvs, corev1.EnvVar{Name: key, Value: value(r, cr)})
		}
	}
	return newAppEnvs
}
