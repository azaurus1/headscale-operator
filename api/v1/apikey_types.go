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

// ApiKeySpec defines the desired state of ApiKey
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="ID",type="string",JSONPath=".status.ID",description="ID"
type ApiKeySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	TimeToExpire int32 `json:"timeToExpire"` //When we initialise an API key, we will use this to create the Expiration

}

// ApiKeyStatus defines the observed state of ApiKey
type ApiKeyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	KeySecret string      `json:"keysecret"` // this is the API_KEY, we could possibly try put this into a Secret
	CreatedAt metav1.Time `json:"createdAt"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ApiKey is the Schema for the apikeys API
type ApiKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiKeySpec   `json:"spec,omitempty"`
	Status ApiKeyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApiKeyList contains a list of ApiKey
type ApiKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiKey `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApiKey{}, &ApiKeyList{})
}
