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

func (e EntandoAppImages) chooseKeyWithTagOrDigest(baseKey string) string {
	key := baseKey
	if utils.IsOlmInstallation() {
		key = key + common.TagKey
	}
	k, _ := e.images[key].(string)
	return k

}

func (e EntandoAppImages) FetchAppBuilder() string {
	return e.chooseKeyWithTagOrDigest(common.AppBuilderKey)
}

func (e EntandoAppImages) FetchComponentManager() string {
	return e.chooseKeyWithTagOrDigest(common.ComponentManagerKey)
}

func (e EntandoAppImages) FetchDeApp() string {
	var imageUrl string
	if utils.IsImageSetTypeCommunity(e.cr) {
		key := e.chooseKeyWithTagOrDigest(common.DeAppKey)
		imageUrl, _ = e.images[key].(string)
	} else {
		key := e.chooseKeyWithTagOrDigest(common.DeAppEapKey)
		imageUrl, _ = e.images[key].(string)
	}
	return imageUrl
}

func (e EntandoAppImages) FetchKeycloak() string {
	var k string
	if utils.IsImageSetTypeCommunity(e.cr) {
		key := e.chooseKeyWithTagOrDigest(common.KeycloakKey)
		k, _ = e.images[key].(string)
	} else {
		key := e.chooseKeyWithTagOrDigest(common.KeycloakSsoKey)
		k, _ = e.images[key].(string)
	}
	return k
}

func (e EntandoAppImages) FetchK8sService() string {
	return e.chooseKeyWithTagOrDigest(common.K8sServiceKey)
}

func (e EntandoAppImages) FetchK8sPluginController() string {
	return e.chooseKeyWithTagOrDigest(common.K8sPluginControllerKey)
}

func (e EntandoAppImages) FetchK8sAppPluginLinkController() string {
	return e.chooseKeyWithTagOrDigest(common.K8sAppPluginLinkControllerKey)
}

type EntandoAppImages struct {
	images map[string]interface{}
	cr     *v1alpha1.EntandoAppV2
}

type entandoAppImages map[string]interface{}

type entandoAppList map[string]entandoAppImages

