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

// EntandoUpgradeSpec defines the desired state of EntandoUpgrade
type EntandoUpgradeSpec struct {

	// FromVersion is an example field of EntandoUpgrade. Edit entandoupgrade_types.go to remove/update
	FromVersion string `json:"fromVersion"`
	// ToVersion is an example field of EntandoUpgrade. Edit entandoupgrade_types.go to remove/update
	ToVersion string `json:"toVersion"`
}

// EntandoUpgradeStatus defines the observed state of EntandoUpgrade
type EntandoUpgradeStatus struct {
	Progress string `json:"progress,omitempty"`
	Patch    string `json:"patch,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EntandoUpgrade is the Schema for the entandoupgrades API
type EntandoUpgrade struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EntandoUpgradeSpec   `json:"spec,omitempty"`
	Status EntandoUpgradeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EntandoUpgradeList contains a list of EntandoUpgrade
type EntandoUpgradeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EntandoUpgrade `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EntandoUpgrade{}, &EntandoUpgradeList{})
}
