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

package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DemoSpec defines the desired state of Demo
type DemoSpec struct {
	Replicas int32              `json:"replicas,omitempty"` // 声明 Pod 副本数
	Template v1.PodTemplateSpec `json:"template,omitempty"` // 声明 Pod 模板配置
}

// DemoStatus defines the observed state of Demo
type DemoStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Demo is the Schema for the demoes API
type Demo struct {
	metav1.TypeMeta   `json:",inline"`            // API 元数据
	metav1.ObjectMeta `json:"metadata,omitempty"` // 对象元数据

	Spec   DemoSpec   `json:"spec,omitempty"` // 自定义部分
	Status DemoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DemoList contains a list of Demo
// Kubernetes 获取对象的 List() 方法，返回 List 类型，而非对象类型的数组
type DemoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Demo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Demo{}, &DemoList{}) // 向 APIServer 注册类型
}
