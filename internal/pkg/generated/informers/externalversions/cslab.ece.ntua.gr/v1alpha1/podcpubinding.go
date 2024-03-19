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
// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	cslabecentuagrv1alpha1 "cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	versioned "cslab.ece.ntua.gr/actimanager/internal/pkg/generated/clientset/versioned"
	internalinterfaces "cslab.ece.ntua.gr/actimanager/internal/pkg/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "cslab.ece.ntua.gr/actimanager/internal/pkg/generated/listers/cslab.ece.ntua.gr/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// PodCPUBindingInformer provides access to a shared informer and lister for
// PodCPUBindings.
type PodCPUBindingInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.PodCPUBindingLister
}

type podCPUBindingInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewPodCPUBindingInformer constructs a new informer for PodCPUBinding type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewPodCPUBindingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredPodCPUBindingInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredPodCPUBindingInformer constructs a new informer for PodCPUBinding type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredPodCPUBindingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CslabV1alpha1().PodCPUBindings(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CslabV1alpha1().PodCPUBindings(namespace).Watch(context.TODO(), options)
			},
		},
		&cslabecentuagrv1alpha1.PodCPUBinding{},
		resyncPeriod,
		indexers,
	)
}

func (f *podCPUBindingInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredPodCPUBindingInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *podCPUBindingInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&cslabecentuagrv1alpha1.PodCPUBinding{}, f.defaultInformer)
}

func (f *podCPUBindingInformer) Lister() v1alpha1.PodCPUBindingLister {
	return v1alpha1.NewPodCPUBindingLister(f.Informer().GetIndexer())
}
