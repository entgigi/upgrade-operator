package legacy

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/entgigi/upgrade-operator.git/common"
	"github.com/entgigi/upgrade-operator.git/service"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ImageInfoEntry struct {
	Version        string `json:"version"`
	ExecutableType string `json:"executable-type"`
	Registry       string `json:"registry"`
	Organization   string `json:"organization"`
	Repository     string `json:"repository"`
}

func (r *LegacyReconcileManager) ReconcileImageInfo(ctx context.Context, req ctrl.Request, appImages service.EntandoAppImages) error {
	r.Log.Info("Starting ImageInfo configMap reconciliation flow")

	configMap := &corev1.ConfigMap{}
	var imageInfo types.NamespacedName = types.NamespacedName{
		Name:      common.ImageInfoConfigMap,
		Namespace: req.NamespacedName.Namespace,
	}
	if err := r.Client.Get(ctx, imageInfo, configMap); err != nil {
		r.Log.Error(err, "Error get ImageInfo configMap")
		return err
	}

	for key, value := range configMap.Data {
		var err error
		switch key {
		case "app-builder-6-4":
			r.Log.Info(fmt.Sprintf("app-builder-6-4 %+v\n", appImages.FetchAppBuilder()))
			err = r.buildAndSetNewValue(configMap, key, value, appImages.FetchAppBuilder())

		case "entando-component-manager-6-4":
			err = r.buildAndSetNewValue(configMap, key, value, appImages.FetchComponentManager())

		case "entando-de-app-wildfly-6-4":
			err = r.buildAndSetNewValue(configMap, key, value, appImages.FetchDeApp())

		case "entando-k8s-keycloak-controller":
			err = r.buildAndSetNewValue(configMap, key, value, appImages.FetchKeycloak())

		case "entando-k8s-service":
			err = r.buildAndSetNewValue(configMap, key, value, appImages.FetchK8sService())

		case "entando-k8s-plugin-controller":
			err = r.buildAndSetNewValue(configMap, key, value, appImages.FetchK8sPluginController())

		case "entando-k8s-app-plugin-link-controller":
			err = r.buildAndSetNewValue(configMap, key, value, appImages.FetchK8sAppPluginLinkController())

		}
		if err != nil {
			return err
		}
	}

	r.Log.Info(fmt.Sprintf("%+v\n", configMap.Data))

	if err := r.Client.Update(ctx, configMap); err != nil {
		r.Log.Error(err, "Error update ImageInfo configMap")
		return err
	}

	r.Log.Info("Finished ImageInfo configMap reconciliation flow")
	return nil
}

func (r *LegacyReconcileManager) buildAndSetNewValue(configMap *corev1.ConfigMap,
	key string,
	value string,
	imageUrl string) error {

	image, err := service.NewImageInfo(imageUrl)
	if err != nil {
		r.Log.Error(err, "Error parse fully qualified image url", "imageUrl", imageUrl)
		return err
	}

	newValue, err := convertJsonStringToStruct(value)
	if err != nil {
		r.Log.Error(err, "Error convert json string to data structure ImageInfoEntry", "value", value)
		return err
	}

	newValue = composeNewValue(newValue, image)

	newStringvalue, err := convertStructToJsonString(newValue)
	if err != nil {
		r.Log.Error(err, "Error convert data structure ImageInfoEntry to json string")
		return err
	}

	configMap.Data[key] = newStringvalue

	return nil
}

func convertJsonStringToStruct(s string) (ImageInfoEntry, error) {
	data := ImageInfoEntry{}
	err := json.Unmarshal([]byte(s), &data)
	return data, err
}

func convertStructToJsonString(data ImageInfoEntry) (string, error) {
	b, err := json.Marshal(&data)
	return string(b), err
}

func composeNewValue(oldValue ImageInfoEntry, image service.ImageInfo) ImageInfoEntry {
	oldValue.Registry = image.Hostname()
	oldValue.Organization = image.Org()
	oldValue.Repository = image.Name()
	if image.IsTag() {
		oldValue.Version = image.Tag()
	} else {
		oldValue.Version = image.Digest()
	}
	return oldValue
}

func retrieveDigestFromImageFullUrl(url string) string {
	return strings.Split(url, "@")[1]
}
