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

package v1alpha1

// SocketApplyConfiguration represents an declarative configuration of the Socket type for use
// with apply.
type SocketApplyConfiguration struct {
	Cores    map[string]CoreApplyConfiguration `json:"cores,omitempty"`
	ListCpus []int                             `json:"listCpus,omitempty"`
}

// SocketApplyConfiguration constructs an declarative configuration of the Socket type for use with
// apply.
func Socket() *SocketApplyConfiguration {
	return &SocketApplyConfiguration{}
}

// WithCores puts the entries into the Cores field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Cores field,
// overwriting an existing map entries in Cores field with the same key.
func (b *SocketApplyConfiguration) WithCores(entries map[string]CoreApplyConfiguration) *SocketApplyConfiguration {
	if b.Cores == nil && len(entries) > 0 {
		b.Cores = make(map[string]CoreApplyConfiguration, len(entries))
	}
	for k, v := range entries {
		b.Cores[k] = v
	}
	return b
}

// WithListCpus adds the given value to the ListCpus field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the ListCpus field.
func (b *SocketApplyConfiguration) WithListCpus(values ...int) *SocketApplyConfiguration {
	for i := range values {
		b.ListCpus = append(b.ListCpus, values[i])
	}
	return b
}
