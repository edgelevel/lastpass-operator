/*


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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LastPassSpec defines the desired state of LastPass
// +k8s:openapi-gen=true
// +kubebuilder:resource:path=lastpass,scope=Namespaced
type LastPassSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	SecretRef  SecretRef  `json:"secretRef,required"`
	SyncPolicy SyncPolicy `json:"syncPolicy,omitempty"`
}

type SecretRef struct {
	Group        string `json:"group,omitempty"`
	Name         string `json:"name,required"`
	WithUsername bool   `json:"withUsername,omitempty"`
	WithPassword bool   `json:"withPassword,omitempty"`
	WithUrl      bool   `json:"withUrl,omitempty"`
	WithNote     bool   `json:"withNote,omitempty"`
}

// LastPassStatus defines the observed state of LastPass
// +k8s:openapi-gen=true
type LastPassStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LastPass is the Schema for the lastpasses API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type LastPass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LastPassSpec   `json:"spec,omitempty"`
	Status LastPassStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true

// LastPassList contains a list of LastPass
type LastPassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LastPass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LastPass{}, &LastPassList{})
}
