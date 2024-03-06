package client

import (
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
)

type PodClient struct {
	client          kubernetes.Clientset
	informer        cache.SharedIndexInformer
	informerFactory informers.SharedInformerFactory
	stopCh          *chan struct{}
	logger          logr.Logger
}

func NewPodClient(clientset kubernetes.Clientset, logger logr.Logger) (*PodClient, error) {
	client := &PodClient{}
	informerFactory := informers.NewSharedInformerFactory(&clientset, 30*time.Second)
	informer := informerFactory.Core().V1().Pods().Informer()
	err := informer.AddIndexers(cache.Indexers{
		"nodeName": func(obj interface{}) ([]string, error) {
			pod, ok := obj.(*corev1.Pod)
			if !ok {
				return []string{}, errors.New("failed to use pod index")
			}
			return []string{pod.Spec.NodeName}, nil
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create node name index for pod: %v", err)
	}
	client.client = clientset
	client.informerFactory = informerFactory
	client.informer = informer
	client.logger = logger.WithName("pod-client")
	return client, nil
}

func (c *PodClient) Start(stopCh *chan struct{}) error {
	c.stopCh = stopCh
	c.informerFactory.Start(*stopCh)
	c.logger.Info("Starting Pod informer")
	c.logger.Info("Waiting for cache to sync")
	if ok := cache.WaitForCacheSync(*stopCh, c.informer.HasSynced); !ok {
		return errors.New("failed to sync cache")
	}
	return nil
}

func (c *PodClient) Stop() {
	c.logger.Info("Stopping Pod informer")
	*c.stopCh <- struct{}{}
}

func (c *PodClient) PodsForNode(nodeName string) ([]corev1.Pod, error) {
	pods, err := c.informer.GetIndexer().ByIndex("nodeName", nodeName)
	if err != nil {
		return nil, err
	}
	var res []corev1.Pod
	for _, pod := range pods {
		res = append(res, *pod.(*corev1.Pod).DeepCopy())
	}
	return res, nil
}
