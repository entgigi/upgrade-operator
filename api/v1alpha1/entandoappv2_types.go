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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EntandoAppV2Spec defines the desired state of EntandoAppV2
type EntandoAppV2Spec struct {
	// Version is the field used to upgrade version of EntandoApp
	Version string `json:"version"`
	// ImagesOverride is a map usefull to override EntandoApp component images
	ImagesOverride map[string]string `json:"imagesOverride,omitempty"`
}

// EntandoAppV2Status defines the observed state of EntandoAppV2
type EntandoAppV2Status struct {
	ObservedGeneration int64  `json:"observedGeneration,omitempty"`
	Progress           string `json:"progress,omitempty"`
	Patch              string `json:"patch,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EntandoAppV2 is the Schema for the entandoappv2s API
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
