package common

import (
	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/utils"
	"github.com/go-logr/logr"
)

type ImageManager struct {
	Log logr.Logger
}

func NewImageManager(log logr.Logger) *ImageManager {
	return &ImageManager{Log: log}
}

type EntandoAppImages map[string]interface{}

func (e EntandoAppImages) FetchAppBuilder() string {
	k, _ := e[AppBuilderKey].(string)
	return k
}

func (e EntandoAppImages) FetchComponentManager() string {
	k, _ := e[ComponentManagerKey].(string)
	return k
}

func (e EntandoAppImages) FetchDeApp() string {
	var k string
	if utils.IsOlmInstallation() {
		k, _ = e[DeAppEapKey].(string)
	} else {
		k, _ = e[DeAppKey].(string)
	}
	return k
}

func (e EntandoAppImages) FetchKeycloak() string {
	var k string
	if utils.IsOlmInstallation() {
		k, _ = e[KeycloakSsoKey].(string)
	} else {
		k, _ = e[KeycloakKey].(string)
	}
	return k
}

func (e EntandoAppImages) FetchK8sService() string {
	k, _ := e[K8sService].(string)
	return k
}

func (e EntandoAppImages) FetchK8sPluginController() string {
	k, _ := e[K8sPluginController].(string)
	return k
}

func (e EntandoAppImages) FetchK8sAppPluginLinkController() string {
	k, _ := e[K8sAppPluginLinkController].(string)
	return k
}

type EntandoAppList map[string]EntandoAppImages

