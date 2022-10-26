package common

import (
	"github.com/go-logr/logr"
)

type ImageManager struct {
	Log logr.Logger
}

type EntandoAppImages map[string]string

type EntandoAppList map[string]EntandoAppImages

// FIXME! read from yaml ???
// FIXME! use struct not map[string] to obtain immutability with constants
var apps = EntandoAppList{
	"7.1.0": EntandoAppImages{
		AppBuilderKey: "registry.hub.docker.com/entando/app-builder:7.1.0",
	},
	"7.1.1": EntandoAppImages{
		AppBuilderKey: "registry.hub.docker.com/entando/app-builder:7.1.1",
	},
}

func (i *ImageManager) FetchImagesByAppVersion(version string) EntandoAppImages {
	log := i.Log.WithName("common")

	log.Info("Fetch entando app images", "version", version)

	if images, ok := apps[version]; ok {
		return images

	} else {
		return nil

	}

}
