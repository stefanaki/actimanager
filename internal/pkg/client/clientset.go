package client

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
)

func New() (*kubernetes.Clientset, error) {
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
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
