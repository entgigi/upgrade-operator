package utils

import (
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
)

func CopyMap(m map[string]interface{}) map[string]interface{} {
	cp := make(map[string]interface{})
	for k, v := range m {
		vm, ok := v.(map[string]interface{})
		if ok {
			cp[k] = CopyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

// GetWatchNamespace returns the Namespace the operator should be watching for changes
func GetWatchNamespace() (string, error) {
	// WatchNamespaceEnvVar is the constant for env variable WATCH_NAMESPACE
	// which specifies the Namespace to watch.
	// An empty value means the operator is running with cluster scope.
	var watchNamespaceEnvVar = "WATCH_NAMESPACE"

	ns, found := os.LookupEnv(watchNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", watchNamespaceEnvVar)
	}
	return ns, nil
}

func IsOlmInstallation() bool {
	var operatorTypeEnvVar = "ENTANDO_K8S_OPERATOR_DEPLOYMENT_TYPE"

	operatorType, found := os.LookupEnv(operatorTypeEnvVar)
	if !found || operatorType == "olm" {
		return true
	}
	return false
}

// TODO make more generic with interface parameter

func ConvertEnvVarSliceToMap(src []corev1.EnvVar) map[string]corev1.EnvVar {
	elementMap := make(map[string]corev1.EnvVar)
	for _, item := range src {
		elementMap[item.Name] = item
	}
	return elementMap
}

func ConvertEnvVarMapToSlice(src map[string]corev1.EnvVar) []corev1.EnvVar {
	elementSlice := make([]corev1.EnvVar, 0)
	for _, item := range src {
		elementSlice = append(elementSlice, item)
	}
	return elementSlice
}
