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
// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1alpha1 "cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/internal/pkg/gen/applyconfiguration/cslab.ece.ntua.gr/v1alpha1"
	scheme "cslab.ece.ntua.gr/actimanager/internal/pkg/gen/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PodCpuBindingsGetter has a method to return a PodCpuBindingInterface.
// A group's client should implement this interface.
type PodCpuBindingsGetter interface {
	PodCpuBindings(namespace string) PodCpuBindingInterface
}

// PodCpuBindingInterface has methods to work with PodCpuBinding resources.
type PodCpuBindingInterface interface {
	Create(ctx context.Context, podCpuBinding *v1alpha1.PodCpuBinding, opts v1.CreateOptions) (*v1alpha1.PodCpuBinding, error)
	Update(ctx context.Context, podCpuBinding *v1alpha1.PodCpuBinding, opts v1.UpdateOptions) (*v1alpha1.PodCpuBinding, error)
	UpdateStatus(ctx context.Context, podCpuBinding *v1alpha1.PodCpuBinding, opts v1.UpdateOptions) (*v1alpha1.PodCpuBinding, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.PodCpuBinding, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.PodCpuBindingList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.PodCpuBinding, err error)
	Apply(ctx context.Context, podCpuBinding *cslabecentuagrv1alpha1.PodCpuBindingApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.PodCpuBinding, err error)
	ApplyStatus(ctx context.Context, podCpuBinding *cslabecentuagrv1alpha1.PodCpuBindingApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.PodCpuBinding, err error)
	PodCpuBindingExpansion
}

// podCpuBindings implements PodCpuBindingInterface
type podCpuBindings struct {
	client rest.Interface
	ns     string
}

// newPodCpuBindings returns a PodCpuBindings
func newPodCpuBindings(c *CslabV1alpha1Client, namespace string) *podCpuBindings {
	return &podCpuBindings{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the podCpuBinding, and returns the corresponding podCpuBinding object, and an error if there is any.
func (c *podCpuBindings) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.PodCpuBinding, err error) {
	result = &v1alpha1.PodCpuBinding{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("podcpubindings").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PodCpuBindings that match those selectors.
func (c *podCpuBindings) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.PodCpuBindingList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.PodCpuBindingList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("podcpubindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested podCpuBindings.
func (c *podCpuBindings) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("podcpubindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a podCpuBinding and creates it.  Returns the server's representation of the podCpuBinding, and an error, if there is any.
func (c *podCpuBindings) Create(ctx context.Context, podCpuBinding *v1alpha1.PodCpuBinding, opts v1.CreateOptions) (result *v1alpha1.PodCpuBinding, err error) {
	result = &v1alpha1.PodCpuBinding{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("podcpubindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(podCpuBinding).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a podCpuBinding and updates it. Returns the server's representation of the podCpuBinding, and an error, if there is any.
func (c *podCpuBindings) Update(ctx context.Context, podCpuBinding *v1alpha1.PodCpuBinding, opts v1.UpdateOptions) (result *v1alpha1.PodCpuBinding, err error) {
	result = &v1alpha1.PodCpuBinding{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("podcpubindings").
		Name(podCpuBinding.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(podCpuBinding).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *podCpuBindings) UpdateStatus(ctx context.Context, podCpuBinding *v1alpha1.PodCpuBinding, opts v1.UpdateOptions) (result *v1alpha1.PodCpuBinding, err error) {
	result = &v1alpha1.PodCpuBinding{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("podcpubindings").
		Name(podCpuBinding.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(podCpuBinding).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the podCpuBinding and deletes it. Returns an error if one occurs.
func (c *podCpuBindings) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("podcpubindings").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *podCpuBindings) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("podcpubindings").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched podCpuBinding.
func (c *podCpuBindings) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.PodCpuBinding, err error) {
	result = &v1alpha1.PodCpuBinding{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("podcpubindings").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied podCpuBinding.
func (c *podCpuBindings) Apply(ctx context.Context, podCpuBinding *cslabecentuagrv1alpha1.PodCpuBindingApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.PodCpuBinding, err error) {
	if podCpuBinding == nil {
		return nil, fmt.Errorf("podCpuBinding provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(podCpuBinding)
	if err != nil {
		return nil, err
	}
	name := podCpuBinding.Name
	if name == nil {
		return nil, fmt.Errorf("podCpuBinding.Name must be provided to Apply")
	}
	result = &v1alpha1.PodCpuBinding{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("podcpubindings").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *podCpuBindings) ApplyStatus(ctx context.Context, podCpuBinding *cslabecentuagrv1alpha1.PodCpuBindingApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.PodCpuBinding, err error) {
	if podCpuBinding == nil {
		return nil, fmt.Errorf("podCpuBinding provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(podCpuBinding)
	if err != nil {
		return nil, err
	}

	name := podCpuBinding.Name
	if name == nil {
		return nil, fmt.Errorf("podCpuBinding.Name must be provided to Apply")
	}

	result = &v1alpha1.PodCpuBinding{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("podcpubindings").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}