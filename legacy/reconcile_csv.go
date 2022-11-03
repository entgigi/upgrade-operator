package legacy

import (
	"context"
	"strings"

	"github.com/entgigi/upgrade-operator.git/common"
	"github.com/entgigi/upgrade-operator.git/service"
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

func (r *LegacyReconcileManager) ReconcileClusterServiceVersion(ctx context.Context,
	req ctrl.Request,
	appImages service.EntandoAppImages) error {

	r.Log.Info("Starting ClusterServiceVersion reconciliation flow")

	csvList := &csv.ClusterServiceVersionList{}

	if err := r.Client.List(ctx, csvList); err != nil {
		r.Log.Error(err, "Error List ClusterServiceVersion")
		return err
	}

	var csv *csv.ClusterServiceVersion
	for _, item := range csvList.Items {
		if strings.HasPrefix(item.ObjectMeta.Name, "entando-k8s-operator") {
			csv = &item
			break
		}
	}

	if csv != nil {
		r.setRelatedImages(csv, appImages)

		r.setCoordinatorEnvs(csv, appImages)

		r.setK8sServiceDeployment(csv, appImages)

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

func (r *LegacyReconcileManager) setRelatedImages(csv *csv.ClusterServiceVersion, appImages service.EntandoAppImages) {
	for i, entry := range csv.Spec.RelatedImages {
		switch entry.Name {
		case "app-builder-6-4":
			csv.Spec.RelatedImages[i].Image = appImages.FetchAppBuilder()
		case "entando-component-manager-6-4":
			csv.Spec.RelatedImages[i].Image = appImages.FetchComponentManager()
		case "entando-de-app-eap-6-4":
			csv.Spec.RelatedImages[i].Image = appImages.FetchDeApp()
		case "entando-redhat-sso":
			csv.Spec.RelatedImages[i].Image = appImages.FetchKeycloak()
		case "entando-k8s-service":
			csv.Spec.RelatedImages[i].Image = appImages.FetchK8sService()
		case "entando-k8s-plugin-controller":
			csv.Spec.RelatedImages[i].Image = appImages.FetchK8sPluginController()
		case "entando-k8s-app-plugin-link-controller":
			csv.Spec.RelatedImages[i].Image = appImages.FetchK8sAppPluginLinkController()
		}
	}
}

func (r *LegacyReconcileManager) setCoordinatorEnvs(csv *csv.ClusterServiceVersion, appImages service.EntandoAppImages) {
	for j, deploy := range (*csv).Spec.InstallStrategy.StrategySpec.DeploymentSpecs {
		if deploy.Name == "entando-operator" {
			r.Log.Info("ClusterServiceVersion deployment entando-operator found")
			for i, env := range deploy.Spec.Template.Spec.Containers[0].Env {
				switch env.Name {
				case "RELATED_IMAGE_APP_BUILDER_6_4":
					/*
						r.Log.Info("AppBuilder image replace",
							"oldVersion",
							deploy.Spec.Template.Spec.Containers[0].Env[i].Value,
							"newVersion",
							appImages.FetchAppBuilder())
					*/
					setCsvEnvValue(csv, appImages.FetchAppBuilder(), j, i)
				case "RELATED_IMAGE_ENTANDO_COMPONENT_MANAGER_6_4":
					setCsvEnvValue(csv, appImages.FetchComponentManager(), j, i)
				case "RELATED_IMAGE_ENTANDO_DE_APP_EAP_6_4":
					setCsvEnvValue(csv, appImages.FetchDeApp(), j, i)
				case "RELATED_IMAGE_ENTANDO_REDHAT_SSO":
					setCsvEnvValue(csv, appImages.FetchKeycloak(), j, i)
				case "RELATED_IMAGE_ENTANDO_K8S_SERVICE":
					setCsvEnvValue(csv, appImages.FetchK8sService(), j, i)
				case "RELATED_IMAGE_ENTANDO_K8S_PLUGIN_CONTROLLER":
					setCsvEnvValue(csv, appImages.FetchK8sPluginController(), j, i)
				case "RELATED_IMAGE_ENTANDO_K8S_APP_PLUGIN_LINK_CONTROLLER":
					setCsvEnvValue(csv, appImages.FetchK8sAppPluginLinkController(), j, i)
				}
			}
		}
	}
	//r.Log.Info(fmt.Sprintf("\n%+v\n\n", *csv))
}

// TODO verify range copy
func setCsvEnvValue(csv *csv.ClusterServiceVersion, image string, indexDeployment int, indexEnv int) {
	(*csv).Spec.InstallStrategy.StrategySpec.DeploymentSpecs[indexDeployment].
		Spec.Template.Spec.Containers[0].Env[indexEnv].Value = image
}

func (r *LegacyReconcileManager) setK8sServiceDeployment(csv *csv.ClusterServiceVersion, appImages service.EntandoAppImages) {
	for j, deploy := range (*csv).Spec.InstallStrategy.StrategySpec.DeploymentSpecs {
		if deploy.Name == "entando-k8s-service" {
			r.Log.Info("ClusterServiceVersion deployment entando-k8s-service found")
			(*csv).Spec.InstallStrategy.StrategySpec.DeploymentSpecs[j].
				Spec.Template.Spec.Containers[0].Image = appImages.FetchK8sService()
		}
	}
}
