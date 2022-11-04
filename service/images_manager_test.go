package service

import (
	"os"
	"testing"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var imageManager *ImageManager

func setup() {
	opts := zap.Options{
		Development: true,
		TimeEncoder: zapcore.ISO8601TimeEncoder,
	}
	logf.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	imageManager = NewImageManager(logf.Log)
}

func TestImagemapFromVersion(t *testing.T) {
	setup()

	// OLM + RedhatCertified
	os.Setenv("ENTANDO_K8S_OPERATOR_DEPLOYMENT_TYPE", "OLM")
	cr := &v1alpha1.EntandoAppV2{Spec: v1alpha1.EntandoAppV2Spec{Version: "7.1.1", ImageSetType: "RedhatCertified"}}
	images, _ := imageManager.FetchAndComposeImagesMap(cr)

	//fmt.Println(fmt.Sprintf("\n\n%+v\n\n", images.images))

	require.Equal(t, true, utils.IsOlmInstallation())
	require.Equal(t, false, utils.IsImageSetTypeCommunity(cr))

	require.Equal(t,
		"registry.hub.docker.com/entando/app-builder@sha256:33ea636090352a919735aa44cc2aaf2c79e8cb15b19216574964ec41c98f5c58",
		images.FetchAppBuilder())

	require.Equal(t,
		"registry.hub.docker.com/entando/entando-de-app-eap@sha256:21c87b5f0069e38b864211427a39d6d135205d81e24371ddc987f50faa0c21b0",
		images.FetchDeApp())

	require.Equal(t,
		"registry.hub.docker.com/entando/entando-redhat-sso@sha256:b5afae1e933d432ccab84e31a9140f8fd2a51517b95a3a373a29bbe88a62d900",
		images.FetchKeycloak())

	// Plain + Community
	os.Setenv("ENTANDO_K8S_OPERATOR_DEPLOYMENT_TYPE", "Plain")

	cr = &v1alpha1.EntandoAppV2{Spec: v1alpha1.EntandoAppV2Spec{Version: "7.1.1", ImageSetType: "Community"}}
	images, _ = imageManager.FetchAndComposeImagesMap(cr)

	//fmt.Println(fmt.Sprintf("\n\n%+v\n\n", images.images))

	require.Equal(t, false, utils.IsOlmInstallation())
	require.Equal(t, true, utils.IsImageSetTypeCommunity(cr))

	require.Equal(t,
		"registry.hub.docker.com/entando/app-builder:7.1.1",
		images.FetchAppBuilder())

	require.Equal(t,
		"registry.hub.docker.com/entando/entando-de-app-wildfly:7.1.1",
		images.FetchDeApp())

	require.Equal(t,
		"registry.hub.docker.com/entando/entando-keycloak:7.1.1",
		images.FetchKeycloak())

}

func TestImagemapFromVersionAndOverride(t *testing.T) {
	setup()

	// OLM + RedhatCertified
	os.Setenv("ENTANDO_K8S_OPERATOR_DEPLOYMENT_TYPE", "OLM")

	const (
		appBuilder = "myrepo.com/myorg/app-builder@sha256:33ea636090352a919735aa44cc2aaf2c79e8cb15b19216574964ec41c98f5c58"
		deApp      = "myrepo.com/myorg/de-app@sha256:33ea636090352a919735aa44cc2aaf2c79e8cb15b19216574964ec41c98f5c58"
		keycloak   = "myrepo.com/myorg/keycloak@sha256:33ea636090352a919735aa44cc2aaf2c79e8cb15b19216574964ec41c98f5c58"
	)

	cr := &v1alpha1.EntandoAppV2{
		Spec: v1alpha1.EntandoAppV2Spec{
			Version:      "7.1.1",
			ImageSetType: "RedhatCertified",
			AppBuilder:   v1alpha1.AppBuilder{EntandoComponent: v1alpha1.EntandoComponent{ImageOverride: appBuilder}},
			DeApp:        v1alpha1.DeApp{EntandoComponent: v1alpha1.EntandoComponent{ImageOverride: deApp}},
			Keycloak:     v1alpha1.Keycloak{EntandoComponent: v1alpha1.EntandoComponent{ImageOverride: keycloak}},
		},
	}
	images, _ := imageManager.FetchAndComposeImagesMap(cr)

	//fmt.Println(fmt.Sprintf("\n\n%+v\n\n", images.images))

	require.Equal(t, true, utils.IsOlmInstallation())
	require.Equal(t, false, utils.IsImageSetTypeCommunity(cr))

	require.Equal(t,
		appBuilder,
		images.FetchAppBuilder())

	require.Equal(t,
		deApp,
		images.FetchDeApp())

	require.Equal(t,
		keycloak,
		images.FetchKeycloak())

}

func TestImagemapFromOverrideAndNonExistentVersion(t *testing.T) {
	// do we need it ?
}

func TestImagemapOlmWithTag(t *testing.T) {
	setup()

	// OLM + RedhatCertified
	os.Setenv("ENTANDO_K8S_OPERATOR_DEPLOYMENT_TYPE", "OLM")
	imageManager := NewImageManager(logf.Log)

	const (
		appBuilder = "myrepo.com/myorg/app-builder:7.1.1"
	)

	cr := &v1alpha1.EntandoAppV2{
		Spec: v1alpha1.EntandoAppV2Spec{
			Version:      "7.1.1",
			ImageSetType: "RedhatCertified",
			AppBuilder:   v1alpha1.AppBuilder{EntandoComponent: v1alpha1.EntandoComponent{ImageOverride: appBuilder}},
		},
	}
	_, err := imageManager.FetchAndComposeImagesMap(cr)

	require.Equal(t, true, err != nil)

}
