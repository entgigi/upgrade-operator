package legacy

import (
	"context"
	"strings"

	"github.com/entgigi/upgrade-operator.git/common"
	csv "github.com/operator-framework/api/pkg/operators/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type LegacyReconcileManager struct {
	common.BaseK8sStructure
}

func (r *LegacyReconcileManager) ReconcileClusterServiceVersion(ctx context.Context, req ctrl.Request, appImages common.EntandoAppImages) error {
	csvList := &csv.ClusterServiceVersionList{}

	if err := r.Client.List(ctx, csvList); err != nil {
		return nil
	}

	var csv *csv.ClusterServiceVersion
	for _, item := range csvList.Items {
		if strings.HasPrefix(item.ObjectMeta.Name, "entando-k8s-operator") {
			csv = &item
			break
		}
	}

	if csv != nil {
		setRelatedImages(csv, appImages)

		setCoordinatorEnvs(csv, appImages)

		if err := r.Client.Update(ctx, csv); err != nil {
			return err
		}
	}
	return nil
}

func setRelatedImages(csv *csv.ClusterServiceVersion, appImages common.EntandoAppImages) {
	for _, image := range csv.Spec.RelatedImages {
		switch {
		case image.Name == "entando-component-manager-6-4":
			image.Name = appImages.FetchComponentManager()
		}
	}

}

func setCoordinatorEnvs(csv *csv.ClusterServiceVersion, appImages common.EntandoAppImages) {
	for _, deploy := range csv.Spec.InstallStrategy.StrategySpec.DeploymentSpecs {
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
