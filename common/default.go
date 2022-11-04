package common

const (
	AppName            string = "EntandoAppV2"
	Version            string = "0.0.1"
	ImageInfoConfigMap string = "entando-docker-image-info"

	// images internal table keys
	TagKey                           string = "Tag"
	AppBuilderKey                    string = "AppBuilder"
	AppBuilderKeyTag                 string = "AppBuilder" + TagKey
	ComponentManagerKey              string = "ComponentManager"
	ComponentManagerKeyTag           string = "ComponentManager" + TagKey
	DeAppKey                         string = "DeApp"
	DeAppKeyTag                      string = "DeApp" + TagKey
	DeAppEapKey                      string = "DeAppEap"
	DeAppEapKeyTag                   string = "DeAppEap" + TagKey
	KeycloakKey                      string = "Keycloak"
	KeycloakKeyTag                   string = "Keycloak" + TagKey
	KeycloakSsoKey                   string = "KeycloakSso"
	KeycloakSsoKeyTag                string = "KeycloakSso" + TagKey
	K8sServiceKey                    string = "K8sService"
	K8sServiceKeyTag                 string = "K8sService" + TagKey
	K8sPluginControllerKey           string = "K8sPluginController"
	K8sPluginControllerKeyTag        string = "K8sPluginController" + TagKey
	K8sAppPluginLinkControllerKey    string = "K8sAppPluginLinkController"
	K8sAppPluginLinkControllerkeyTag string = "K8sAppPluginLinkController" + TagKey

	// Operator configuration constants
	WatchNamespaceEnvVar        string = "WATCH_NAMESPACE"
	OperatorTypeEnvVar          string = "ENTANDO_K8S_OPERATOR_DEPLOYMENT_TYPE"
	OperatorTypeOlm             string = "OLM"
	OperatorTypeCommunity       string = "Plain"
	ImageSetTypeRedHatCertified string = "Community"
	ImageSetTypeCommunity       string = "RedhatCertified"
)
