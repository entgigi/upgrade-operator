/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EntandoAppV2Spec defines the desired state of EntandoAppV2
type EntandoAppV2Spec struct {
	// Version is the field used to upgrade version of EntandoApp
	Version string `json:"version"`
	// EntandoAppName is the name of the EntandoApp CR to update
	EntandoAppName string `json:"entandoAppName,omitempty"`
	// ImageSetType is used to select the list of image to use, the admitted values are: RedhatCertified or Community
	ImageSetType string `json:"imageSetType,omitempty"`
	// IngressHostName is the host name that Ingress uses to control access to the app
	IngressHostName string `json:"ingressHostName,omitempty"`
	// UpdateStrategy the strategy to apply during the upgrade of the deployments, the admitted values are: Recreate or RollingUpdate
	// Recreate: deployments will be updated and scaled to 0, then they will be scaled up to 1 again
	// RollingUpdate: deployments will be updated without scaling down to 0, trying to ensure 0 downtime
	UpdateStrategy string `json:"updateStrategy,omitempty"`
	// used to add Environment Variables to all EntandoApp components
	CommonEnvironmentVariables []corev1.EnvVar `json:"commonEnvironmentVariables,omitempty"`
	// Section used to configure AppBuilder
	AppBuilder AppBuilder `json:"appBuilder,omitempty"`
	// Section used to configure ComponentManager
	ComponentManager ComponentManager `json:"componentManager,omitempty"`
	// Section used to configure DeApp
	DeApp DeApp `json:"deApp,omitempty"`
	// Section used to configure Keycloak
	Keycloak Keycloak `json:"keycloak,omitempty"`
	// Section used to configure k8s-service
	K8sService K8sService `json:"k8sService,omitempty"`
	// Section used to configure legacy operator plugin-controller
	K8sPluginController K8sPluginController `json:"k8sPluginController,omitempty"`
	// Section used to configure legacy operator app-plugin-link-controller
	K8sAppPluginLinkController K8sAppPluginLinkController `json:"k8sAppPluginLinkController,omitempty"`
}
type EntandoComponent struct {
	// used to override the component image
	ImageOverride string `json:"imageOverride,omitempty"`
	// used to add Environment Variables to component
	EnvironmentVariables []corev1.EnvVar `json:"environmentVariables,omitempty"`
}

type AppBuilder struct {
	// Empty JSON tag is needed to avoid the error 'encountered struct field "" without JSON tag'
	EntandoComponent `json:",omitempty"`
}

type ComponentManager struct {
	// Empty JSON tag is needed to avoid the error 'encountered struct field "" without JSON tag'
	EntandoComponent `json:",omitempty"`
}

type DeApp struct {
	// Empty JSON tag is needed to avoid the error 'encountered struct field "" without JSON tag'
	EntandoComponent `json:",omitempty"`
}

type Keycloak struct {
	// Empty JSON tag is needed to avoid the error 'encountered struct field "" without JSON tag'
	EntandoComponent `json:",omitempty"`
	ExternalService  bool `json:"externalService,omitempty"`
}

type K8sService struct {
	EntandoComponent `json:",omitempty"`
}

type K8sPluginController struct {
	EntandoComponent `json:",omitempty"`
}

type K8sAppPluginLinkController struct {
	EntandoComponent `json:",omitempty"`
}

// EntandoAppV2Status defines the observed state of EntandoAppV2
type EntandoAppV2Status struct {
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
	Progress           int                `json:"progress,omitempty"`
	Total              int                `json:"total,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// An EntandoAppV2 deploys the components required to upgrade an Entando App. The server side
// components that are deployed include the Entando App Engine, the Entando Component Manager,
// the Entando App Builder, and the user facing application.
type EntandoAppV2 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EntandoAppV2Spec   `json:"spec,omitempty"`
	Status EntandoAppV2Status `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EntandoAppV2List contains a list of EntandoAppV2
type EntandoAppV2List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EntandoAppV2 `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EntandoAppV2{}, &EntandoAppV2List{})
}
