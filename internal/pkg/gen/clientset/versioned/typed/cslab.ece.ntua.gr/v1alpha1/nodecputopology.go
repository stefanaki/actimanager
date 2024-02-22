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

// NodeCpuTopologiesGetter has a method to return a NodeCpuTopologyInterface.
// A group's client should implement this interface.
type NodeCpuTopologiesGetter interface {
	NodeCpuTopologies() NodeCpuTopologyInterface
}

// NodeCpuTopologyInterface has methods to work with NodeCpuTopology resources.
type NodeCpuTopologyInterface interface {
	Create(ctx context.Context, nodeCpuTopology *v1alpha1.NodeCpuTopology, opts v1.CreateOptions) (*v1alpha1.NodeCpuTopology, error)
	Update(ctx context.Context, nodeCpuTopology *v1alpha1.NodeCpuTopology, opts v1.UpdateOptions) (*v1alpha1.NodeCpuTopology, error)
	UpdateStatus(ctx context.Context, nodeCpuTopology *v1alpha1.NodeCpuTopology, opts v1.UpdateOptions) (*v1alpha1.NodeCpuTopology, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.NodeCpuTopology, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.NodeCpuTopologyList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.NodeCpuTopology, err error)
	Apply(ctx context.Context, nodeCpuTopology *cslabecentuagrv1alpha1.NodeCpuTopologyApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.NodeCpuTopology, err error)
	ApplyStatus(ctx context.Context, nodeCpuTopology *cslabecentuagrv1alpha1.NodeCpuTopologyApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.NodeCpuTopology, err error)
	NodeCpuTopologyExpansion
}

// nodeCpuTopologies implements NodeCpuTopologyInterface
type nodeCpuTopologies struct {
	client rest.Interface
}

// newNodeCpuTopologies returns a NodeCpuTopologies
func newNodeCpuTopologies(c *CslabV1alpha1Client) *nodeCpuTopologies {
	return &nodeCpuTopologies{
		client: c.RESTClient(),
	}
}

// Get takes name of the nodeCpuTopology, and returns the corresponding nodeCpuTopology object, and an error if there is any.
func (c *nodeCpuTopologies) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.NodeCpuTopology, err error) {
	result = &v1alpha1.NodeCpuTopology{}
	err = c.client.Get().
		Resource("nodecputopologies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of NodeCpuTopologies that match those selectors.
func (c *nodeCpuTopologies) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.NodeCpuTopologyList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.NodeCpuTopologyList{}
	err = c.client.Get().
		Resource("nodecputopologies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested nodeCpuTopologies.
func (c *nodeCpuTopologies) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("nodecputopologies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a nodeCpuTopology and creates it.  Returns the server's representation of the nodeCpuTopology, and an error, if there is any.
func (c *nodeCpuTopologies) Create(ctx context.Context, nodeCpuTopology *v1alpha1.NodeCpuTopology, opts v1.CreateOptions) (result *v1alpha1.NodeCpuTopology, err error) {
	result = &v1alpha1.NodeCpuTopology{}
	err = c.client.Post().
		Resource("nodecputopologies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(nodeCpuTopology).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a nodeCpuTopology and updates it. Returns the server's representation of the nodeCpuTopology, and an error, if there is any.
func (c *nodeCpuTopologies) Update(ctx context.Context, nodeCpuTopology *v1alpha1.NodeCpuTopology, opts v1.UpdateOptions) (result *v1alpha1.NodeCpuTopology, err error) {
	result = &v1alpha1.NodeCpuTopology{}
	err = c.client.Put().
		Resource("nodecputopologies").
		Name(nodeCpuTopology.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(nodeCpuTopology).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *nodeCpuTopologies) UpdateStatus(ctx context.Context, nodeCpuTopology *v1alpha1.NodeCpuTopology, opts v1.UpdateOptions) (result *v1alpha1.NodeCpuTopology, err error) {
	result = &v1alpha1.NodeCpuTopology{}
	err = c.client.Put().
		Resource("nodecputopologies").
		Name(nodeCpuTopology.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(nodeCpuTopology).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the nodeCpuTopology and deletes it. Returns an error if one occurs.
func (c *nodeCpuTopologies) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("nodecputopologies").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *nodeCpuTopologies) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("nodecputopologies").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched nodeCpuTopology.
func (c *nodeCpuTopologies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.NodeCpuTopology, err error) {
	result = &v1alpha1.NodeCpuTopology{}
	err = c.client.Patch(pt).
		Resource("nodecputopologies").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied nodeCpuTopology.
func (c *nodeCpuTopologies) Apply(ctx context.Context, nodeCpuTopology *cslabecentuagrv1alpha1.NodeCpuTopologyApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.NodeCpuTopology, err error) {
	if nodeCpuTopology == nil {
		return nil, fmt.Errorf("nodeCpuTopology provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(nodeCpuTopology)
	if err != nil {
		return nil, err
	}
	name := nodeCpuTopology.Name
	if name == nil {
		return nil, fmt.Errorf("nodeCpuTopology.Name must be provided to Apply")
	}
	result = &v1alpha1.NodeCpuTopology{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("nodecputopologies").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *nodeCpuTopologies) ApplyStatus(ctx context.Context, nodeCpuTopology *cslabecentuagrv1alpha1.NodeCpuTopologyApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.NodeCpuTopology, err error) {
	if nodeCpuTopology == nil {
		return nil, fmt.Errorf("nodeCpuTopology provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(nodeCpuTopology)
	if err != nil {
		return nil, err
	}

	name := nodeCpuTopology.Name
	if name == nil {
		return nil, fmt.Errorf("nodeCpuTopology.Name must be provided to Apply")
	}

	result = &v1alpha1.NodeCpuTopology{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("nodecputopologies").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