// TODO read from yaml ???
// FIXME! use struct not map[string] to obtain immutability with constants
var apps = EntandoAppList{
	"7.0.2": EntandoAppImages{
		// "registry.hub.docker.com/entando/app-builder:7.0.2"
		AppBuilderKey: "registry.hub.docker.com/entando/app-builder@sha256:9d54b125b96c861dcf4af311dc47e1f54ee58810c3b16bd50b9f4a15fcf85edd",
		//"registry.hub.docker.com/entando/entando-component-manager:7.0.2"
		ComponentManagerKey: "registry.hub.docker.com/entando/entando-component-manager@sha256:6dfa8c2c910e0d36feee404fa72fafc81d138607724fc5a017aeea461716ceba",
		// "registry.hub.docker.com/entando/entando-de-app-eap:7.0.2"
		DeAppEapKey: "registry.hub.docker.com/entando/entando-de-app-eap@sha256:9aef7599961026b0f6f037e2e66089ce2a0d24055678adc68c2a262842e679f0",
		// "registry.hub.docker.com/entando/entando-de-app-wildfly:7.0.2"
		DeAppKey: "registry.hub.docker.com/entando/entando-de-app-wildfly@sha256:4905d244a62778ccf4990cae20e35fcd9af6f7662fddab80a9757dcecfd18699",
		// "registry.hub.docker.com/entando/entando-redhat-sso:7.0.0"
		KeycloakSsoKey: "registry.hub.docker.com/entando/entando-redhat-sso@sha256:d91c8472a8676d884e789758391cfc36dbfc89318e5293e82d04411160bd132a",
		// "registry.hub.docker.com/entando/entando-keycloak:7.0.1"
		KeycloakKey: "registry.hub.docker.com/entando/entando-keycloak@sha256:e251ba0bf83c529bd1818d4b78cc74ed25b6573d8dc55f59d72b86194932c3ac",
		// "registry.hub.docker.com/entando/entando-k8s-service:7.0.2"
		K8sService: "registry.hub.docker.com/entando/entando-k8s-service@sha256:d1e022d6c1eb1f20268372493115fb34066b52245995eee11b77b7c36bc66431",
		// "registry.hub.docker.com/entando/entando-k8s-plugin-controller:7.0.2"
		K8sPluginController: "registry.hub.docker.com/entando/entando-k8s-plugin-controller@sha256:8fa2a1f44ca9b94797bb5d65d409c58d924cabb60dada6bd2cb1a654adf605e3",
		// "registry.hub.docker.com/entando/entando-k8s-app-plugin-link-controller:7.0.2"
		K8sAppPluginLinkController: "registry.hub.docker.com/entando/entando-k8s-app-plugin-link-controller@sha256:bf71073f3e4d6c6815cbdc582a984695b5944cf3c82e603386a96bc0a1b4bd74",
	},
	"7.1.0": EntandoAppImages{
		// "registry.hub.docker.com/entando/app-builder:7.1.0"
		AppBuilderKey: "registry.hub.docker.com/entando/app-builder@sha256:b09c3d47d667f0f58c1837e6427b37231f90a6c71b9e62c2bbc5c2633b9b3a55",
		//"registry.hub.docker.com/entando/entando-component-manager:7.1.0"
		ComponentManagerKey: "registry.hub.docker.com/entando/entando-component-manager@sha256:7dfb73e37b14396450c695bcaa2b3a3a0b06b1cdfcbdad72d1d50e3150d387e4",
		// "registry.hub.docker.com/entando/entando-de-app-eap:7.1.0"
		DeAppEapKey: "registry.hub.docker.com/entando/entando-de-app-eap@sha256:b64761667234af704088c88ac97dfe25cad531e2390bb585d49f90d18e2bf535",
		// "registry.hub.docker.com/entando/entando-de-app-wildfly:7.1.0"
		DeAppKey: "registry.hub.docker.com/entando/entando-de-app-wildfly@sha256:ac062b7b86fd0e0c64af6426137fcb54e44eed2edcea3b717df5414a2e111d35",
		// "registry.hub.docker.com/entando/entando-redhat-sso:7.1.0"
		KeycloakSsoKey: "registry.hub.docker.com/entando/entando-redhat-sso@sha256:d91c8472a8676d884e789758391cfc36dbfc89318e5293e82d04411160bd132a",
		// "registry.hub.docker.com/entando/entando-keycloak:7.0.1"
		KeycloakKey: "registry.hub.docker.com/entando/entando-keycloak@sha256:e251ba0bf83c529bd1818d4b78cc74ed25b6573d8dc55f59d72b86194932c3ac",
		// "registry.hub.docker.com/entando/entando-k8s-service:7.1.0"
		K8sService: "registry.hub.docker.com/entando/entando-k8s-service@sha256:d7f953ae5627f35a22ee1264e2017f8965b7db6df788a2d2d80e5be82eb7d52b",
		// "registry.hub.docker.com/entando/entando-k8s-plugin-controller:7.1.0"
		K8sPluginController: "registry.hub.docker.com/entando/entando-k8s-plugin-controller@sha256:fd58e7d59a20643735b9bbae6dac8d3b9f2f44c494c674ba13dd0bcc40bf66a9",
		// "registry.hub.docker.com/entando/entando-k8s-app-plugin-link-controller:7.1.0"
		K8sAppPluginLinkController: "registry.hub.docker.com/entando/entando-k8s-app-plugin-link-controller@sha256:d9e9cdcbf2abec4b0d1253955d4dffd01de284d32f40ac42ec18aa3e94e32db4",
	},
	"7.1.1": EntandoAppImages{
		// "registry.hub.docker.com/entando/app-builder:7.1.1"
		AppBuilderKey: "registry.hub.docker.com/entando/app-builder@sha256:33ea636090352a919735aa44cc2aaf2c79e8cb15b19216574964ec41c98f5c58",
		// "registry.hub.docker.com/entando/entando-component-manager:7.1.1"
		ComponentManagerKey: "registry.hub.docker.com/entando/entando-component-manager@sha256:13d9ba8a9a3cb52f4f999ece27f3ba8470fabd77a8c31290b5fdcc615e1dff11",
		// "registry.hub.docker.com/entando/entando-de-app-eap:7.1.1"
		DeAppEapKey: "registry.hub.docker.com/entando/entando-de-app-eap@sha256:21c87b5f0069e38b864211427a39d6d135205d81e24371ddc987f50faa0c21b0",
		// "registry.hub.docker.com/entando/entando-de-app-wildfly:7.1.1"
		DeAppKey: "registry.hub.docker.com/entando/entando-de-app-wildfly@sha256:ee2966c9cc6fe3a258f8113af8801e830a0029181b23abebfc30cb8f16cdd14f",
		// "registry.hub.docker.com/entando/entando-redhat-sso:7.1.1"
		KeycloakSsoKey: "registry.hub.docker.com/entando/entando-redhat-sso@sha256:b5afae1e933d432ccab84e31a9140f8fd2a51517b95a3a373a29bbe88a62d900",
		// "registry.hub.docker.com/entando/entando-keycloak:7.1.1"
		KeycloakKey: "registry.hub.docker.com/entando/entando-keycloak@sha256:4b5b81a6f233e070a747541ec1fa30cdaeca78feefa1542cb117aaaf2079863c",
		// "registry.hub.docker.com/entando/entando-k8s-service:7.1.1"
		K8sService: "registry.hub.docker.com/entando/entando-k8s-service@sha256:df0473993a7eb6dd71fd06dcd3b31efebd470e9c9962fb2ccc85e6ee356de3cd",
		// "registry.hub.docker.com/entando/entando-k8s-plugin-controller:7.1.1"
		K8sPluginController: "registry.hub.docker.com/entando/entando-k8s-plugin-controller@sha256:1084d03bd9ecf2a720390b7e1543f60b43c4baa523b97a85e7e54590d81d2574",
		// "registry.hub.docker.com/entando/entando-k8s-app-plugin-link-controller:7.1.0"
		K8sAppPluginLinkController: "registry.hub.docker.com/entando/entando-k8s-app-plugin-link-controller@sha256:d9e9cdcbf2abec4b0d1253955d4dffd01de284d32f40ac42ec18aa3e94e32db4",
	},
}

