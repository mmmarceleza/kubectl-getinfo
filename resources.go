package main

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// getGVR returns the GroupVersionResource for a given resource type
// It uses the Kubernetes API discovery to resolve resource names, kinds, and short names
func getGVR(resourceType string, config *rest.Config) (schema.GroupVersionResource, bool, error) {
	// Create discovery client to query API resources
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return schema.GroupVersionResource{}, false, fmt.Errorf("error creating discovery client: %v", err)
	}

	// Get all API resources from the cluster
	_, apiResourceLists, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		// Handle partial discovery errors (some groups may fail but others succeed)
		if apiResourceLists == nil {
			return schema.GroupVersionResource{}, false, fmt.Errorf("API discovery failed: %v", err)
		}
		// Continue with partial results
	}

	// Normalize resource type for comparison (case-insensitive)
	resourceTypeLower := strings.ToLower(resourceType)

	// Search for the resource type across all API groups
	for _, apiResourceList := range apiResourceLists {
		if apiResourceList == nil {
			continue
		}

		for _, apiResource := range apiResourceList.APIResources {
			// Skip subresources (e.g., pods/status, pods/log)
			if strings.Contains(apiResource.Name, "/") {
				continue
			}

			// Check if the input matches:
			// 1. Resource name (e.g., "pods", "deployments")
			// 2. Kind (e.g., "Pod", "Deployment")
			// 3. Short names (e.g., "po", "deploy", "svc")
			resourceNameLower := strings.ToLower(apiResource.Name)
			kindLower := strings.ToLower(apiResource.Kind)

			matched := resourceNameLower == resourceTypeLower || kindLower == resourceTypeLower

			// Check short names if not matched yet
			if !matched {
				for _, shortName := range apiResource.ShortNames {
					if strings.ToLower(shortName) == resourceTypeLower {
						matched = true
						break
					}
				}
			}

			if matched {
				// Parse group and version from the group version string
				gv, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
				if err != nil {
					continue
				}

				gvr := schema.GroupVersionResource{
					Group:    gv.Group,
					Version:  gv.Version,
					Resource: apiResource.Name,
				}

				return gvr, apiResource.Namespaced, nil
			}
		}
	}

	return schema.GroupVersionResource{}, false, fmt.Errorf("resource type '%s' not found in cluster", resourceType)
}

// getResources retrieves resources from the Kubernetes API
func getResources(
	client dynamic.Interface,
	gvr schema.GroupVersionResource,
	namespaced bool,
	namespace string,
	resourceNames []string,
	labelSelector labels.Selector,
) ([]unstructured.Unstructured, error) {
	ctx := context.Background()

	var resourceInterface dynamic.ResourceInterface
	if namespaced {
		if namespace == "" {
			resourceInterface = client.Resource(gvr)
		} else {
			resourceInterface = client.Resource(gvr).Namespace(namespace)
		}
	} else {
		resourceInterface = client.Resource(gvr)
	}

	var items []unstructured.Unstructured

	// If specific resource names are provided, get them individually
	if len(resourceNames) > 0 {
		for _, name := range resourceNames {
			item, err := resourceInterface.Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				return nil, fmt.Errorf("error getting %s: %v", name, err)
			}
			items = append(items, *item)
		}
	} else {
		// List all resources
		listOptions := metav1.ListOptions{}
		if labelSelector != nil {
			listOptions.LabelSelector = labelSelector.String()
		}

		list, err := resourceInterface.List(ctx, listOptions)
		if err != nil {
			return nil, fmt.Errorf("error listing resources: %v", err)
		}

		items = list.Items
	}

	return items, nil
}
