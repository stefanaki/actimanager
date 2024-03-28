package client

import (
	clientset "cslab.ece.ntua.gr/actimanager/internal/pkg/generated/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
)

func NewClient() (*kubernetes.Clientset, error) {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Printf("error creating clientset from flags: %v\n", err.Error())

		config, err = rest.InClusterConfig()
		if err != nil {
			log.Printf("error creating clientset from in-cluster config: %v\n", err.Error())
			return nil, err
		}
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewCSLabClient() (*clientset.Clientset, error) {
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Printf("error creating cslab clientset from flags: %v\n", err.Error())

		config, err = rest.InClusterConfig()
		if err != nil {
			log.Printf("error creating cslab clientset from in-cluster config: %v\n", err.Error())
			return nil, err
		}
	}
	client, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
