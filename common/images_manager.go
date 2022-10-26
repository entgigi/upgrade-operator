package common

import (
	"github.com/entgigi/upgrade-operator.git/utils"
	"github.com/go-logr/logr"
)

type ImageManager struct {
	Log logr.Logger
}

type EntandoAppImages map[string]interface{}

func (e EntandoAppImages) FetchAppBuilder() string {
	return e[AppBuilderKey].(string)
}

func (e EntandoAppImages) FetchComponentManager() string {
	return e[ComponentManagerKey].(string)
}

func (e EntandoAppImages) FetchDeApp() string {
	return e[DeAppKey].(string)
}

func (e EntandoAppImages) FetchKeycloak() string {
	return e[KeycloakKey].(string)
}

type EntandoAppList map[string]EntandoAppImages

// FIXME! read from yaml ???
// FIXME! use struct not map[string] to obtain immutability with constants
var apps = EntandoAppList{
	"7.1.0": EntandoAppImages{
		AppBuilderKey:       "registry.hub.docker.com/entando/app-builder:7.1.0",
		ComponentManagerKey: "registry.hub.docker.com/entando/entando-component-manager:7.1.0",
		DeAppKey:            "registry.hub.docker.com/entando/entando-de-app-eap:7.1.0",
		KeycloakKey:         "registry.hub.docker.com/entando/entando-keycloak:7.1.0",
	},
	"7.1.1": EntandoAppImages{
		AppBuilderKey:       "registry.hub.docker.com/entando/app-builder:7.1.1",
		ComponentManagerKey: "registry.hub.docker.com/entando/entando-component-manager:7.1.1",
		DeAppKey:            "registry.hub.docker.com/entando/entando-de-app-eap:7.1.1",
		KeycloakKey:         "registry.hub.docker.com/entando/entando-keycloak:7.1.1",
	},
}

func (i *ImageManager) FetchImagesByAppVersion(version string) EntandoAppImages {
	log := i.Log.WithName("common")

	log.Info("Fetch entando app images", "version", version)

	if images, ok := apps[version]; ok {
		return utils.CopyMap(images)

	} else {
		return nil

	}

}
