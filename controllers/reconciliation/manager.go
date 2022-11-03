package reconciliation

import (
	"context"
	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/common"
	"github.com/entgigi/upgrade-operator.git/legacy"
	"github.com/go-logr/logr"
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
	cr := &v1alpha1.EntandoAppV2{}
	if err := r.Client.Get(ctx, req.NamespacedName, cr); err != nil {
		return err
	}
	images := r.fetchImages(*cr)

	//TODO reconcile secrets for ca before EntandoApp components

	if _, err = r.reconcileComponent(ctx, req, "Keycloak", r.reconcileKeycloak, images.FetchKeycloak(), cr); err != nil {
		return err
	}

	if _, err = r.reconcileComponent(ctx, req, "DeApp", r.reconcileDeApp, images.FetchDeApp(), cr); err != nil {
		return err
	}

	if _, err = r.reconcileComponent(ctx, req, "AppBuilder", r.reconcileAppBuilder, images.FetchAppBuilder(), cr); err != nil {
		return err
	}

	// TODO before check entando-k8s-service app ready
	if cr, err = r.reconcileComponent(ctx, req, "ComponentManager", r.reconcileComponentManager, images.FetchComponentManager(), cr); err != nil {
		return err
	}

	// progress step not added because is not a business step but jsut technical
	r.statusUpdater.SetReconcileProcessingComponent(ctx, req.NamespacedName, "Csv")
	csvReconcile := legacy.NewLegacyReconcileManager(r.Client, r.Log)
	if err = csvReconcile.ReconcileClusterServiceVersion(ctx, req, images); err != nil {
		r.statusUpdater.SetReconcileFailed(ctx, req.NamespacedName, "CsvReconciliationFailed")
		return err
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

// fetchImages fetch and return the images to update to
func (r *ReconcileManager) fetchImages(entandoAppV2 v1alpha1.EntandoAppV2) common.EntandoAppImages {

	imageManager := &common.ImageManager{Log: r.Log}
	images := imageManager.FetchImagesByAppVersion(entandoAppV2.Spec.Version)
	if images == nil {
		images = common.EntandoAppImages{}
		r.Log.Info("The catalog does not contain the requested App Version ",
			"version", entandoAppV2.Spec.Version)
	}

	images[common.AppBuilderKey] = entandoAppV2.Spec.AppBuilder.ImageOverride
	images[common.ComponentManagerKey] = entandoAppV2.Spec.ComponentManager.ImageOverride
	images[common.DeAppKey] = entandoAppV2.Spec.DeApp.ImageOverride
	images[common.KeycloakKey] = entandoAppV2.Spec.Keycloak.ImageOverride

	r.Log.Info("image", "appbuilder", images.FetchAppBuilder())
	r.Log.Info("image", "cm", images.FetchComponentManager())
	r.Log.Info("image", "de-app", images.FetchDeApp())
	r.Log.Info("image", "kc", images.FetchKeycloak())

	return images

}