// TODO read from yaml ???
var apps = entandoAppList{
	"7.0.2": entandoAppImages{
		common.AppBuilderKeyTag:                 "registry.hub.docker.com/entando/app-builder:7.0.2",
		common.AppBuilderKey:                    "registry.hub.docker.com/entando/app-builder@sha256:9d54b125b96c861dcf4af311dc47e1f54ee58810c3b16bd50b9f4a15fcf85edd",
		common.ComponentManagerKeyTag:           "registry.hub.docker.com/entando/entando-component-manager:7.0.2",
		common.ComponentManagerKey:              "registry.hub.docker.com/entando/entando-component-manager@sha256:6dfa8c2c910e0d36feee404fa72fafc81d138607724fc5a017aeea461716ceba",
		common.DeAppEapKeyTag:                   "registry.hub.docker.com/entando/entando-de-app-eap:7.0.2",
		common.DeAppEapKey:                      "registry.hub.docker.com/entando/entando-de-app-eap@sha256:9aef7599961026b0f6f037e2e66089ce2a0d24055678adc68c2a262842e679f0",
		common.DeAppKeyTag:                      "registry.hub.docker.com/entando/entando-de-app-wildfly:7.0.2",
		common.DeAppKey:                         "registry.hub.docker.com/entando/entando-de-app-wildfly@sha256:4905d244a62778ccf4990cae20e35fcd9af6f7662fddab80a9757dcecfd18699",
		common.KeycloakSsoKeyTag:                "registry.hub.docker.com/entando/entando-redhat-sso:7.0.0",
		common.KeycloakSsoKey:                   "registry.hub.docker.com/entando/entando-redhat-sso@sha256:d91c8472a8676d884e789758391cfc36dbfc89318e5293e82d04411160bd132a",
		common.KeycloakKeyTag:                   "registry.hub.docker.com/entando/entando-keycloak:7.0.1",
		common.KeycloakKey:                      "registry.hub.docker.com/entando/entando-keycloak@sha256:e251ba0bf83c529bd1818d4b78cc74ed25b6573d8dc55f59d72b86194932c3ac",
		common.K8sServiceKeyTag:                 "registry.hub.docker.com/entando/entando-k8s-service:7.0.2",
		common.K8sServiceKey:                    "registry.hub.docker.com/entando/entando-k8s-service@sha256:d1e022d6c1eb1f20268372493115fb34066b52245995eee11b77b7c36bc66431",
		common.K8sPluginControllerKeyTag:        "registry.hub.docker.com/entando/entando-k8s-plugin-controller:7.0.2",
		common.K8sPluginControllerKey:           "registry.hub.docker.com/entando/entando-k8s-plugin-controller@sha256:8fa2a1f44ca9b94797bb5d65d409c58d924cabb60dada6bd2cb1a654adf605e3",
		common.K8sAppPluginLinkControllerkeyTag: "registry.hub.docker.com/entando/entando-k8s-app-plugin-link-controller:7.0.2",
		common.K8sAppPluginLinkControllerKey:    "registry.hub.docker.com/entando/entando-k8s-app-plugin-link-controller@sha256:bf71073f3e4d6c6815cbdc582a984695b5944cf3c82e603386a96bc0a1b4bd74",
	},
	"7.1.0": entandoAppImages{
		common.AppBuilderKeyTag:                 "registry.hub.docker.com/entando/app-builder:7.1.0",
		common.AppBuilderKey:                    "registry.hub.docker.com/entando/app-builder@sha256:b09c3d47d667f0f58c1837e6427b37231f90a6c71b9e62c2bbc5c2633b9b3a55",
		common.ComponentManagerKeyTag:           "registry.hub.docker.com/entando/entando-component-manager:7.1.0",
		common.ComponentManagerKey:              "registry.hub.docker.com/entando/entando-component-manager@sha256:7dfb73e37b14396450c695bcaa2b3a3a0b06b1cdfcbdad72d1d50e3150d387e4",
		common.DeAppEapKeyTag:                   "registry.hub.docker.com/entando/entando-de-app-eap:7.1.0",
		common.DeAppEapKey:                      "registry.hub.docker.com/entando/entando-de-app-eap@sha256:b64761667234af704088c88ac97dfe25cad531e2390bb585d49f90d18e2bf535",
		common.DeAppKeyTag:                      "registry.hub.docker.com/entando/entando-de-app-wildfly:7.1.0",
		common.DeAppKey:                         "registry.hub.docker.com/entando/entando-de-app-wildfly@sha256:ac062b7b86fd0e0c64af6426137fcb54e44eed2edcea3b717df5414a2e111d35",
		common.KeycloakSsoKeyTag:                "registry.hub.docker.com/entando/entando-redhat-sso:7.1.0",
		common.KeycloakSsoKey:                   "registry.hub.docker.com/entando/entando-redhat-sso@sha256:d91c8472a8676d884e789758391cfc36dbfc89318e5293e82d04411160bd132a",
		common.KeycloakKeyTag:                   "registry.hub.docker.com/entando/entando-keycloak:7.0.1",
		common.KeycloakKey:                      "registry.hub.docker.com/entando/entando-keycloak@sha256:e251ba0bf83c529bd1818d4b78cc74ed25b6573d8dc55f59d72b86194932c3ac",
		common.K8sServiceKeyTag:                 "registry.hub.docker.com/entando/entando-k8s-service:7.1.0",
		common.K8sServiceKey:                    "registry.hub.docker.com/entando/entando-k8s-service@sha256:d7f953ae5627f35a22ee1264e2017f8965b7db6df788a2d2d80e5be82eb7d52b",
		common.K8sPluginControllerKeyTag:        "registry.hub.docker.com/entando/entando-k8s-plugin-controller:7.1.0",
		common.K8sPluginControllerKey:           "registry.hub.docker.com/entando/entando-k8s-plugin-controller@sha256:fd58e7d59a20643735b9bbae6dac8d3b9f2f44c494c674ba13dd0bcc40bf66a9",
		common.K8sAppPluginLinkControllerkeyTag: "registry.hub.docker.com/entando/entando-k8s-app-plugin-link-controller:7.1.0",
		common.K8sAppPluginLinkControllerKey:    "registry.hub.docker.com/entando/entando-k8s-app-plugin-link-controller@sha256:d9e9cdcbf2abec4b0d1253955d4dffd01de284d32f40ac42ec18aa3e94e32db4",
	},
	"7.1.1": entandoAppImages{
		common.AppBuilderKeyTag:                 "registry.hub.docker.com/entando/app-builder:7.1.1",
		common.AppBuilderKey:                    "registry.hub.docker.com/entando/app-builder@sha256:33ea636090352a919735aa44cc2aaf2c79e8cb15b19216574964ec41c98f5c58",
		common.ComponentManagerKeyTag:           "registry.hub.docker.com/entando/entando-component-manager:7.1.1",
		common.ComponentManagerKey:              "registry.hub.docker.com/entando/entando-component-manager@sha256:13d9ba8a9a3cb52f4f999ece27f3ba8470fabd77a8c31290b5fdcc615e1dff11",
		common.DeAppEapKeyTag:                   "registry.hub.docker.com/entando/entando-de-app-eap:7.1.1",
		common.DeAppEapKey:                      "registry.hub.docker.com/entando/entando-de-app-eap@sha256:21c87b5f0069e38b864211427a39d6d135205d81e24371ddc987f50faa0c21b0",
		common.DeAppKeyTag:                      "registry.hub.docker.com/entando/entando-de-app-wildfly:7.1.1",
		common.DeAppKey:                         "registry.hub.docker.com/entando/entando-de-app-wildfly@sha256:ee2966c9cc6fe3a258f8113af8801e830a0029181b23abebfc30cb8f16cdd14f",
		common.KeycloakSsoKeyTag:                "registry.hub.docker.com/entando/entando-redhat-sso:7.1.1",
		common.KeycloakSsoKey:                   "registry.hub.docker.com/entando/entando-redhat-sso@sha256:b5afae1e933d432ccab84e31a9140f8fd2a51517b95a3a373a29bbe88a62d900",
		common.KeycloakKeyTag:                   "registry.hub.docker.com/entando/entando-keycloak:7.1.1",
		common.KeycloakKey:                      "registry.hub.docker.com/entando/entando-keycloak@sha256:4b5b81a6f233e070a747541ec1fa30cdaeca78feefa1542cb117aaaf2079863c",
		common.K8sServiceKeyTag:                 "registry.hub.docker.com/entando/entando-k8s-service:7.1.1",
		common.K8sServiceKey:                    "registry.hub.docker.com/entando/entando-k8s-service@sha256:df0473993a7eb6dd71fd06dcd3b31efebd470e9c9962fb2ccc85e6ee356de3cd",
		common.K8sPluginControllerKeyTag:        "registry.hub.docker.com/entando/entando-k8s-plugin-controller:7.1.1",
		common.K8sPluginControllerKey:           "registry.hub.docker.com/entando/entando-k8s-plugin-controller@sha256:1084d03bd9ecf2a720390b7e1543f60b43c4baa523b97a85e7e54590d81d2574",
		common.K8sAppPluginLinkControllerkeyTag: "registry.hub.docker.com/entando/entando-k8s-app-plugin-link-controller:7.1.0",
		common.K8sAppPluginLinkControllerKey:    "registry.hub.docker.com/entando/entando-k8s-app-plugin-link-controller@sha256:d9e9cdcbf2abec4b0d1253955d4dffd01de284d32f40ac42ec18aa3e94e32db4",
	},
}

