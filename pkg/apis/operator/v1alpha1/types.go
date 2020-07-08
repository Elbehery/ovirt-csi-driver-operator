package v1alpha1

import (
	openshiftapi "github.com/openshift/api/operator/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OvirtCSIOperatorSpec defines the desired state of OvirtCSIOperator
type OvirtCSIOperatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// managementState indicates whether and how the operator should manage the component
	ManagementState openshiftapi.ManagementState `json:"managementState"`

}

// OvirtCSIOperatorStatus defines the observed state of OvirtCSIOperator
type OvirtCSIOperatorStatus struct {
	openshiftapi.OperatorStatus `json:",inline"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Generation of API objects that the operator has created / updated.
	// For internal operator bookkeeping purposes.
	Children []openshiftapi.GenerationHistory `json:"children,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OvirtCSIOperator is the Schema for the ovirtcsioperators API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=ovirtcsioperators,scope=Namespaced
type OvirtCSIOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OvirtCSIOperatorSpec   `json:"spec,omitempty"`
	Status OvirtCSIOperatorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OvirtCSIOperatorList contains a list of OvirtCSIOperator
type OvirtCSIOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OvirtCSIOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OvirtCSIOperator{}, &OvirtCSIOperatorList{})
}

// CSIDeploymentContainerImages specifies custom sidecar container image names. This should be used only to override the default operator image names.
type CSIDeploymentContainerImages struct {
	// Name of CSI Attacher sidecar container image.
	// Optional.
	AttacherImage *string `json:"attacherImage,omitempty" yaml:"attacherImage"` // Use yaml tags for reading config from a file.

	// Name of CSI Provisioner sidecar container image.
	// Optional.
	ProvisionerImage *string `json:"provisionerImage,omitempty" yaml:"provisionerImage"`

	// Name of CSI Driver Registrar sidecar container image.
	// Optional.
	DriverRegistrarImage *string `json:"driverRegistrarImage,omitempty" yaml:"driverRegistrarImage"`

	// Name of CSI Liveness Probe sidecar container image.
	// Optional.
	LivenessProbeImage *string `json:"livenessProbeImage,omitempty" yaml:"livenessProbeImage"`
}