func (i *ImageManager) FetchImagesByAppVersion(version string) EntandoAppImages {
	log := i.Log.WithName("common")

	log.Info("Fetch entando app images", "version", version)

	if images, ok := apps[version]; ok {
		return utils.CopyMap(images)
	} else {
		log.Info("Entando version not found")
		return nil

	}
}

// FetchAndComposeImagesMap fetch and return the images to update to
func (r *ImageManager) FetchAndComposeImagesMap(entandoAppV2 v1alpha1.EntandoAppV2) EntandoAppImages {

	images := r.FetchImagesByAppVersion(entandoAppV2.Spec.Version)
	if images == nil {
		images = EntandoAppImages{}
		r.Log.Info("The catalog does not contain the requested App Version ",
			"version", entandoAppV2.Spec.Version)
	}

	if len(entandoAppV2.Spec.AppBuilder.ImageOverride) > 0 {
		images[AppBuilderKey] = entandoAppV2.Spec.AppBuilder.ImageOverride
	}
	if len(entandoAppV2.Spec.ComponentManager.ImageOverride) > 0 {
		images[ComponentManagerKey] = entandoAppV2.Spec.ComponentManager.ImageOverride
	}
	if len(entandoAppV2.Spec.DeApp.ImageOverride) > 0 {
		images[DeAppKey] = entandoAppV2.Spec.DeApp.ImageOverride
	}
	if len(entandoAppV2.Spec.Keycloak.ImageOverride) > 0 {
		images[KeycloakKey] = entandoAppV2.Spec.Keycloak.ImageOverride
	}
	if len(entandoAppV2.Spec.K8sService.ImageOverride) > 0 {
		images[K8sService] = entandoAppV2.Spec.K8sService.ImageOverride
	}
	if len(entandoAppV2.Spec.K8sPluginController.ImageOverride) > 0 {
		images[K8sPluginController] = entandoAppV2.Spec.K8sPluginController.ImageOverride
	}
	if len(entandoAppV2.Spec.K8sAppPluginLinkController.ImageOverride) > 0 {
		images[K8sAppPluginLinkController] = entandoAppV2.Spec.K8sAppPluginLinkController.ImageOverride
	}

	r.Log.Info("image", "app-builder", images.FetchAppBuilder())
	r.Log.Info("image", "component-manager", images.FetchComponentManager())
	r.Log.Info("image", "de-app", images.FetchDeApp())
	r.Log.Info("image", "keycloak", images.FetchKeycloak())
	r.Log.Info("image", "k8s-service", images.FetchK8sService())
	r.Log.Info("image", "k8s-plugin-controller", images.FetchK8sPluginController())
	r.Log.Info("image", "k8s-app-plugin-link-controller", images.FetchK8sAppPluginLinkController())

	return images
}
