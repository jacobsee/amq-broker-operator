package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AMQBrokerSpec defines the desired state of AMQBroker
type AMQBrokerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	Username string `json:"username"`
	Password string `json:"password"`
}

// AMQBrokerStatus defines the observed state of AMQBroker
type AMQBrokerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AMQBroker is the Schema for the amqbrokers API
// +k8s:openapi-gen=true
type AMQBroker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AMQBrokerSpec   `json:"spec,omitempty"`
	Status AMQBrokerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AMQBrokerList contains a list of AMQBroker
type AMQBrokerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AMQBroker `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AMQBroker{}, &AMQBrokerList{})
}
