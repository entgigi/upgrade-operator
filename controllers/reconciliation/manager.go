package reconciliation

import (
	"context"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/common"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
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
