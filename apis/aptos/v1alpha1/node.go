package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AptosNetwork is Aptos network
type AptosNetwork string

const (
	Devnet  AptosNetwork = "devnet"
	Testnet AptosNetwork = "testnet"
)

// NodeSpec defines the desired state of Node
type NodeSpec struct {
	// Image is Aptos node client image
	Image *string `json:"image,omitempty"`
	// Network is Aptos network to join and sync
	// +kubebuilder:validation:Enum=devnet;testnet
	Network AptosNetwork `json:"network"`
	// Validator enables validator mode
	Validator bool `json:"validator,omitempty"`
	// Waypoint provides an off-chain mechanism to verify the sync process after restart or epoch change
	Waypoint string `json:"waypoint"`
	// GenesisConfigmapName is Kubernetes configmap name holding genesis blob
	GenesisConfigmapName string `json:"genesisConfigmapName"`
	// NodePrivateKeySecretName is the secret name holding node private key
	NodePrivateKeySecretName string `json:"nodePrivateKeySecretName,omitempty"`
	// PeerId is the node identity
	PeerId string `json:"peerId,omitempty"`
}

// NodeStatus defines the observed state of Node
type NodeStatus struct {
}

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