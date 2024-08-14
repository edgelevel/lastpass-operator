package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LastPassGroupSpec defines the desired state of LastPassGroup
type LastPassGroupSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	GroupRef   GroupRef   `json:"groupRef,required"`
	SyncPolicy SyncPolicy `json:"syncPolicy,omitempty"`
}

type GroupRef struct {
	Group        string `json:"group,omitempty"`
	WithUsername bool   `json:"withUsername,omitempty"`
	WithPassword bool   `json:"withPassword,omitempty"`
	WithUrl      bool   `json:"withUrl,omitempty"`
	WithNote     bool   `json:"withNote,omitempty"`
}

// LastPassGroupStatus defines the observed state of LastPassGroup
type LastPassGroupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LastPassGroup is the Schema for the lastpassgroups API
type LastPassGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LastPassGroupSpec   `json:"spec,omitempty"`
	Status LastPassGroupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LastPassGroupList contains a list of LastPassGroup
type LastPassGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LastPassGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LastPassGroup{}, &LastPassGroupList{})
}
