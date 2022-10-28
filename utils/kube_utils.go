package utils

import (
	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// MergeEnvVars merge the environment variable of the starting deployment with the ones specified in the EntandoAppV2 CR
func MergeEnvVars(entandoAppV2 v1alpha1.EntandoAppV2, deployment *appsv1.Deployment) []corev1.EnvVar {

	envVarMap := ConvertEnvVarSliceToMap(deployment.Spec.Template.Spec.Containers[0].Env)

	for _, e := range entandoAppV2.Spec.CommonEnvironmentVariables {
		envVarMap[e.Name] = e
	}

	for _, e := range entandoAppV2.Spec.ComponentManager.EnvironmentVariables {
		envVarMap[e.Name] = e
	}

	return ConvertEnvVarMapToSlice(envVarMap)
}
