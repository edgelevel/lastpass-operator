package v1alpha1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LastPassSpec defines the desired state of LastPass
// +k8s:openapi-gen=true
// +kubebuilder:resource:path=lastpass,scope=Namespaced
type LastPassSpec struct {
	// SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	SecretRef  SecretRef  `json:"secretRef,required"`
	SyncPolicy SyncPolicy `json:"syncPolicy"`
}

type SecretRef struct {
	Group        string `json:"group,omitempty"`
	Name         string `json:"name,required"`
	WithUsername bool   `json:"withUsername,omitempty"`
	WithPassword bool   `json:"withPassword,omitempty"`
	WithUrl      bool   `json:"withUrl,omitempty"`
	WithNote     bool   `json:"withNote,omitempty"`
}

type SyncPolicy struct {
	Enabled bool          `json:"enabled,required"`
	Refresh time.Duration `json:"refresh,required"`
}

// LastPassStatus defines the observed state of LastPass
// +k8s:openapi-gen=true
type LastPassStatus struct {
	// STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

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

// LastPassList contains a list of LastPass
type LastPassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LastPass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LastPass{}, &LastPassList{})
}