func (i *ImageManager) FetchImagesByAppVersion(cr *v1alpha1.EntandoAppV2) EntandoAppImages {
	log := i.Log.WithName("common")
	version := cr.Spec.Version
	log.Info("Fetch entando app images", "version", version)

	if images, ok := apps[version]; ok {
		// WARNING! we do map deep copy to grant immutabiolity to original map
		return EntandoAppImages{utils.CopyMap(images), cr}
	} else {
		log.Info("The catalog does not contain the requested App Version ", "version", cr.Spec.Version)
		return EntandoAppImages{nil, cr}
	}
}

// FetchAndComposeImagesMap fetch and return the images to update to
func (r *ImageManager) FetchAndComposeImagesMap(cr *v1alpha1.EntandoAppV2) (EntandoAppImages, error) {

	images := r.FetchImagesByAppVersion(cr)

	if err := r.checkTagOrDigest(*cr); err != nil {
		return EntandoAppImages{}, err
	}

	if len(cr.Spec.AppBuilder.ImageOverride) > 0 {
		images.images[common.AppBuilderKey] = cr.Spec.AppBuilder.ImageOverride
	}
	if len(cr.Spec.ComponentManager.ImageOverride) > 0 {
		images.images[common.ComponentManagerKey] = cr.Spec.ComponentManager.ImageOverride
	}
	if len(cr.Spec.DeApp.ImageOverride) > 0 {
		images.images[common.DeAppKey] = cr.Spec.DeApp.ImageOverride
	}
	if len(cr.Spec.Keycloak.ImageOverride) > 0 {
		images.images[common.KeycloakKey] = cr.Spec.Keycloak.ImageOverride
	}
	if len(cr.Spec.K8sService.ImageOverride) > 0 {
		images.images[common.K8sServiceKey] = cr.Spec.K8sService.ImageOverride
	}
	if len(cr.Spec.K8sPluginController.ImageOverride) > 0 {
		images.images[common.K8sPluginControllerKey] = cr.Spec.K8sPluginController.ImageOverride
	}
	if len(cr.Spec.K8sAppPluginLinkController.ImageOverride) > 0 {
		images.images[common.K8sAppPluginLinkControllerKey] = cr.Spec.K8sAppPluginLinkController.ImageOverride
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
