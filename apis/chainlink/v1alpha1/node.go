package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// APICredentials is api credentials
type APICredentials struct {
	// Email is user email
	Email string `json:"email"`
	// PasswordSecretName is the k8s secret name that holds password
	PasswordSecretName string `json:"passwordSecretName"`
}

// NodeSpec defines the desired state of Node
type NodeSpec struct {
	// EthereumChainId is ethereum chain id
	EthereumChainId uint `json:"ethereumChainId"`
	// EthereumWSEndpoint is ethereum websocket endpoint
	EthereumWSEndpoint string `json:"ethereumWsEndpoint"`
	// LinkContractAddress is link contract address
	LinkContractAddress string `json:"linkContractAddress"`
	// DatabaseURL is postgres database connection URL
	DatabaseURL string `json:"databaseURL"`
	// KeystorePasswordSecretName is k8s secret name that holds keystore password
	KeystorePasswordSecretName string `json:"keystorePasswordSecretName"`
	// APICredentials is api credentials
	APICredentials APICredentials `json:"apiCredentials"`
}

// NodeStatus defines the observed state of Node
type NodeStatus struct{}

// +kubebuilder:object:root=true

// Node is the Schema for the nodes API
type Node struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeSpec   `json:"spec,omitempty"`
	Status NodeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NodeList contains a list of Node
type NodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Node `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Node{}, &NodeList{})
}