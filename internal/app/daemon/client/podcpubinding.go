package client

import (
	"errors"
	"fmt"
	"time"

	"cslab.ece.ntua.gr/actimanager/api/cslab.ece.ntua.gr/v1alpha1"
	clientset "cslab.ece.ntua.gr/actimanager/internal/pkg/generated/clientset/versioned"
	"cslab.ece.ntua.gr/actimanager/internal/pkg/generated/informers/externalversions"
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/cache"
)

type PodCPUBindingClient struct {
	client          clientset.Clientset
	informer        cache.SharedIndexInformer
	informerFactory externalversions.SharedInformerFactory
	stopCh          chan struct{}
	logger          logr.Logger
}

func NewPodCPUBindingClient(cslabClient clientset.Clientset, logger logr.Logger) (*PodCPUBindingClient, error) {
	client := &PodCPUBindingClient{}
	informerFactory := externalversions.NewSharedInformerFactory(&cslabClient, 30*time.Second)
	informer := informerFactory.Cslab().V1alpha1().PodCPUBindings().Informer()
	err := informer.AddIndexers(cache.Indexers{
		"nodeName": func(obj interface{}) ([]string, error) {
			pcb, ok := obj.(*v1alpha1.PodCPUBinding)
			if !ok {
				return []string{}, fmt.Errorf("failed to use podcpubinding index")
			}
			return []string{pcb.Status.NodeName}, nil
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create node name index for podcpubinding: %v", err)
	}
	client.client = cslabClient
	client.informer = informer
	client.informerFactory = informerFactory
	client.logger = logger.WithName("podcpubinding-client")
	return client, nil
}

func (c *PodCPUBindingClient) Start() error {
	c.stopCh = make(chan struct{})
	c.informerFactory.Start(c.stopCh)
	c.logger.Info("Starting PodCPUBinding informer")
	c.logger.Info("Waiting for cache to sync")
	if ok := cache.WaitForCacheSync(c.stopCh, c.informer.HasSynced); !ok {
		return errors.New("failed to sync cache")
	}
	return nil
}

func (c *PodCPUBindingClient) Stop() {
	c.logger.Info("Stopping PodCPUBinding informer")
	c.stopCh <- struct{}{}
}

func (c *PodCPUBindingClient) PodCPUBindingsForNode(nodeName string) ([]v1alpha1.PodCPUBinding, error) {
	b, err := c.informer.GetIndexer().ByIndex("nodeName", nodeName)
	if err != nil {
		return nil, fmt.Errorf("failed to get podcpubinding for node %s: %v", nodeName, err)
	}
	bindings := make([]v1alpha1.PodCPUBinding, 0)
	for _, obj := range b {
		binding, ok := obj.(*v1alpha1.PodCPUBinding)
		if !ok {
			return nil, fmt.Errorf("failed to cast podcpubinding for node %s", nodeName)
		}
		bindings = append(bindings, *binding)
	}
	return bindings, nil
}
