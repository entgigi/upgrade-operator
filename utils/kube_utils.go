package utils

import (
	"context"
	"errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const DeploymentLabelKey = "entando.org/deployment"

// MergeEnvVars merge the environment variables of the starting deployment with the ones specified in the EntandoAppV2 CR
func MergeEnvVars(deployment *appsv1.Deployment,
	genericEnvVars []corev1.EnvVar,
	specificEnvVars []corev1.EnvVar) []corev1.EnvVar {

	envVarMap := ConvertEnvVarSliceToMap(deployment.Spec.Template.Spec.Containers[0].Env)

	for _, e := range genericEnvVars {
		envVarMap[e.Name] = e
	}

	for _, e := range specificEnvVars {
		envVarMap[e.Name] = e
	}

	// TODO find a way to automatically inject the env vars required by the new version to update to

	return ConvertEnvVarMapToSlice(envVarMap)
}

// BuildDeploymentLabelSelectorWithAppName build and return the label to select an entando deployment
func BuildDeploymentLabelSelectorWithAppName(appName string, componentName string) map[string]string {
	labelValue := appName
	if componentName != "" {
		labelValue += "-" + componentName
	}

	return BuildDeploymentLabelSelector(labelValue)
}

// BuildDeploymentLabelSelector build and return the label to select an entando deployment
func BuildDeploymentLabelSelector(componentName string) map[string]string {
	return map[string]string{DeploymentLabelKey: componentName}
}

func MustGetFirstDeploymentByLabels(ctx context.Context, kubeClient client.Client, labelsMap map[string]string) (*appsv1.Deployment, error) {

	deployments, err := FindDeploymentsByLabels(ctx, kubeClient, labelsMap)
	if err != nil {
		return nil, err
	}

	if len(deployments.Items) <= 0 {
		return nil, errors.New("no deployment found matching the received labels") // TODO add labels detail?
	}

	if len(deployments.Items) > 1 {
		return nil, errors.New("more than 1 deployment found matching the received labels") // TODO add labels detail?
	}

	return &deployments.Items[0], err
}

// FindDeploymentsByLabels find and return the list of deployment corresponding to the received labels
func FindDeploymentsByLabels(ctx context.Context, kubeClient client.Client, labelsMap map[string]string) (*appsv1.DeploymentList, error) {

	selector, err := labels.ValidatedSelectorFromSet(labelsMap)
	if err != nil {
		return nil, err
	}

	//namespaceOptions := client.InNamespace(namespace)
	labelOptions := client.MatchingLabelsSelector{
		Selector: selector,
	}

	deployments := &appsv1.DeploymentList{}
	if err = kubeClient.List(ctx, deployments, labelOptions); err != nil {
		return nil, err
	}

	return deployments, nil
}
