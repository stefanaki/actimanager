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

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	v1alpha1 "cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/internal/pkg/generated/applyconfiguration/cslab.ece.ntua.gr/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePodCPUBindings implements PodCPUBindingInterface
type FakePodCPUBindings struct {
	Fake *FakeCslabV1alpha1
	ns   string
}

var podcpubindingsResource = v1alpha1.SchemeGroupVersion.WithResource("podcpubindings")

var podcpubindingsKind = v1alpha1.SchemeGroupVersion.WithKind("PodCPUBinding")

// Get takes name of the podCPUBinding, and returns the corresponding podCPUBinding object, and an error if there is any.
func (c *FakePodCPUBindings) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.PodCPUBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(podcpubindingsResource, c.ns, name), &v1alpha1.PodCPUBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PodCPUBinding), err
}

// List takes label and field selectors, and returns the list of PodCPUBindings that match those selectors.
func (c *FakePodCPUBindings) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.PodCPUBindingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(podcpubindingsResource, podcpubindingsKind, c.ns, opts), &v1alpha1.PodCPUBindingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.PodCPUBindingList{ListMeta: obj.(*v1alpha1.PodCPUBindingList).ListMeta}
	for _, item := range obj.(*v1alpha1.PodCPUBindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested podCPUBindings.
func (c *FakePodCPUBindings) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(podcpubindingsResource, c.ns, opts))

}

// Create takes the representation of a podCPUBinding and creates it.  Returns the server's representation of the podCPUBinding, and an error, if there is any.
func (c *FakePodCPUBindings) Create(ctx context.Context, podCPUBinding *v1alpha1.PodCPUBinding, opts v1.CreateOptions) (result *v1alpha1.PodCPUBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(podcpubindingsResource, c.ns, podCPUBinding), &v1alpha1.PodCPUBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PodCPUBinding), err
}

// Update takes the representation of a podCPUBinding and updates it. Returns the server's representation of the podCPUBinding, and an error, if there is any.
func (c *FakePodCPUBindings) Update(ctx context.Context, podCPUBinding *v1alpha1.PodCPUBinding, opts v1.UpdateOptions) (result *v1alpha1.PodCPUBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(podcpubindingsResource, c.ns, podCPUBinding), &v1alpha1.PodCPUBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PodCPUBinding), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakePodCPUBindings) UpdateStatus(ctx context.Context, podCPUBinding *v1alpha1.PodCPUBinding, opts v1.UpdateOptions) (*v1alpha1.PodCPUBinding, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(podcpubindingsResource, "status", c.ns, podCPUBinding), &v1alpha1.PodCPUBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PodCPUBinding), err
}

// Delete takes name of the podCPUBinding and deletes it. Returns an error if one occurs.
func (c *FakePodCPUBindings) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(podcpubindingsResource, c.ns, name, opts), &v1alpha1.PodCPUBinding{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePodCPUBindings) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(podcpubindingsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.PodCPUBindingList{})
	return err
}

// Patch applies the patch and returns the patched podCPUBinding.
func (c *FakePodCPUBindings) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.PodCPUBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(podcpubindingsResource, c.ns, name, pt, data, subresources...), &v1alpha1.PodCPUBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PodCPUBinding), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied podCPUBinding.
func (c *FakePodCPUBindings) Apply(ctx context.Context, podCPUBinding *cslabecentuagrv1alpha1.PodCPUBindingApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.PodCPUBinding, err error) {
	if podCPUBinding == nil {
		return nil, fmt.Errorf("podCPUBinding provided to Apply must not be nil")
	}
	data, err := json.Marshal(podCPUBinding)
	if err != nil {
		return nil, err
	}
	name := podCPUBinding.Name
	if name == nil {
		return nil, fmt.Errorf("podCPUBinding.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(podcpubindingsResource, c.ns, *name, types.ApplyPatchType, data), &v1alpha1.PodCPUBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PodCPUBinding), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakePodCPUBindings) ApplyStatus(ctx context.Context, podCPUBinding *cslabecentuagrv1alpha1.PodCPUBindingApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.PodCPUBinding, err error) {
	if podCPUBinding == nil {
		return nil, fmt.Errorf("podCPUBinding provided to Apply must not be nil")
	}
	data, err := json.Marshal(podCPUBinding)
	if err != nil {
		return nil, err
	}
	name := podCPUBinding.Name
	if name == nil {
		return nil, fmt.Errorf("podCPUBinding.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(podcpubindingsResource, c.ns, *name, types.ApplyPatchType, data, "status"), &v1alpha1.PodCPUBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PodCPUBinding), err
}
