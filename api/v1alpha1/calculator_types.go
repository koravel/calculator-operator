/*
Copyright 2022.

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

// CalculatorSpec defines the desired state of Calculator
type CalculatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Size defines the number of Calculator instances
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=3
	// +kubebuilder:validation:ExclusiveMaximum=false
	Size int32 `json:"size,omitempty"`

	X int32 `json:"x,omitempty"`
	Y int32 `json:"y,omitempty"`
}

// CalculatorStatus defines the observed state of Calculator
type CalculatorStatus struct {
	// Represents the observations of a Calculator's current state.
	// Calculator.status.conditions.type are: "Available", "Progressing", and "Degraded"
	// Calculator.status.conditions.status are one of True, False, Unknown.
	// Calculator.status.conditions.reason the value should be a CamelCase string and producers of specific
	// condition types may define expected values and meanings for this field, and whether the values
	// are considered a guaranteed API.
	// Calculator.status.conditions.Message is a human readable message indicating details about the transition.
	// For further information see: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
	Result     int32              `json:"sum,omitempty" managed-by:"calc-operator"`
	Processed  bool               `json:"processed,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Calculator is the Schema for the calculators API
type Calculator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CalculatorSpec   `json:"spec,omitempty"`
	Status CalculatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CalculatorList contains a list of Calculator
type CalculatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Calculator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Calculator{}, &CalculatorList{})
}
