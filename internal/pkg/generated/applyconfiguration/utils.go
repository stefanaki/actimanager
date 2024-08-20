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
// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	v1alpha1 "cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/internal/pkg/generated/applyconfiguration/cslab.ece.ntua.gr/v1alpha1"
	internal "cslab.ece.ntua.gr/actimanager/internal/pkg/generated/applyconfiguration/internal"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=cslab.ece.ntua.gr, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("Core"):
		return &cslabecentuagrv1alpha1.CoreApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("CPU"):
		return &cslabecentuagrv1alpha1.CPUApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("CPUTopology"):
		return &cslabecentuagrv1alpha1.CPUTopologyApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("NodeCPUTopology"):
		return &cslabecentuagrv1alpha1.NodeCPUTopologyApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("NodeCPUTopologySpec"):
		return &cslabecentuagrv1alpha1.NodeCPUTopologySpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("NodeCPUTopologyStatus"):
		return &cslabecentuagrv1alpha1.NodeCPUTopologyStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("NUMANode"):
		return &cslabecentuagrv1alpha1.NUMANodeApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PodCPUBinding"):
		return &cslabecentuagrv1alpha1.PodCPUBindingApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PodCPUBindingSpec"):
		return &cslabecentuagrv1alpha1.PodCPUBindingSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PodCPUBindingStatus"):
		return &cslabecentuagrv1alpha1.PodCPUBindingStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Socket"):
		return &cslabecentuagrv1alpha1.SocketApplyConfiguration{}

	}
	return nil
}

func NewTypeConverter(scheme *runtime.Scheme) *testing.TypeConverter {
	return &testing.TypeConverter{Scheme: scheme, TypeResolver: internal.Parser()}
}
