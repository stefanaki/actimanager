package client

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func GetClientSet() (*kubernetes.Clientset, error) {
	home := homedir.HomeDir()
	var kubeconfig = filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		errors.Errorf("error creating clientset from flags: %v\n", err.Error())

		config, err = rest.InClusterConfig()
		if err != nil {
			errors.Errorf("error creating clientset from in-cluster config: %v\n", err.Error())
			return nil, err
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
