package legacy

import (
	"context"
	"strings"

	"github.com/entgigi/upgrade-operator.git/common"
	"github.com/go-logr/logr"
	csv "github.com/operator-framework/api/pkg/operators/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type LegacyReconcileManager struct {
	common.BaseK8sStructure
}

func NewLegacyReconcileManager(client client.Client, log logr.Logger) *LegacyReconcileManager {
	return &LegacyReconcileManager{BaseK8sStructure: common.BaseK8sStructure{Client: client, Log: log}}
}

func (r *LegacyReconcileManager) ReconcileClusterServiceVersion(ctx context.Context, req ctrl.Request, appImages common.EntandoAppImages) error {
	r.Log.Info("Starting ClusterServiceVersion reconciliation flow")

	csvList := &csv.ClusterServiceVersionList{}

	if err := r.Client.List(ctx, csvList); err != nil {
		r.Log.Error(err, "Error List ClusterServiceVersion")
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

		r.setCoordinatorEnvs(csv, appImages)

		//r.Log.Info(fmt.Sprintf("%+v\n", *csv))

		if err := r.Client.Update(ctx, csv); err != nil {
			r.Log.Error(err, "Error Update ClusterServiceVersion")
			return err
		}
	} else {
		r.Log.Info("No ClusterServiceVersion for entando-k8s-operator found")
	}
	r.Log.Info("Finished ClusterServiceVersion reconciliation flow")
	return nil
}

func setRelatedImages(csv *csv.ClusterServiceVersion, appImages common.EntandoAppImages) {
	for i, entry := range csv.Spec.RelatedImages {
		switch {
		case entry.Name == "app-builder-6-4":
			csv.Spec.RelatedImages[i].Image = appImages.FetchAppBuilder()
		case entry.Name == "entando-component-manager-6-4":
			csv.Spec.RelatedImages[i].Image = appImages.FetchComponentManager()
		case entry.Name == "entando-de-app-eap-6-4":
			csv.Spec.RelatedImages[i].Image = appImages.FetchDeApp()
		case entry.Name == "entando-redhat-sso":
			csv.Spec.RelatedImages[i].Image = appImages.FetchKeycloak()
		case entry.Name == "entando-k8s-service":
			csv.Spec.RelatedImages[i].Image = appImages.FetchK8sService()
		}
	}
}

func (r *LegacyReconcileManager) setCoordinatorEnvs(csv *csv.ClusterServiceVersion, appImages common.EntandoAppImages) {
	for _, deploy := range (*csv).Spec.InstallStrategy.StrategySpec.DeploymentSpecs {
		if deploy.Name == "entando-operator" {
			r.Log.Info("ClusterServiceVersion deployment entando-operator found")
			for i, env := range deploy.Spec.Template.Spec.Containers[0].Env {
				switch {
				case env.Name == "RELATED_IMAGE_APP_BUILDER_6_4":
					/*
						r.Log.Info("AppBuilder image replace",
							"oldVersion",
							deploy.Spec.Template.Spec.Containers[0].Env[i].Value,
							"newVersion",
							appImages.FetchAppBuilder())
					*/
					deploy.Spec.Template.Spec.Containers[0].Env[i].Value = appImages.FetchAppBuilder()
				case env.Name == "RELATED_IMAGE_ENTANDO_COMPONENT_MANAGER_6_4":
					deploy.Spec.Template.Spec.Containers[0].Env[i].Value = appImages.FetchComponentManager()
				case env.Name == "RELATED_IMAGE_ENTANDO_DE_APP_EAP_6_4":
					deploy.Spec.Template.Spec.Containers[0].Env[i].Value = appImages.FetchDeApp()
				case env.Name == "RELATED_IMAGE_ENTANDO_REDHAT_SSO":
					deploy.Spec.Template.Spec.Containers[0].Env[i].Value = appImages.FetchKeycloak()
				case env.Name == "RELATED_IMAGE_ENTANDO_K8S_SERVICE":
					deploy.Spec.Template.Spec.Containers[0].Env[i].Value = appImages.FetchK8sService()
				}
			}
		}
	}
	//r.Log.Info(fmt.Sprintf("\n%+v\n\n", *csv))
}
