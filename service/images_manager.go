package service

import (
	"fmt"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/common"
	"github.com/entgigi/upgrade-operator.git/utils"
	"github.com/go-logr/logr"
)

type ImageManager struct {
	Log logr.Logger
}

func NewImageManager(log logr.Logger) *ImageManager {
	return &ImageManager{Log: log}
}

func (i *ImageManager) FetchImagesByAppVersion(cr *v1alpha1.EntandoAppV2) EntandoAppImages {
	log := i.Log.WithName("common")
	version := cr.Spec.Version
	log.Info("Fetch entando app images", "version", version)

	if images, ok := apps[version]; ok {
		// WARNING! we do map deep copy to grant immutability to original map
		return EntandoAppImages{utils.CopyMap(images), cr}
	} else {
		log.Info("The catalog does not contain the requested App Version ", "version", cr.Spec.Version)
		return EntandoAppImages{make(entandoAppImages, 0), cr}
	}
}

// FetchAndComposeImagesMap fetch and return the images to update to
func (r *ImageManager) FetchAndComposeImagesMap(cr *v1alpha1.EntandoAppV2) (EntandoAppImages, error) {

	images := r.FetchImagesByAppVersion(cr)

	if err := r.checkTagOrDigest(*cr); err != nil {
		return EntandoAppImages{}, err
	}

	if len(cr.Spec.AppBuilder.ImageOverride) > 0 {
		key := ChooseKeyWithTagOrDigest(common.AppBuilderKey)
		images.images[key] = cr.Spec.AppBuilder.ImageOverride
	}
	if len(cr.Spec.ComponentManager.ImageOverride) > 0 {
		key := ChooseKeyWithTagOrDigest(common.ComponentManagerKey)
		images.images[key] = cr.Spec.ComponentManager.ImageOverride
	}
	if len(cr.Spec.DeApp.ImageOverride) > 0 {
		key := ChooseKeyByTypeAndVersion(cr, common.DeAppEapKey, common.DeAppKey)
		images.images[key] = cr.Spec.DeApp.ImageOverride
	}
	if len(cr.Spec.Keycloak.ImageOverride) > 0 {
		key := ChooseKeyByTypeAndVersion(cr, common.KeycloakSsoKey, common.KeycloakKey)
		images.images[key] = cr.Spec.Keycloak.ImageOverride
	}
	if len(cr.Spec.K8sService.ImageOverride) > 0 {
		key := ChooseKeyWithTagOrDigest(common.K8sServiceKey)
		images.images[key] = cr.Spec.K8sService.ImageOverride
	}
	if len(cr.Spec.K8sPluginController.ImageOverride) > 0 {
		key := ChooseKeyWithTagOrDigest(common.K8sPluginControllerKey)
		images.images[key] = cr.Spec.K8sPluginController.ImageOverride
	}
	if len(cr.Spec.K8sAppPluginLinkController.ImageOverride) > 0 {
		key := ChooseKeyWithTagOrDigest(common.K8sAppPluginLinkControllerKey)
		images.images[key] = cr.Spec.K8sAppPluginLinkController.ImageOverride
	}

	r.Log.Info("image", "app-builder", images.FetchAppBuilder())
	r.Log.Info("image", "component-manager", images.FetchComponentManager())
	r.Log.Info("image", "de-app", images.FetchDeApp())
	r.Log.Info("image", "keycloak", images.FetchKeycloak())
	r.Log.Info("image", "k8s-service", images.FetchK8sService())
	r.Log.Info("image", "k8s-plugin-controller", images.FetchK8sPluginController())
	r.Log.Info("image", "k8s-app-plugin-link-controller", images.FetchK8sAppPluginLinkController())

	return images, nil
}

func (r *ImageManager) checkTagOrDigest(entandoAppV2 v1alpha1.EntandoAppV2) error {
	imagesToCheck := []string{
		entandoAppV2.Spec.AppBuilder.ImageOverride,
		entandoAppV2.Spec.ComponentManager.ImageOverride,
		entandoAppV2.Spec.DeApp.ImageOverride,
		entandoAppV2.Spec.Keycloak.ImageOverride,
		entandoAppV2.Spec.K8sService.ImageOverride,
		entandoAppV2.Spec.K8sPluginController.ImageOverride,
		entandoAppV2.Spec.K8sAppPluginLinkController.ImageOverride,
	}
	for _, image := range imagesToCheck {
		if len(image) > 0 {
			imageInfo, err := NewImageInfo(image)
			if err != nil {
				r.Log.Error(err, "Error parsing image url", "imageUrl", image)
				return err
			}
			if utils.IsOlmInstallation() && imageInfo.IsTag() {
				r.Log.Error(err, "Error image url contains tag in an OLM installation", "imageUrl", image)
				return fmt.Errorf("Error image url:'%s' contains tag in an OLM installation", image)
			}
			if !utils.IsOlmInstallation() && imageInfo.IsDigest() {
				r.Log.Error(err, "Error image url contains digest in an non-OLM installation", "imageUrl", image)
				return fmt.Errorf("Error image url:'%s' contains tag in an non-OLM installation", image)
			}
		}
	}

	return nil
}
