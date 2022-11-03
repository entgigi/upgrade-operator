package common

const (
	AppName                    string = "EntandoAppV2"
	Version                    string = "0.0.1"
	ImageInfoConfigMap         string = "entando-docker-image-info"
	AppBuilderKey              string = "AppBuilder"
	ComponentManagerKey        string = "ComponentManager"
	DeAppKey                   string = "DeApp"
	DeAppEapKey                string = "DeAppEap"
	KeycloakKey                string = "Keycloak"
	KeycloakSsoKey             string = "KeycloakSso"
	K8sService                 string = "K8sService"
	K8sPluginController        string = "K8sPluginController"
	K8sAppPluginLinkController string = "K8sAppPluginLinkController"

	WatchNamespaceEnvVar  string = "WATCH_NAMESPACE"
	OperatorTypeEnvVar    string = "ENTANDO_K8S_OPERATOR_DEPLOYMENT_TYPE"
	OperatorTypeOlm       string = "olm"
	OperatorTypeCommunity string = "community"
)
