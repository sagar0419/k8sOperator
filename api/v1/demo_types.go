/*
Copyright 2023.

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

// package v1 line refers to the package name in which the API type of this version of the custom resource are defined.
package v1

// package that are required to work with kubernetes API
import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// This spec struct is used to define the specification or desired state of the resource
type ServiceSpec struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Protocol   string `json:"protocol"`
	Port       int32  `json:"port"`
	TargetPort int32  `json:"targetPort"`
	NodePort   int32  `json:"nodePort"`
	Type       string `json:"type"`
}

// DemoSpec defines the desired state of Demo
type DemoSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Demo. Edit demo_types.go to remove/update
	// Foo string `json:"foo,omitempty"`
	CompanyName            string      `json:"companyName,omitempty"`
	ApplicationDescription string      `json:"applicationDescription,omitempty"`
	AppContainerName       string      `json:"appContainerName"`
	AppImage               string      `json:"appImage"`
	AppPort                int32       `json:"appPort"`
	MonitorContainerName   string      `json:"monitorContainerName"`
	MonitorImage           string      `json:"monitorImage"`
	MonitorCommand         string      `json:"monitorCommand"`
	Size                   int32       `json:"size"`
	Service                ServiceSpec `json:"service"`
}

// DemoStatus defines the observed state/status of Demo
type DemoStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	PodList []string `json:"podList"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Demo is the Schema for the demoes API
type Demo struct {
	metav1.TypeMeta `json:",inline"`
	// metav1.ObjectMeta: This field includes metadata about the resource instance itself. It contains fields like name, namespace, labels, annotations, and more.
	// The json:"metadata,omitempty" tag indicates that this field is represented in the JSON serialization of the custom resource as "metadata", and it's optional (the omitempty part).
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Spec: This field of type DemoSpec is used to define the specification of the custom resource.
	Spec DemoSpec `json:"spec,omitempty"`
	// Status: This field of type DemoStatus is used to define the status of the custom resource.
	Status DemoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DemoList contains a list of Demo
type DemoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Demo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Demo{}, &DemoList{})
}
