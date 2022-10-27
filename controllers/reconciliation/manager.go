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
)

type ReconcileManager struct {
	client.Client
	log    logr.Logger
	scheme *runtime.Scheme
}

// NewReconcileManager initialize a ReconcileManager
func NewReconcileManager(client client.Client, log logr.Logger, scheme *runtime.Scheme) *ReconcileManager {
	logger := log.WithName(reconciliationLogName)
	return &ReconcileManager{
		client,
		logger,
		scheme,
	}
}

func (r *ReconcileManager) MainReconcile(ctx context.Context, req ctrl.Request) error {

	r.log.Info("Starting main reconciliation flow")

	cr := &v1alpha1.EntandoAppV2{}
	if err := r.Client.Get(ctx, req.NamespacedName, cr); err != nil {
		return err
	}
	images := r.fetchImages(*cr)

	if err := r.reconcileKeycloak(ctx, images.FetchKeycloak(), req); err != nil {
		return err
	}

	if err := r.reconcileDeApp(ctx, images.FetchDeApp(), req); err != nil {
		return err
	}

	if err := r.reconcileAppBuilder(ctx, images.FetchAppBuilder(), req); err != nil {
		return err
	}

	// TODO before check entando-k8s-service app ready
	if err := r.reconcileComponentManager(ctx, images.FetchComponentManager(), req); err != nil {
		return err
	}

	// TODO manage CSV

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
