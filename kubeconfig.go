package main

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// getKubeconfig returns the Kubernetes REST config
func getKubeconfig() (*rest.Config, error) {
	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// Try kubeconfig file
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting home directory: %v", err)
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("error building config from kubeconfig: %v", err)
	}

	return config, nil
}

// getCurrentNamespace returns the namespace from the current kubeconfig context
func getCurrentNamespace() string {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "default"
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return "default"
	}

	contextName := config.CurrentContext
	if contextName == "" {
		return "default"
	}

	context, exists := config.Contexts[contextName]
	if !exists || context == nil {
		return "default"
	}

	if context.Namespace != "" {
		return context.Namespace
	}

	return "default"
}

