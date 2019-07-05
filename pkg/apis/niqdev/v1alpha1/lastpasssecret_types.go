package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LastPassSecretSpec defines the desired state of LastPassSecret
// +k8s:openapi-gen=true
type LastPassSecretSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	ItemRef ItemRef `json:"itemRef,required"`
}

type ItemRef struct {
	Group        string `json:"group,omitempty"`
	Name         string `json:"name,required"`
	WithUsername bool   `json:"withUsername,omitempty"`
	WithPassword bool   `json:"withPassword,omitempty"`
	WithUrl      bool   `json:"withUrl,omitempty"`
	WithNote     bool   `json:"withNote,omitempty"`
}

// LastPassSecretStatus defines the observed state of LastPassSecret
// +k8s:openapi-gen=true
type LastPassSecretStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LastPassSecret is the Schema for the lastpasssecrets API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type LastPassSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LastPassSecretSpec   `json:"spec,omitempty"`
	Status LastPassSecretStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LastPassSecretList contains a list of LastPassSecret
type LastPassSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LastPassSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LastPassSecret{}, &LastPassSecretList{})
}
