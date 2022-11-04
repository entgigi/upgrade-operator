package utils

import (
	"fmt"
	"os"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/entgigi/upgrade-operator.git/common"
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

	ns, found := os.LookupEnv(common.WatchNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", common.WatchNamespaceEnvVar)
	}
	return ns, nil
}

func IsOlmInstallation() bool {
	operatorType := GetOperatorDeploymentType()
	if operatorType == common.OperatorTypeOlm {
		return true
	}
	return false
}

func IsImageSetTypeCommunity(cr *v1alpha1.EntandoAppV2) bool {
	return cr.Spec.ImageSetType == common.ImageSetTypeCommunity
}

func GetOperatorDeploymentType() string {
	operatorType, found := os.LookupEnv(common.OperatorTypeEnvVar)
	if found {
		return operatorType
	} else {
		// default
		return common.OperatorTypeOlm
	}
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
