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
	k, _ := e[AppBuilderKey].(string)
	return k
}

func (e EntandoAppImages) FetchComponentManager() string {
	k, _ := e[ComponentManagerKey].(string)
	return k
}

func (e EntandoAppImages) FetchDeApp() string {
	k, _ := e[DeAppKey].(string)
	return k
}

func (e EntandoAppImages) FetchKeycloak() string {
	k, _ := e[KeycloakKey].(string)
	return k
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
