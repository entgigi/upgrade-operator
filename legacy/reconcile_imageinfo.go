package legacy

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/entgigi/upgrade-operator.git/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

func (r *LegacyReconcileManager) ReconcileImageInfo(ctx context.Context, req ctrl.Request, appImages common.EntandoAppImages) error {
	r.Log.Info("Starting ImageInfo configMap reconciliation flow")

	configMap := &corev1.ConfigMap{}
	var imageInfo types.NamespacedName = types.NamespacedName{
		Name:      common.ImageInfoConfigMap,
		Namespace: req.NamespacedName.Namespace,
	}
	if err := r.Client.Get(ctx, imageInfo, configMap); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("ImageInfo configMap not found")
			return nil
		}
		r.Log.Error(err, "Error get ImageInfo configMap")
		return err
	}

	for value, key := range configMap.Data {
		switch {
		case key == "app-builder-6-4":
			configMap.Data[key] = buildNewValue(value, appImages.FetchAppBuilder())
		case key == "entando-component-manager-6-4":
			configMap.Data[key] = buildNewValue(value, appImages.FetchComponentManager())
		case key == "entando-de-app-wildfly-6-4":
			configMap.Data[key] = buildNewValue(value, appImages.FetchDeApp())
		case key == "entando-k8s-keycloak-controller":
			configMap.Data[key] = buildNewValue(value, appImages.FetchKeycloak())
		case key == "entando-k8s-service":
			configMap.Data[key] = buildNewValue(value, appImages.FetchK8sService())
		case key == "entando-k8s-plugin-controller":
			configMap.Data[key] = buildNewValue(value, appImages.FetchK8sPluginController())
		case key == "entando-k8s-app-plugin-link-controller":
			configMap.Data[key] = buildNewValue(value, appImages.FetchK8sAppPluginLinkController())
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

func buildNewValue(value string, imageUrl string) string {
	newValue := convertJsonStringToStruct(value)
	newValue.Version = retrieveDigestFromImageFullUrl(imageUrl)
	return convertStructToJsonString(newValue)
}

func convertJsonStringToStruct(s string) ImageInfoEntry {
	data := ImageInfoEntry{}
	json.Unmarshal([]byte(s), &data)
	return data
}

func convertStructToJsonString(data ImageInfoEntry) string {
	b, _ := json.Marshal(&data)
	return string(b)
}

func retrieveDigestFromImageFullUrl(url string) string {
	return strings.Split(url, "@")[1]
}
