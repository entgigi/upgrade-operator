package utils

import (
	"fmt"
	"testing"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestShouldConvertAnEnvVarSliceToMap(t *testing.T) {

	appCommonEnvVar1 := corev1.EnvVar{
		Name:  "common_name_1",
		Value: "SUCCESS",
	}
	appCommonEnvVar2 := corev1.EnvVar{
		Name:  "common_name_2",
		Value: "FAIL",
	}
	appCommonEnvVar3 := corev1.EnvVar{
		Name:  "double_update",
		Value: "FAIL_2",
	}
	appCommonEnvVars := []corev1.EnvVar{appCommonEnvVar1, appCommonEnvVar2, appCommonEnvVar3}

	appCMEnvVar1 := corev1.EnvVar{
		Name:  "cm_var_name_1",
		Value: "SUCCESS",
	}
	appCMEnvVar2 := corev1.EnvVar{
		Name:  "common_name_2",
		Value: "SUCCESS",
	}
	appCMEnvVar3 := corev1.EnvVar{
		Name:  "double_update",
		Value: "SUCCESS",
	}
	appCMEnvVars := []corev1.EnvVar{appCMEnvVar1, appCMEnvVar2, appCMEnvVar3}

	componentManager := v1alpha1.ComponentManager{}
	componentManager.EnvironmentVariables = appCMEnvVars

	appV2 := v1alpha1.EntandoAppV2{
		Spec: v1alpha1.EntandoAppV2Spec{
			CommonEnvironmentVariables: appCommonEnvVars,
			ComponentManager:           componentManager,
		},
	}

	deplEnvVar0 := corev1.EnvVar{
		Name:  "var_name_1",
		Value: "SUCCESS",
	}
	deplEnvVar1 := corev1.EnvVar{
		Name:  "common_name_1",
		Value: "FAIL",
	}
	deplEnvVar2 := corev1.EnvVar{
		Name:  "cm_var_name_1",
		Value: "FAIL",
	}
	deplEnvVar3 := corev1.EnvVar{
		Name:  "double_update",
		Value: "FAIL",
	}
	deplEnvVars := []corev1.EnvVar{deplEnvVar0, deplEnvVar1, deplEnvVar2, deplEnvVar3}

	cmDepl := v1.Deployment{
		Spec: v1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Env: deplEnvVars,
						},
					},
				},
			},
		},
	}

	vars := MergeEnvVars(appV2, &cmDepl)
	for _, a := range vars {
		fmt.Println("### " + a.Name + " --- " + a.Value)
	}
	assertOnEnvVar(t, vars, appCommonEnvVar1)
	assertOnEnvVar(t, vars, appCommonEnvVar2)
	assertOnEnvVar(t, vars, deplEnvVar0)
	assertOnEnvVar(t, vars, appCMEnvVar1)
	assertOnEnvVar(t, vars, appCMEnvVar3)
	assertOnEnvVar(t, vars, appCommonEnvVar3)

	// FIXME LET IT WORK
}

func assertOnEnvVar(t *testing.T, envVars []corev1.EnvVar, expected corev1.EnvVar) {
	actual := containsEnvVar(envVars, expected)
	assert.NotNil(t, actual)
	assert.Equal(t, "SUCCESS", actual.Value)
}

func containsEnvVar(s []corev1.EnvVar, e corev1.EnvVar) *corev1.EnvVar {
	for _, a := range s {
		if a.Name == e.Name {
			return &a
		}
	}
	return nil
}
