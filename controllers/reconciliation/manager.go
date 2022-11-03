package reconciliation

import (
	"context"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/common"
	"github.com/entgigi/upgrade-operator.git/legacy"
	"github.com/entgigi/upgrade-operator.git/service"
	"github.com/entgigi/upgrade-operator.git/utils"
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
	crReadOnly := &v1alpha1.EntandoAppV2{}
	if err := r.Client.Get(ctx, req.NamespacedName, crReadOnly); err != nil {
		return err
	}

	imageManager := service.NewImageManager(r.Log)
	images := imageManager.FetchAndComposeImagesMap(*crReadOnly)
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
		
		if err = r.reconcileK8sService(ctx, req, images.FetchK8sService(), *crReadOnly); err != nil {
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
