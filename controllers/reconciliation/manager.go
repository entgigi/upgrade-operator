package reconciliation

import (
	"context"
	"strings"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/common"
	"github.com/go-logr/logr"
	csv "github.com/operator-framework/api/pkg/operators/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	reconciliationLogName = "ReconcileManager"
	numberOfSteps         = 4
)

type ReconcileManager struct {
	client.Client
	log           logr.Logger
	scheme        *runtime.Scheme
	statusUpdater *StatusUpdater
}

// NewReconcileManager initialize a ReconcileManager
func NewReconcileManager(client client.Client, log logr.Logger, scheme *runtime.Scheme) *ReconcileManager {
	logger := log.WithName(reconciliationLogName)
	statusUpdater := NewStatusUpdater(client, logger)
	return &ReconcileManager{
		client,
		logger,
		scheme,
		statusUpdater,
	}
}

func (r *ReconcileManager) MainReconcile(ctx context.Context, req ctrl.Request) error {

	r.log.Info("Starting main reconciliation flow")
	r.statusUpdater.SetReconcileStarted(ctx, req.NamespacedName, numberOfSteps)

	cr := &v1alpha1.EntandoAppV2{}
	if err := r.Client.Get(ctx, req.NamespacedName, cr); err != nil {
		return err
	}
	images := r.fetchImages(*cr)

	r.statusUpdater.SetReconcileProcessingComponent(ctx, req.NamespacedName, "Keycloak")
	if err := r.reconcileKeycloak(ctx, images.FetchKeycloak(), req); err != nil {
		r.statusUpdater.SetReconcileFailed(ctx, req.NamespacedName, "KeycloakReconciliationFailed")
		return err
	}
	r.statusUpdater.IncrementProgress(ctx, req.NamespacedName)

	r.statusUpdater.SetReconcileProcessingComponent(ctx, req.NamespacedName, "DeApp")
	if err := r.reconcileDeApp(ctx, images.FetchDeApp(), req); err != nil {
		r.statusUpdater.SetReconcileFailed(ctx, req.NamespacedName, "DeAppReconciliationFailed")
		return err
	}
	r.statusUpdater.IncrementProgress(ctx, req.NamespacedName)

	r.statusUpdater.SetReconcileProcessingComponent(ctx, req.NamespacedName, "AppBuilder")
	if err := r.reconcileAppBuilder(ctx, images.FetchAppBuilder(), req); err != nil {
		r.statusUpdater.SetReconcileFailed(ctx, req.NamespacedName, "AppBuilderReconciliationFailed")
		return err
	}
	r.statusUpdater.IncrementProgress(ctx, req.NamespacedName)

	// TODO before check entando-k8s-service app ready
	r.statusUpdater.SetReconcileProcessingComponent(ctx, req.NamespacedName, "ComponentManager")
	if err := r.reconcileComponentManager(ctx, images.FetchComponentManager(), req); err != nil {
		r.statusUpdater.SetReconcileFailed(ctx, req.NamespacedName, "ComponentManagerReconciliationFailed")
		return err
	}
	cr, _ = r.statusUpdater.IncrementProgress(ctx, req.NamespacedName)

	// Check for progress/total mismatch
	if cr.Status.Progress != numberOfSteps {
		r.log.Info("WARNING: progress different from total at the end of reconciliation", "progress", cr.Status.Progress, "total", numberOfSteps)
	}

	// TODO manage CSV

	r.statusUpdater.SetReconcileSuccessfullyCompleted(ctx, req.NamespacedName)
	return nil
}

// fetchImages fetch and return the images to update to
func (r *ReconcileManager) fetchImages(entandoAppV2 v1alpha1.EntandoAppV2) common.EntandoAppImages {

	imageManager := &common.ImageManager{Log: r.log}
	images := imageManager.FetchImagesByAppVersion(entandoAppV2.Spec.Version)
	if images == nil {
		images = common.EntandoAppImages{}
		r.log.Info("The catalog does not contain the requested App Version ",
			"version", entandoAppV2.Spec.Version)
	}

	for k, v := range entandoAppV2.Spec.ImagesOverride {
		images[k] = v
	}

	r.log.Info("image", "appbuilder", images.FetchAppBuilder())
	r.log.Info("image", "cm", images.FetchComponentManager())
	r.log.Info("image", "de-app", images.FetchDeApp())
	r.log.Info("image", "kc", images.FetchKeycloak())

	return images

}

func (r *ReconcileManager) getDeployment(ctx context.Context, namespace string, deploymentName string) (appsv1.Deployment, error) {
	deployment := appsv1.Deployment{}
	lookupKey := types.NamespacedName{Namespace: namespace, Name: deploymentName}

	err := r.Get(ctx, lookupKey, &deployment)

	return deployment, err
}

func (r *ReconcileManager) reconcileCsv(ctx context.Context, req ctrl.Request, appImages common.EntandoAppImages) error {
	csv := &csv.ClusterServiceVersionList{}

	if err := r.Client.List(ctx, csv); err != nil {
		return err
	}

	for _, item := range csv.Items {
		if strings.HasPrefix(item.ObjectMeta.Name, "entando-k8s-operator") {
			for _, image := range item.Spec.RelatedImages {
				switch {
				case image.Name == "entando-component-manager-6-4":
					image.Name = appImages.FetchComponentManager()
				}
			}
			for _, deploy := range item.Spec.InstallStrategy.StrategySpec.DeploymentSpecs {
				if deploy.Name == "entando-operator" {
					for _, env := range deploy.Spec.Template.Spec.Containers[0].Env {
						switch {
						case env.Name == "RELATED_IMAGE_ENTANDO_COMPONENT_MANAGER_6_4":
							env.Value = appImages.FetchComponentManager()
						}
					}
				}
			}
		}
	}

	return nil
}
