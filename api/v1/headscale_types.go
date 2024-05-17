/*
Copyright 2024.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:resource:scope=Namespaced
// HeadscaleSpec defines the desired state of Headscale
type HeadscaleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Version string `json:"version"`
	Name    string `json:"name"`
	Config  string `json:"config"`
}

// HeadscaleStatus defines the observed state of Headscale
type HeadscaleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Version string `json:"version"`
	Name    string `json:"name"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Headscale is the Schema for the headscales API
type Headscale struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HeadscaleSpec   `json:"spec,omitempty"`
	Status HeadscaleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HeadscaleList contains a list of Headscale
type HeadscaleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Headscale `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Headscale{}, &HeadscaleList{})
}
