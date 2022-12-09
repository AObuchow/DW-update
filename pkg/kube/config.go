package kube

import (
	"fmt"
	"path/filepath"

	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Get kube config depending on whether we're in a pod or running locally
func GetKubeConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		if err == rest.ErrNotInCluster {
			config, err = outClusterConfig()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return config, nil
}

// TODO: Cleanup this function..
func outClusterConfig() (config *rest.Config, err error) {
	if home := homedir.HomeDir(); home != "" {
		kubeconfig := filepath.Join(home, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)

		if err != nil {
			return nil, err
		}
	} else {
		fmt.Fprintf(os.Stderr, "Couldn't find ~/.kube/config file and not running in a pod. Exiting.")
		os.Exit(1)
	}
	return config, nil
}
